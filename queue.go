package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"integrator/internal/database"
	"io"
	"log"
	"net/http"
	"objects"
	"shopify"
	"strconv"
	"sync"
	"time"
	"utils"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// TODO where are queue errors inside the queue worker logged?
// are they logged as part of the error message for the queue_item?
// retryable?  - no dont want to make this a function, but it should return the status if it failed...

func QueueWorker(dbconfig *DbConfig) {
	if dbconfig.Valid {
		go LoopQueueWorker(dbconfig)
	}
}

func LoopQueueWorker(dbconfig *DbConfig) {
	interval := 5
	db_interval, err := dbconfig.DB.GetAppSettingByKey(context.Background(), "app_queue_cron_time")
	if err != nil {
		interval = 5
	}
	interval, err = strconv.Atoi(db_interval.Value)
	if err != nil {
		interval = 5
	}
	timer := time.Duration(interval * int(time.Second))
	ticker := time.NewTicker(timer)
	for ; ; <-ticker.C {
		QueueWaitGroup(dbconfig)
	}
}

func QueueWaitGroup(dbconfig *DbConfig) {
	fmt.Println("inside gorountine")
	waitgroup := &sync.WaitGroup{}
	process_limit := 0
	process_limit_db, err := dbconfig.DB.GetAppSettingByKey(context.Background(), "app_queue_process_limit")
	if err != nil {
		process_limit = 20
	}
	process_limit, err = strconv.Atoi(process_limit_db.Value)
	if err != nil {
		process_limit = 20
	}
	queue_items, err := dbconfig.DB.GetQueueItemsByDate(context.Background(), database.GetQueueItemsByDateParams{
		Status: "in-queue",
		Limit:  int32(process_limit),
		Offset: 0,
	})
	if err != nil {
		log.Println(err)
		return
	}
	for _, queue_item := range queue_items {
		waitgroup.Add(1)
		go dbconfig.QueuePopAndProcess(queue_item.QueueType)
	}
	waitgroup.Wait()
}

// program starts
// items are added to the queue inside the database
// worker will then pick up the items in the queue that has been newly added
// should run whenever there is a new item added to the queue
// until the queue is empty
// worker will process them and mark them completed

// POST /api/queue
func (dbconfig *DbConfig) QueuePush(w http.ResponseWriter, r *http.Request, user database.User) {
	queue_size_int := 0
	queue_size_db, err := dbconfig.DB.GetAppSettingByKey(context.Background(), "app_queue_size")
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		queue_size_int = 100
	}
	queue_size_int, err = strconv.Atoi(queue_size_db.Value)
	if err != nil {
		queue_size_int = 100
	}
	size, err := dbconfig.DB.GetQueueSize(context.Background())
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if size >= int64(queue_size_int) {
		RespondWithError(w, http.StatusBadRequest, "queue is full, please wait")
		return
	}
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "error checking queue size")
		return
	}
	body, err := DecodeQueueItem(r)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	err = QueueItemValidation(body)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	raw, err := json.Marshal(&body.Object)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	queue_id, err := dbconfig.DB.CreateQueueItem(r.Context(), database.CreateQueueItemParams{
		ID:          uuid.New(),
		Object:      raw,
		QueueType:   body.Type,
		Instruction: body.Instruction,
		Status:      body.Status,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	})
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusCreated, objects.ResponseQueueItem{
		ID:     queue_id,
		Object: body,
	})
}

// Not an endpoint anymore, it's processes automatically in the background each 5 seconds
// POST /api/queue/worker?type=orders,products
func (dbconfig *DbConfig) QueuePopAndProcess(worker_type string) (database.QueueItem, error) {
	err := CheckWorkerType(worker_type)
	if err != nil {
		return database.QueueItem{}, err
	}
	queue_item, err := dbconfig.DB.GetNextQueueItem(context.Background())
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return database.QueueItem{}, errors.New("queue is empty")
		}
		return database.QueueItem{}, err
	}
	if worker_type == "product" {
		push_enabled, err := dbconfig.DB.GetAppSettingByKey(context.Background(), "app_enable_shopify_push")
		if err != nil {
			return database.QueueItem{}, err
		}
		is_enabled, err := strconv.ParseBool(push_enabled.Value)
		if err != nil {
			return database.QueueItem{}, err
		}
		if !is_enabled {
			return database.QueueItem{}, errors.New("product push is disabled")
		}
		item, err := dbconfig.DB.GetQueueItemsByStatusAndType(context.Background(), database.GetQueueItemsByStatusAndTypeParams{
			Status:    "processing",
			QueueType: "product",
			Limit:     1,
			Offset:    0,
		})
		if err != nil {
			if err.Error() != "sql: no rows in result set" {
				return database.QueueItem{}, errors.New("product queue is currently empty")
			}
		}
		if len(item) != 0 {
			if item[0].QueueType == "product" {
				return database.QueueItem{}, errors.New("product queue is currently busy")
			}
		}
	} else if worker_type == "order" {
		item, err := dbconfig.DB.GetQueueItemsByStatusAndType(context.Background(), database.GetQueueItemsByStatusAndTypeParams{
			Status:    "processing",
			QueueType: "order",
			Limit:     1,
			Offset:    0,
		})
		if err != nil {
			if err.Error() != "sql: no rows in result set" {
				return database.QueueItem{}, errors.New("order queue is currently empty")
			}
		}
		if len(item) != 0 {
			if item[0].QueueType == "order" {
				return database.QueueItem{}, errors.New("order queue is currently busy")
			}
		}
	}
	_, err = dbconfig.DB.UpdateQueueItem(context.Background(), database.UpdateQueueItemParams{
		Status:    "processing",
		UpdatedAt: time.Now().UTC(),
		ID:        queue_item.ID,
	})
	if err != nil {
		return database.QueueItem{}, err
	}
	err = ProcessQueueItem(dbconfig, queue_item)
	if err != nil {
		return database.QueueItem{}, err
	}
	updated_queue_item, err := dbconfig.DB.UpdateQueueItem(context.Background(), database.UpdateQueueItemParams{
		Status:    "completed",
		UpdatedAt: time.Now().UTC(),
		ID:        queue_item.ID,
	})
	if err != nil {
		return database.QueueItem{}, err
	}
	return updated_queue_item, nil
}

// field can be:
// - status: processing, completed, in-queue
// - instruction: add_order, update_order, add_product, update_product, add_variant, update_variant
// - type: product, order

// GET /api/queue/filter?key=value
func (dbconfig *DbConfig) FilterQueueItems(w http.ResponseWriter, r *http.Request, user database.User) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}
	param_type := utils.ConfirmFilters(r.URL.Query().Get("type"))
	param_instruction := utils.ConfirmFilters(r.URL.Query().Get("instruction"))
	param_status := utils.ConfirmFilters(r.URL.Query().Get("status"))
	result, err := CompileQueueFilterSearch(dbconfig, r.Context(), page, param_type, param_status, param_instruction)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusOK, result)
}

// GET /api/queue/processing
func (dbconfig *DbConfig) QueueViewCurrentItem(w http.ResponseWriter, r *http.Request, user database.User) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}
	queue_items, err := dbconfig.DB.GetQueueItemsByDate(r.Context(), database.GetQueueItemsByDateParams{
		Status: "in-queue",
		Limit:  10,
		Offset: int32((page - 1) * 10),
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
	}
	RespondWithJSON(w, http.StatusOK, queue_items)
}

// GET /api/queue?page=1
func (dbconfig *DbConfig) QueueViewNextItems(w http.ResponseWriter, r *http.Request, user database.User) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}
	queue_items, err := dbconfig.DB.GetQueueItemsByDate(r.Context(), database.GetQueueItemsByDateParams{
		Status: "in-queue",
		Limit:  10,
		Offset: int32((page - 1) * 10),
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
	}
	if queue_items == nil {
		RespondWithJSON(w, http.StatusOK, []string{})
		return
	} else {
		RespondWithJSON(w, http.StatusOK, queue_items)
	}
}

// GET /api/queue/view
func (dbconfig *DbConfig) QueueView(
	w http.ResponseWriter,
	r *http.Request,
	user database.User) {
	response, err := dbconfig.DisplayQueueCount()
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusOK, response)
}

// DELETE /api/queue/{id}
func (dbconfig *DbConfig) ClearQueueByID(
	w http.ResponseWriter,
	r *http.Request,
	user database.User) {
	id := chi.URLParam(r, "id")
	id_uuid, err := uuid.Parse(id)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	err = dbconfig.DB.RemoveQueueItemByID(r.Context(), id_uuid)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusOK, "success")
}

// DELETE /api/queue?key=value
func (dbconfig *DbConfig) ClearQueueByFilter(
	w http.ResponseWriter,
	r *http.Request,
	user database.User) {
	param_type := utils.ConfirmFilters(r.URL.Query().Get("type"))
	param_instruction := utils.ConfirmFilters(r.URL.Query().Get("instruction"))
	param_status := utils.ConfirmFilters(r.URL.Query().Get("status"))
	response, err := CompileRemoveQueueFilter(dbconfig, r.Context(), param_type, param_status, param_instruction)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusOK, response)
}

// Process a queue item
func ProcessQueueItem(dbconfig *DbConfig, queue_item database.QueueItem) error {
	if queue_item.QueueType == "product" {
		queue_object, err := DecodeQueueItemProduct(queue_item.Object)
		if err != nil {
			return err
		}
		product_id, err := uuid.Parse(queue_object.SystemProductID)
		if err != nil {
			return errors.New("could not decode product_id: " + queue_object.SystemProductID)
		}
		shopifyConfig := shopify.InitConfigShopify()
		product, err := CompileProductData(dbconfig, product_id, context.Background(), false)
		if err != nil {
			return err
		}
		if queue_item.Instruction == "add_product" {
			return dbconfig.PushProduct(&shopifyConfig, product)
		} else if queue_item.Instruction == "update_product" {
			return dbconfig.PushProduct(&shopifyConfig, product)
		} else {
			return errors.New("invalid product instruction")
		}
	} else if queue_item.QueueType == "order" {
		if queue_item.Instruction == "add_order" {
			queue_object, err := DecodeQueueItemOrder(queue_item.Object)
			if err != nil {
				return err
			}
			return dbconfig.AddOrder(queue_object)
		} else if queue_item.Instruction == "update_order" {
			queue_object, err := DecodeQueueItemOrder(queue_item.Object)
			if err != nil {
				return err
			}
			return dbconfig.UpdateOrder(queue_object)
		} else {
			return errors.New("invalid order instruction")
		}
	} else if queue_item.QueueType == "product_variant" {
		shopifyConfig := shopify.InitConfigShopify()
		queue_object, err := DecodeQueueItemProduct(queue_item.Object)
		if err != nil {
			return err
		}
		variant_id, err := uuid.Parse(queue_object.SystemVariantID)
		if err != nil {
			return errors.New("could not decode variant_id: " + queue_object.SystemVariantID)
		}
		variant, err := CompileVariantData(dbconfig, variant_id, context.Background())
		if err != nil {
			return err
		}
		if queue_item.Instruction == "add_variant" {
			return dbconfig.PushVariant(
				&shopifyConfig,
				variant,
				queue_object.Shopify.ProductID,
				queue_object.Shopify.VariantID,
			)

		} else if queue_item.Instruction == "update_variant" {
			return dbconfig.PushVariant(
				&shopifyConfig,
				variant,
				queue_object.Shopify.ProductID,
				queue_object.Shopify.VariantID,
			)
		} else {
			return errors.New("invalid product_variant instruction")
		}
	}
	return errors.New("invalid queue item type")
}

// helper function: Displays count of different instructions
func (dbconfig *DbConfig) DisplayQueueCount() (objects.ResponseQueueCount, error) {
	add_order, err := dbconfig.DB.GetQueueItemsCount(context.Background(), "add_order")
	if err != nil {
		return objects.ResponseQueueCount{}, err
	}
	add_product, err := dbconfig.DB.GetQueueItemsCount(context.Background(), "add_product")
	if err != nil {
		return objects.ResponseQueueCount{}, err
	}
	add_variant, err := dbconfig.DB.GetQueueItemsCount(context.Background(), "add_variant")
	if err != nil {
		return objects.ResponseQueueCount{}, err
	}
	update_order, err := dbconfig.DB.GetQueueItemsCount(context.Background(), "update_order")
	if err != nil {
		return objects.ResponseQueueCount{}, err
	}
	update_product, err := dbconfig.DB.GetQueueItemsCount(context.Background(), "update_product")
	if err != nil {
		return objects.ResponseQueueCount{}, err
	}
	update_variant, err := dbconfig.DB.GetQueueItemsCount(context.Background(), "update_variant")
	if err != nil {
		return objects.ResponseQueueCount{}, err
	}
	return objects.ResponseQueueCount{
		AddOrder:      int(add_order),
		AddProduct:    int(add_product),
		AddVariant:    int(add_variant),
		UpdateOrder:   int(update_order),
		UpdateProduct: int(update_product),
		UpdateVariant: int(update_variant),
	}, nil
}

// Compile Queue Filter Search into a single object (variable)
func CompileRemoveQueueFilter(
	dbconfig *DbConfig,
	ctx context.Context,
	queue_type,
	status,
	instruction string) (string, error) {
	if queue_type == "" {
		if status == "" {
			err := dbconfig.DB.RemoveQueueItemsByInstruction(ctx, instruction)
			if err != nil {
				return "error", err
			}
			return "success", nil
		} else {
			if instruction == "" {
				err := dbconfig.DB.RemoveQueueItemsByStatus(ctx, status)
				if err != nil {
					return "error", err
				}
				return "success", nil
			}
			err := dbconfig.DB.RemoveQueueItemsByStatusAndInstruction(
				ctx,
				database.RemoveQueueItemsByStatusAndInstructionParams{
					Instruction: instruction,
					Status:      status,
				})
			if err != nil {
				return "error", err
			}
			return "success", nil
		}
	}
	if status == "" {
		if instruction == "" {
			err := dbconfig.DB.RemoveQueueItemsByType(ctx, queue_type)
			if err != nil {
				return "error", err
			}
			return "success", nil
		} else {
			if queue_type == "" {
				err := dbconfig.DB.RemoveQueueItemsByInstruction(ctx, instruction)
				if err != nil {
					return "error", err
				}
				return "success", nil
			}
			err := dbconfig.DB.RemoveQueueItemsByTypeAndInstruction(
				ctx,
				database.RemoveQueueItemsByTypeAndInstructionParams{
					Instruction: instruction,
					QueueType:   queue_type,
				})
			if err != nil {
				return "error", err
			}
			return "success", nil
		}
	}
	if instruction == "" {
		if queue_type == "" {
			err := dbconfig.DB.RemoveQueueItemsByStatus(ctx, status)
			if err != nil {
				return "error", err
			}
			return "success", nil
		} else {
			if status == "" {
				err := dbconfig.DB.RemoveQueueItemsByType(ctx, queue_type)
				if err != nil {
					return "error", err
				}
				return "success", nil
			}
			err := dbconfig.DB.RemoveQueueItemsByStatusAndType(
				ctx,
				database.RemoveQueueItemsByStatusAndTypeParams{
					Status:    status,
					QueueType: queue_type,
				})
			if err != nil {
				return "error", err
			}
			return "success", nil
		}
	}
	err := dbconfig.DB.RemoveQueueItemsFilter(
		ctx,
		database.RemoveQueueItemsFilterParams{
			Status:      status,
			QueueType:   queue_type,
			Instruction: instruction,
		})
	if err != nil {
		return "error", err
	}
	return "success", nil
}

// Helper function: Checks if the worker type is valid
func CheckWorkerType(worker_type string) error {
	worker_types := []string{"product", "order"}
	for _, value := range worker_types {
		if value == worker_type {
			return nil
		}
	}
	return errors.New("invalid worker type")
}

// Helper function: Posts internal requests to the queue endpoint
func (dbconfig *DbConfig) QueueHelper(request_data objects.RequestQueueHelper, body io.Reader) (objects.ResponseQueueItem, error) {
	httpClient := http.Client{
		Timeout: time.Second * 10,
	}
	var buffer bytes.Buffer
	err := json.NewEncoder(&buffer).Encode(request_data)
	if err != nil {
		return objects.ResponseQueueItem{}, err
	}
	req, err := http.NewRequest(
		request_data.Method,
		"http://localhost:"+utils.LoadEnv("port")+"/api/"+request_data.Endpoint,
		&buffer,
	)
	if request_data.ApiKey != "" {
		req.Header.Add("Authorization", "ApiKey "+request_data.ApiKey)
	}
	if err != nil {
		return objects.ResponseQueueItem{}, err
	}
	res, err := httpClient.Do(req)
	if err != nil {
		return objects.ResponseQueueItem{}, err
	}
	defer res.Body.Close()
	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return objects.ResponseQueueItem{}, err
	}
	if res.StatusCode != 200 {
		return objects.ResponseQueueItem{}, err
	}
	queue_response := objects.ResponseQueueItem{}
	err = json.Unmarshal(respBody, &queue_response)
	if err != nil {
		return objects.ResponseQueueItem{}, err
	}
	return queue_response, nil
}

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
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

func QueueWorker(dbconfig *DbConfig) {
	go LoopQueueWorker(dbconfig)
}

func LoopQueueWorker(dbconfig *DbConfig) {
	interval := 5
	db_interval, err := dbconfig.DB.GetAppSettingByKey(context.Background(), "app_queue_cron_time")
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			log.Println(err)
			return
		}
		interval = 5
	}
	interval, err = strconv.Atoi(db_interval.Value)
	if err != nil {
		interval = 5
	}
	timer := time.Duration(interval * int(time.Second))
	ticker := time.NewTicker(timer)
	for ; ; <-ticker.C {
		queue_enabled := false
		queue_enabled_db, err := dbconfig.DB.GetAppSettingByKey(context.Background(), "app_enable_queue_worker")
		if err != nil {
			if err.Error() != "sql: no rows in result set" {
				log.Println(err)
				return
			}
			queue_enabled = false
		}
		queue_enabled, err = strconv.ParseBool(queue_enabled_db.Value)
		if err != nil {
			queue_enabled = false
		}
		if queue_enabled {
			QueueWaitGroup(dbconfig)
		}
	}
}

func QueueWaitGroup(dbconfig *DbConfig) {
	waitgroup := &sync.WaitGroup{}
	process_limit := int64(0)
	process_limit_db, err := dbconfig.DB.GetAppSettingByKey(context.Background(), "app_queue_process_limit")
	if err != nil {
		process_limit = 20
	}
	process_limit, err = strconv.ParseInt(process_limit_db.Value, 10, 32)
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
		dbconfig.QueuePopAndProcess(queue_item.QueueType, waitgroup)
	}
	waitgroup.Wait()
}

// POST /api/shopify/sync
func (dbconfig *DbConfig) Synchronize(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	// check if the syncro queue item exists in the queue already
	// if it does then it should throw an error
	item, err := dbconfig.DB.GetQueueItemsByInstructionAndStatus(r.Context(), database.GetQueueItemsByInstructionAndStatusParams{
		Instruction: "zsync_channel",
		Status:      "in-queue",
		Limit:       1,
		Offset:      0,
	})
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	if len(item) != 0 {
		RespondWithError(w, http.StatusInternalServerError, "sync in progress")
		return
	}
	page := 0
	for {
		// fetch all products paginated
		products, err := dbconfig.DB.GetActiveProducts(context.Background(), database.GetActiveProductsParams{
			Limit:  1,
			Offset: (int32(page) * 50),
		})
		if err != nil {
			if err.Error() != "sql: no rows in result set" {
				RespondWithError(w, http.StatusInternalServerError, err.Error())
				return
			}
		}
		if len(products) == 0 || page == 1 {
			break
		}
		for _, product := range products {
			product_compiled, err := CompileProductData(dbconfig, product.ID, r.Context(), false)
			if err != nil {
				RespondWithError(w, http.StatusInternalServerError, err.Error())
				return
			}
			// add product queue_items to the queue
			err = CompileInstructionProduct(dbconfig, product_compiled, dbUser)
			if err != nil {
				RespondWithError(w, http.StatusInternalServerError, err.Error())
				return
			}
			for _, variant := range product_compiled.Variants {
				// add variant queue_items to queue
				err = CompileInstructionVariant(dbconfig, variant, product_compiled, dbUser)
				if err != nil {
					RespondWithError(w, http.StatusInternalServerError, err.Error())
					return
				}
			}
		}
		page += 1
	}
	_, err = dbconfig.QueueHelper(objects.RequestQueueHelper{
		Type:        "product",
		Status:      "in-queue",
		Instruction: "zsync_channel",
		Endpoint:    "queue",
		ApiKey:      dbUser.ApiKey,
		Method:      http.MethodPost,
		Object:      nil,
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusOK, "synconizing started")
}

// GET /api/queue/{id}
func (dbconfig *DbConfig) GetQueueItemByID(
	w http.ResponseWriter,
	r *http.Request,
	user database.User) {
	id := chi.URLParam(r, "id")
	id_uuid, err := uuid.Parse(id)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	queue_item, err := dbconfig.DB.GetQueueItemByID(r.Context(), id_uuid)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusOK, queue_item)
}

// POST /api/queue
func (dbconfig *DbConfig) QueuePush(w http.ResponseWriter, r *http.Request, user database.User) {
	queue_size_int := 500
	queue_size_db, err := dbconfig.DB.GetAppSettingByKey(context.Background(), "app_queue_size")
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		queue_size_int = 500
	}
	queue_size_int, err = strconv.Atoi(queue_size_db.Value)
	if err != nil {
		queue_size_int = 500
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

// Not an endpoint anymore, it's processes automatically in the background
func (dbconfig *DbConfig) QueuePopAndProcess(worker_type string, wait_group *sync.WaitGroup) {
	defer wait_group.Done()
	queue_item, err := dbconfig.DB.GetNextQueueItem(context.Background())
	if err != nil {
		FailedQueueItem(dbconfig, queue_item, err)
		return
	}
	err = CheckWorkerType(worker_type)
	if err != nil {
		FailedQueueItem(dbconfig, queue_item, err)
		return
	}
	if worker_type == "product" || worker_type == "product_variant" {
		is_enabled := false
		push_enabled, err := dbconfig.DB.GetAppSettingByKey(context.Background(), "app_enable_shopify_push")
		if err != nil {
			if err.Error() != "sql: no rows in result set" {
				FailedQueueItem(dbconfig, queue_item, err)
				return
			}
		}
		if push_enabled.Value == "" {
			push_enabled.Value = "false"
		}
		is_enabled, err = strconv.ParseBool(push_enabled.Value)
		if err != nil {
			FailedQueueItem(dbconfig, queue_item, err)
			return
		}
		if !is_enabled {
			FailedQueueItem(dbconfig, queue_item, errors.New("product push is disabled"))
			return
		}
		item, err := dbconfig.DB.GetQueueItemsByStatusAndType(context.Background(), database.GetQueueItemsByStatusAndTypeParams{
			Status:    "processing",
			QueueType: worker_type,
			Limit:     1,
			Offset:    0,
		})
		if err != nil {
			if err.Error() != "sql: no rows in result set" {
				return
			}
		}
		if len(item) != 0 {
			if item[0].QueueType == worker_type {
				FailedQueueItem(dbconfig, queue_item, errors.New("product queue is currently busy"))
				return
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
				return
			}
		}
		if len(item) != 0 {
			if item[0].QueueType == "order" {
				FailedQueueItem(dbconfig, queue_item, errors.New("order queue is currently busy"))
				return
			}
		}
	}
	_, err = dbconfig.DB.UpdateQueueItem(context.Background(), database.UpdateQueueItemParams{
		Status:    "processing",
		UpdatedAt: time.Now().UTC(),
		ID:        queue_item.ID,
	})
	if err != nil {
		FailedQueueItem(dbconfig, queue_item, err)
		return
	}
	err = ProcessQueueItem(dbconfig, queue_item)
	if err != nil {
		FailedQueueItem(dbconfig, queue_item, err)
		return
	}
	_, err = dbconfig.DB.UpdateQueueItem(context.Background(), database.UpdateQueueItemParams{
		Status:    "completed",
		UpdatedAt: time.Now().UTC(),
		ID:        queue_item.ID,
	})
	if err != nil {
		FailedQueueItem(dbconfig, queue_item, err)
		return
	}
}

// field can be:
// - status: processing, completed, in-queue, failed
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
	RespondWithJSON(w, http.StatusOK, objects.ResponseString{
		Message: "success",
	})
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
	RespondWithJSON(w, http.StatusOK, objects.ResponseString{
		Message: response,
	})
}

// Process a queue item
func ProcessQueueItem(dbconfig *DbConfig, queue_item database.QueueItem) error {
	if queue_item.QueueType == "product" {
		if queue_item.Instruction == "zsync_channel" {
			return nil
		}
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
		restrictions, err := dbconfig.DB.GetPushRestriction(context.Background())
		if err != nil {
			return err
		}
		restrictions_map := PushRestrictionsToMap(restrictions)
		shopify_update_variant := ApplyPushRestrictionV(
			restrictions_map,
			ConvertVariantToShopify(variant),
		)
		if queue_item.Instruction == "add_variant" {
			return dbconfig.PushVariant(
				&shopifyConfig,
				variant,
				shopify_update_variant,
				restrictions_map,
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
	worker_types := []string{"product", "product_variant", "order"}
	for _, value := range worker_types {
		if value == worker_type {
			return nil
		}
	}
	return errors.New("invalid worker type")
}

// Updates the queue item if there is an error
func FailedQueueItem(dbconfig *DbConfig, queue_item database.QueueItem, err_r error) {
	_, err := dbconfig.DB.UpdateQueueItem(context.Background(), database.UpdateQueueItemParams{
		Status:      "failed",
		UpdatedAt:   time.Now().UTC(),
		Description: err_r.Error(),
		ID:          queue_item.ID,
	})
	if err != nil {
		log.Println(err)
		return
	}
}

// Helper function: Posts internal requests to the queue endpoint
func (dbconfig *DbConfig) QueueHelper(request_data objects.RequestQueueHelper) (objects.ResponseQueueItem, error) {
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
		// TODO change this call to be more dynamic
		"http://localhost:"+utils.LoadEnv("port")+"/api/"+request_data.Endpoint,
		&buffer,
	)
	if err != nil {
		return objects.ResponseQueueItem{}, err
	}
	if request_data.ApiKey != "" {
		req.Header.Add("Authorization", "ApiKey "+request_data.ApiKey)
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
	if res.StatusCode != 201 {
		return objects.ResponseQueueItem{}, errors.New(string(respBody))
	}
	queue_response := objects.ResponseQueueItem{}
	err = json.Unmarshal(respBody, &queue_response)
	if err != nil {
		return objects.ResponseQueueItem{}, err
	}
	return queue_response, nil
}

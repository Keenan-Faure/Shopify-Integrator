package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"integrator/internal/database"
	"io"
	"net/http"
	"objects"
	"shopify"
	"strconv"
	"time"
	"utils"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// TODO make a queue_size setting - would need another table to store app setting
// how about using the .env file?
const queue_size = 300

// POST /api/queue
func (dbconfig *DbConfig) QueuePush(w http.ResponseWriter, r *http.Request, user database.User) {
	// check if the queue_size has been reached
	size, err := dbconfig.DB.GetQueueSize(context.Background())
	if size >= queue_size {
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
	RespondWithJSON(w, http.StatusOK, objects.ResponseQueueItem{
		ID:     queue_id,
		Object: body,
	})
}

// POST /api/queue/worker?type=orders,products
func (dbconfig *DbConfig) QueuePopAndProcess(w http.ResponseWriter, r *http.Request, user database.User) {
	queue_item, err := dbconfig.DB.GetNextQueueItem(r.Context())
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			RespondWithError(w, http.StatusInternalServerError, "queue is empty.")
			return
		}
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if queue_item.QueueType == "product" {
		item, err := dbconfig.DB.GetQueueItemsByStatusAndType(r.Context(), database.GetQueueItemsByStatusAndTypeParams{
			Status:    "processing",
			QueueType: "product",
			Limit:     1,
			Offset:    0,
		})
		if err != nil {
			if err.Error() != "sql: no rows in result set" {
				RespondWithError(w, http.StatusInternalServerError, "queue is empty.")
				return
			}
		}
		if item[0].QueueType == "product" {
			RespondWithError(w, http.StatusServiceUnavailable, "product queue is currently busy")
			return
		}
	} else {
		item, err := dbconfig.DB.GetQueueItemsByStatusAndType(r.Context(), database.GetQueueItemsByStatusAndTypeParams{
			Status:    "processing",
			QueueType: "order",
			Limit:     1,
			Offset:    0,
		})
		if err != nil {
			if err.Error() != "sql: no rows in result set" {
				RespondWithError(w, http.StatusInternalServerError, "queue is empty.")
				return
			}
		}
		if item[0].QueueType == "order" {
			RespondWithError(w, http.StatusServiceUnavailable, "order queue is currently busy")
			return
		}
	}
	_, err = dbconfig.DB.UpdateQueueItem(r.Context(), database.UpdateQueueItemParams{
		Status:    "processing",
		UpdatedAt: time.Now().UTC(),
		ID:        queue_item.ID,
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = ProcessQueueItem(dbconfig, queue_item)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	updated_queue_item, err := dbconfig.DB.UpdateQueueItem(r.Context(), database.UpdateQueueItemParams{
		Status:    "completed",
		UpdatedAt: time.Now().UTC(),
		ID:        queue_item.ID,
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusOK, updated_queue_item)
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

// TODO make a seperate endpoint?
// GET /api/queue?position=0
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
	RespondWithJSON(w, http.StatusOK, []string{"success"})
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

// Helper function: removes queue items by filters

// Compile Queue Filter Search into a single object (variable)
func CompileRemoveQueueFilter(
	dbconfig *DbConfig,
	ctx context.Context,
	queue_type,
	status,
	instruction string) ([]string, error) {
	if queue_type == "" {
		if status == "" {
			err := dbconfig.DB.RemoveQueueItemsByInstruction(ctx, instruction)
			if err != nil {
				return []string{"error"}, err
			}
			return []string{"success"}, nil
		} else {
			if instruction == "" {
				err := dbconfig.DB.RemoveQueueItemsByStatus(ctx, status)
				if err != nil {
					return []string{"error"}, err
				}
				return []string{"success"}, nil
			}
			err := dbconfig.DB.RemoveQueueItemsByStatusAndInstruction(
				ctx,
				database.RemoveQueueItemsByStatusAndInstructionParams{
					Instruction: instruction,
					Status:      status,
				})
			if err != nil {
				return []string{"error"}, err
			}
			return []string{"success"}, nil
		}
	}
	if status == "" {
		if instruction == "" {
			err := dbconfig.DB.RemoveQueueItemsByType(ctx, queue_type)
			if err != nil {
				return []string{"error"}, err
			}
			return []string{"success"}, nil
		} else {
			if queue_type == "" {
				err := dbconfig.DB.RemoveQueueItemsByInstruction(ctx, instruction)
				if err != nil {
					return []string{"error"}, err
				}
				return []string{"success"}, nil
			}
			err := dbconfig.DB.RemoveQueueItemsByTypeAndInstruction(
				ctx,
				database.RemoveQueueItemsByTypeAndInstructionParams{
					Instruction: instruction,
					QueueType:   queue_type,
				})
			if err != nil {
				return []string{"error"}, err
			}
			return []string{"success"}, nil
		}
	}
	if instruction == "" {
		if queue_type == "" {
			err := dbconfig.DB.RemoveQueueItemsByStatus(ctx, status)
			if err != nil {
				return []string{"error"}, err
			}
			return []string{"success"}, nil
		} else {
			if status == "" {
				err := dbconfig.DB.RemoveQueueItemsByType(ctx, queue_type)
				if err != nil {
					return []string{"error"}, err
				}
				return []string{"success"}, nil
			}
			err := dbconfig.DB.RemoveQueueItemsByStatusAndType(
				ctx,
				database.RemoveQueueItemsByStatusAndTypeParams{
					Status:    status,
					QueueType: queue_type,
				})
			if err != nil {
				return []string{"error"}, err
			}
			return []string{"success"}, nil
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
		return []string{"error"}, err
	}
	return []string{"success"}, nil
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

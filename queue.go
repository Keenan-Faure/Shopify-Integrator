package main

import (
	"context"
	"errors"
	"integrator/internal/database"
	"log"
	"strconv"
	"sync"
	"time"
)

// import (
// 	"bytes"
// 	"context"
// 	"encoding/json"
// 	"errors"
// 	"integrator/internal/database"
// 	"io"
// 	"log"
// 	"net/http"
// 	"objects"
// 	"shopify"
// 	"strconv"
// 	"sync"
// 	"time"
// 	"utils"

// 	"github.com/go-chi/chi/v5"
// 	"github.com/google/uuid"
// )

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

// // POST /api/shopify/sync
// func (dbconfig *DbConfig) Synchronize(w http.ResponseWriter, r *http.Request, dbUser database.User) {
// 	// check if the syncro queue item exists in the queue already
// 	// if it does then it should throw an error
// 	item, err := dbconfig.DB.GetQueueItemsByInstructionAndStatus(r.Context(), database.GetQueueItemsByInstructionAndStatusParams{
// 		Instruction: "zsync_channel",
// 		Status:      "in-queue",
// 		Limit:       1,
// 		Offset:      0,
// 	})
// 	if err != nil {
// 		if err.Error() != "sql: no rows in result set" {
// 			RespondWithError(w, http.StatusInternalServerError, err.Error())
// 			return
// 		}
// 	}
// 	if len(item) != 0 {
// 		RespondWithError(w, http.StatusInternalServerError, "sync in progress")
// 		return
// 	}
// 	page := 0
// 	for {
// 		// fetch all products paginated
// 		products, err := dbconfig.DB.GetActiveProducts(context.Background(), database.GetActiveProductsParams{
// 			Limit:  50,
// 			Offset: (int32(page) * 50),
// 		})
// 		if err != nil {
// 			if err.Error() != "sql: no rows in result set" {
// 				RespondWithError(w, http.StatusInternalServerError, err.Error())
// 				return
// 			}
// 		}
// 		if len(products) == 0 {
// 			break
// 		}
// 		for _, product := range products {
// 			product_compiled, err := CompileProductData(dbconfig, product.ID, r.Context(), false)
// 			if err != nil {
// 				RespondWithError(w, http.StatusInternalServerError, err.Error())
// 				return
// 			}
// 			// add product queue_items to the queue
// 			err = CompileInstructionProduct(dbconfig, product_compiled, dbUser)
// 			if err != nil {
// 				RespondWithError(w, http.StatusInternalServerError, err.Error())
// 				return
// 			}
// 			for _, variant := range product_compiled.Variants {
// 				// add variant queue_items to queue
// 				err = CompileInstructionVariant(dbconfig, variant, product_compiled, dbUser)
// 				if err != nil {
// 					RespondWithError(w, http.StatusInternalServerError, err.Error())
// 					return
// 				}
// 			}
// 		}
// 		page += 1
// 	}
// 	_, err = dbconfig.QueueHelper(objects.RequestQueueHelper{
// 		Type:        "product",
// 		Status:      "in-queue",
// 		Instruction: "zsync_channel",
// 		Endpoint:    "queue",
// 		ApiKey:      dbUser.ApiKey,
// 		Method:      http.MethodPost,
// 		Object:      nil,
// 	})
// 	if err != nil {
// 		RespondWithError(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}
// 	RespondWithJSON(w, http.StatusOK, objects.ResponseString{
// 		Message: "synconizing started",
// 	})
// }

// // GET /api/queue/{id}
// func (dbconfig *DbConfig) GetQueueItemByID(
// 	w http.ResponseWriter,
// 	r *http.Request,
// 	user database.User) {
// 	id := chi.URLParam(r, "id")
// 	id_uuid, err := uuid.Parse(id)
// 	if err != nil {
// 		RespondWithError(w, http.StatusBadRequest, err.Error())
// 		return
// 	}
// 	queue_item, err := dbconfig.DB.GetQueueItemByID(r.Context(), id_uuid)
// 	if err != nil {
// 		RespondWithError(w, http.StatusBadRequest, err.Error())
// 		return
// 	}
// 	RespondWithJSON(w, http.StatusOK, queue_item)
// }

// // POST /api/queue
// func (dbconfig *DbConfig) QueuePush(w http.ResponseWriter, r *http.Request, user database.User) {
// 	queue_size_int := 500
// 	queue_size_db, err := dbconfig.DB.GetAppSettingByKey(context.Background(), "app_queue_size")
// 	if err != nil {
// 		if err.Error() != "sql: no rows in result set" {
// 			RespondWithError(w, http.StatusInternalServerError, err.Error())
// 			return
// 		}
// 		queue_size_int = 500
// 	}
// 	queue_size_int, err = strconv.Atoi(queue_size_db.Value)
// 	if err != nil {
// 		queue_size_int = 500
// 	}
// 	size, err := dbconfig.DB.GetQueueSize(context.Background())
// 	if err != nil {
// 		RespondWithError(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}
// 	if size >= int64(queue_size_int) {
// 		RespondWithError(w, http.StatusBadRequest, "queue is full, please wait")
// 		return
// 	}
// 	if err != nil {
// 		RespondWithError(w, http.StatusBadRequest, "error checking queue size")
// 		return
// 	}
// 	body, err := DecodeQueueItem(r)
// 	if err != nil {
// 		RespondWithError(w, http.StatusBadRequest, err.Error())
// 		return
// 	}
// 	err = QueueItemValidation(body)
// 	if err != nil {
// 		RespondWithError(w, http.StatusBadRequest, err.Error())
// 		return
// 	}
// 	raw, err := json.Marshal(&body.Object)
// 	if err != nil {
// 		RespondWithError(w, http.StatusBadRequest, err.Error())
// 		return
// 	}
// 	queue_id, err := dbconfig.DB.CreateQueueItem(r.Context(), database.CreateQueueItemParams{
// 		ID:          uuid.New(),
// 		Object:      raw,
// 		QueueType:   body.Type,
// 		Instruction: body.Instruction,
// 		Status:      body.Status,
// 		CreatedAt:   time.Now().UTC(),
// 		UpdatedAt:   time.Now().UTC(),
// 	})
// 	if err != nil {
// 		RespondWithError(w, http.StatusBadRequest, err.Error())
// 		return
// 	}
// 	RespondWithJSON(w, http.StatusCreated, objects.ResponseQueueItem{
// 		ID:     queue_id,
// 		Object: body,
// 	})
// }

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

// // field can be:
// // - status: processing, completed, in-queue, failed
// // - instruction: add_order, update_order, add_product, update_product, add_variant, update_variant
// // - type: product, order

// // GET /api/queue/filter?key=value
// func (dbconfig *DbConfig) FilterQueueItems(w http.ResponseWriter, r *http.Request, user database.User) {
// 	page, err := strconv.Atoi(r.URL.Query().Get("page"))
// 	if err != nil {
// 		page = 1
// 	}
// 	param_type := utils.ConfirmFilters(r.URL.Query().Get("type"))
// 	param_instruction := utils.ConfirmFilters(r.URL.Query().Get("instruction"))
// 	param_status := utils.ConfirmFilters(r.URL.Query().Get("status"))
// 	result, err := CompileQueueFilterSearch(dbconfig, r.Context(), page, param_type, param_status, param_instruction)
// 	if err != nil {
// 		RespondWithError(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}
// 	RespondWithJSON(w, http.StatusOK, result)
// }

// // GET /api/queue/processing
// func (dbconfig *DbConfig) QueueViewCurrentItem(w http.ResponseWriter, r *http.Request, user database.User) {
// 	page, err := strconv.Atoi(r.URL.Query().Get("page"))
// 	if err != nil {
// 		page = 1
// 	}
// 	queue_items, err := dbconfig.DB.GetQueueItemsByDate(r.Context(), database.GetQueueItemsByDateParams{
// 		Status: "in-queue",
// 		Limit:  10,
// 		Offset: int32((page - 1) * 10),
// 	})
// 	if err != nil {
// 		RespondWithError(w, http.StatusInternalServerError, err.Error())
// 	}
// 	RespondWithJSON(w, http.StatusOK, queue_items)
// }

// // GET /api/queue?page=1
// func (dbconfig *DbConfig) QueueViewNextItems(w http.ResponseWriter, r *http.Request, user database.User) {
// 	page, err := strconv.Atoi(r.URL.Query().Get("page"))
// 	if err != nil {
// 		page = 1
// 	}
// 	queue_items, err := dbconfig.DB.GetQueueItemsByDate(r.Context(), database.GetQueueItemsByDateParams{
// 		Status: "in-queue",
// 		Limit:  10,
// 		Offset: int32((page - 1) * 10),
// 	})
// 	if err != nil {
// 		RespondWithError(w, http.StatusInternalServerError, err.Error())
// 	}
// 	if queue_items == nil {
// 		RespondWithJSON(w, http.StatusOK, []string{})
// 		return
// 	} else {
// 		RespondWithJSON(w, http.StatusOK, queue_items)
// 	}
// }

// // GET /api/queue/view
// func (dbconfig *DbConfig) QueueView(
// 	w http.ResponseWriter,
// 	r *http.Request,
// 	user database.User) {
// 	response, err := dbconfig.DisplayQueueCount()
// 	if err != nil {
// 		RespondWithError(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}
// 	RespondWithJSON(w, http.StatusOK, response)
// }

// // DELETE /api/queue/{id}
// func (dbconfig *DbConfig) ClearQueueByID(
// 	w http.ResponseWriter,
// 	r *http.Request,
// 	user database.User) {
// 	id := chi.URLParam(r, "id")
// 	id_uuid, err := uuid.Parse(id)
// 	if err != nil {
// 		RespondWithError(w, http.StatusBadRequest, err.Error())
// 		return
// 	}
// 	err = dbconfig.DB.RemoveQueueItemByID(r.Context(), id_uuid)
// 	if err != nil {
// 		RespondWithError(w, http.StatusBadRequest, err.Error())
// 		return
// 	}
// 	RespondWithJSON(w, http.StatusOK, objects.ResponseString{
// 		Message: "success",
// 	})
// }

// // DELETE /api/queue?key=value
// func (dbconfig *DbConfig) ClearQueueByFilter(
// 	w http.ResponseWriter,
// 	r *http.Request,
// 	user database.User) {
// 	param_type := utils.ConfirmFilters(r.URL.Query().Get("type"))
// 	param_instruction := utils.ConfirmFilters(r.URL.Query().Get("instruction"))
// 	param_status := utils.ConfirmFilters(r.URL.Query().Get("status"))
// 	response, err := CompileRemoveQueueFilter(dbconfig, r.Context(), param_type, param_status, param_instruction)
// 	if err != nil {
// 		RespondWithError(w, http.StatusBadRequest, err.Error())
// 		return
// 	}
// 	RespondWithJSON(w, http.StatusOK, objects.ResponseString{
// 		Message: response,
// 	})
// }

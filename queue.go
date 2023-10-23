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
	"time"
	"utils"

	"github.com/google/uuid"
)

// POST /api/queue
func (dbconfig *DbConfig) QueuePush(w http.ResponseWriter, r *http.Request, user database.User) {
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
	// converts the data to rawJSON
	raw, err := json.Marshal(&body.Object)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	queue_id, err := dbconfig.DB.CreateQueueItem(r.Context(), database.CreateQueueItemParams{
		ID:          uuid.New(),
		Object:      raw,
		Type:        body.Type,
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

// POST /api/queue/worker
func (dbconfig *DbConfig) QueuePopAndProcess(w http.ResponseWriter, r *http.Request, user database.User) {
	// fetch the next item from the database table queue_items based on the time it was added into the queue
	queue_item, err := dbconfig.DB.GetNextQueueItem(r.Context())
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// remove the queue item from the database...
	err = dbconfig.DB.RemoveQueueItemByID(r.Context(), queue_item.ID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// change the status inside the database to Processing
	err = dbconfig.DB.UpdateQueueItem(r.Context(), database.UpdateQueueItemParams{
		Status:    "processing",
		UpdatedAt: time.Now().UTC(),
		ID:        queue_item.ID,
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// process the queue item

	// read the type, object_id, and instruction and run the respective function.
	// return the appropriate response
}

// GET /api/queue?position=0
func (dbconfig *DbConfig) QueueViewNextItem(w http.ResponseWriter, r *http.Request, user database.User) {
	// fetch a collection of the next queue items
	// queue items should be grouped with their count inside a map[string]int
	// map is returned as response
}

// GET /api/queue?page=1
func (dbconfig *DbConfig) QueueViewCurrent(w http.ResponseWriter, r *http.Request, user database.User) {
	// fetch a collection of the next queue items
	// queue items should be grouped with their count inside a map[string]int
	// map is returned as response
}

// add rules for processing queue items (like which instructions should be prioritized above the rest etc...)
// should be another function that sorts the queue, but needs to fetch all rows...

// Process the queue item
func ProcessQueueItem(dbconfig *DbConfig, queue_item database.QueueItem) error {
	if queue_item.Type == "product" {
		queue_object, err := DecodeQueueItemProduct(queue_item.Object)
		if err != nil {
			return err
		}
		product_id, err := uuid.Parse(queue_object.Shopify.ProductID)
		if err != nil {
			return errors.New("could not decode feed_id: " + queue_object.SystemProductID)
		}
		shopifyConfig := shopify.InitConfigShopify()
		product, err := CompileProductData(dbconfig, product_id, context.Background(), false)
		if err != nil {
			return err
		}
		if queue_item.Instruction == "add_product" {
			// get new instance of shopifyConfig for each new instruction
			dbconfig.PushProduct(&shopifyConfig, product)
		} else if queue_item.Instruction == "update_product" {
			dbconfig.PushProduct(&shopifyConfig, product)
		} else {
			return errors.New("invalid product instruction")
		}
	} else if queue_item.Type == "order" {
		if queue_item.Instruction == "add_order" {
			queue_object, err := DecodeQueueItemOrder(queue_item.Object)
			if err != nil {
				return err
			}
			dbconfig.AddOrder(queue_object)
		} else if queue_item.Instruction == "update_order" {
			// update the existing order
			// needs to check if its exists
			// should check here...
		} else {
			return errors.New("invalid order instruction")
		}
	} else if queue_item.Type == "product_variant" {
		shopifyConfig := shopify.InitConfigShopify()
		queue_object, err := DecodeQueueItemProduct(queue_item.Object)
		if err != nil {
			return err
		}
		variant_id, err := uuid.Parse(queue_object.SystemProductID)
		if err != nil {
			return errors.New("could not decode feed_id: " + queue_object.SystemVariantID)
		}
		variant, err := CompileVariantData(dbconfig, variant_id, context.Background())
		if err != nil {
			return err
		}
		if queue_item.Instruction == "add_variant" {
			dbconfig.PushVariant(
				&shopifyConfig,
				variant,
				queue_object.Shopify.ProductID,
				queue_object.Shopify.VariantID,
			)
		} else if queue_item.Instruction == "update_variant" {
			dbconfig.PushVariant(
				&shopifyConfig,
				variant,
				queue_object.Shopify.ProductID,
				queue_object.Shopify.VariantID,
			)
		} else {
			return errors.New("invalid product_variant instruction")
		}
	}
	return nil
}

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

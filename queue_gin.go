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
	"time"
	"utils"

	"github.com/google/uuid"
)

/*
Process a queue item inside the queue depending on it's instruction

Note that zsync_channel is not a processable queue instruction
*/
func ProcessQueueItem(dbconfig *DbConfig, queue_item database.QueueItem) error {
	if queue_item.QueueType == "product" {
		if queue_item.Instruction == "zsync_channel" {
			// do not process this instruction, just return nil
			return nil
		}
		queue_object, err := DecodeQueueItemProduct(queue_item.Object)
		if err != nil {
			return err
		}
		product_id, err := uuid.Parse(queue_object.SystemProductID)
		if err != nil {
			return errors.New("could not decode product_id '" + queue_object.SystemProductID + "'")
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
	} else if queue_item.QueueType == "product_variant" {
		shopifyConfig := shopify.InitConfigShopify()
		queue_object, err := DecodeQueueItemProduct(queue_item.Object)
		if err != nil {
			return err
		}
		variant_id, err := uuid.Parse(queue_object.SystemVariantID)
		if err != nil {
			return errors.New("could not decode variant_id '" + queue_object.SystemVariantID + "'")
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
		if queue_item.Instruction == "add_variant" || queue_item.Instruction == "update_variant" {
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
	} else if queue_item.QueueType == "order" {
		queue_object, err := DecodeQueueItemOrder(queue_item.Object)
		if err != nil {
			return err
		}
		if queue_item.Instruction == "add_order" {
			return dbconfig.AddOrder(queue_object)
		} else if queue_item.Instruction == "update_order" {
			return dbconfig.UpdateOrder(queue_object)
		} else {
			return errors.New("invalid order instruction")
		}
	}
	return errors.New("invalid queue item type")
}

// Helper function that displays count of different instructions
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

// Helper function that updates the queue item if there is an error
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

// Helper function used to post internal requests to the queue endpoint
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

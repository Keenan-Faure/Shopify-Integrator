package main

import (
	"errors"
	"integrator/internal/database"
	"net/http"
	"objects"
	"time"

	"github.com/google/uuid"
)

// POST /api/queue
func (dbconfig *DbConfig) QueuePush(w http.ResponseWriter, r *http.Request, user database.User) {
	// read request and determine if it's a product request or order request
	// create respective queue item struct and post it to the database with the respective time.
	// if any errors respond with http.badrequest
	// log the item to the back of the queue
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
	queue_id, err := dbconfig.DB.CreateQueueItem(r.Context(), database.CreateQueueItemParams{
		ID:          uuid.New(),
		ObjectID:    body.ObjectID,
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
		ID:       queue_id,
		ObjectID: body.ObjectID.String(),
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

// Process the queue item
func ProcessQueueItem(dbconfig *DbConfig, queue_item database.QueueItem) error {
	if queue_item.Type == "product" {
		if queue_item.Instruction == "add_product" {
			dbconfig.PushProduct()
		} else if queue_item.Instruction == "add_variant" {
			dbconfig.PushProduct()
		} else if queue_item.Instruction == "update_product" {
			dbconfig.PushProduct()
		} else {
			return errors.New("invalid instruction")
		}
	} else if queue_item.Type == "order" {
		if queue_item.Instruction == "add_order" {
			// add the order
		} else if queue_item.Instruction == "update_order" {
			// update the existing order
			// needs to check if its exists
		}
	} else if queue_item.Type == "product_variant" {
		if queue_item.Instruction == "add_variant" {
			dbconfig.PushVariant()
		} else if queue_item.Instruction == "update_variant" {
			dbconfig.PushVariant()
		} else {
			return errors.New("invalid instruction")
		}
	}
	return nil
}

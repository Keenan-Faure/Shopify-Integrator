package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"integrator/internal/database"
	"log"
	"net/http"
	"objects"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestSynchronize(t *testing.T) {
	dbconfig := setupDatabase("", "", "", false)
	dbUser := createDatabaseUser(&dbconfig)
	defer dbconfig.DB.RemoveUser(context.Background(), dbUser.ApiKey)

	httpmock.Activate()
	InitMockQueue()
	defer httpmock.DeactivateAndReset()

	// Test 1 - valid function params
	response := InitQueueMockRequests(
		MOCK_APP_API_URL+"/api/shopify/sync?api_key="+dbUser.ApiKey,
		http.MethodPost, map[string][]string{}, nil, &dbconfig,
	)
	assert.Equal(t, 200, response.StatusCode)
}

func TestGetQueueItemByID(t *testing.T) {
	/* Test 1 - Invalid authentication */
	dbconfig := setupDatabase("", "", "", false)
	router := setUpAPI(&dbconfig)
	w := Init(
		"/api/queue/{id}",
		http.MethodGet, map[string][]string{}, nil, &dbconfig, router,
	)
	assert.Equal(t, 401, w.Code)

	/* Test 2 - Invalid queue_item_id (malformed) */
	dbUser := createDatabaseUser(&dbconfig)
	defer dbconfig.DB.RemoveUser(context.Background(), dbUser.ApiKey)
	w = Init(
		"/api/queue/abctest123?api_key="+dbUser.ApiKey,
		http.MethodGet, map[string][]string{}, nil, &dbconfig, router,
	)
	assert.Equal(t, 400, w.Code)

	/* Test 3 - Invalid queue_item_id (404) */
	w = Init(
		"/api/queue/"+uuid.New().String()+"?api_key="+dbUser.ApiKey,
		http.MethodGet, map[string][]string{}, nil, &dbconfig, router,
	)
	assert.Equal(t, 400, w.Code)

	/* Test 4 - Valid request */
	queueUUID := CreateDatabaseQueueItem(&dbconfig, "")
	defer ClearQueueMockData(&dbconfig)
	w = Init(
		"/api/queue/"+queueUUID.String()+"?api_key="+dbUser.ApiKey,
		http.MethodGet, map[string][]string{}, nil, &dbconfig, router,
	)
	assert.Equal(t, 200, w.Code)
}

func TestQueuePush(t *testing.T) {
	/* Test 1 - Invalid authentication */
	dbconfig := setupDatabase("", "", "", false)
	router := setUpAPI(&dbconfig)
	w := Init(
		"/api/queue",
		http.MethodPost, map[string][]string{}, nil, &dbconfig, router,
	)

	assert.Equal(t, 401, w.Code)

	/* Test 2 - valid request */
	dbUser := createDatabaseUser(&dbconfig)
	defer dbconfig.DB.RemoveUser(context.Background(), dbUser.ApiKey)
	queueItemData := QueuePayload("queue_add_product.json")
	w = Init(
		"/api/queue?api_key="+dbUser.ApiKey,
		http.MethodPost, map[string][]string{}, queueItemData, &dbconfig, router,
	)

	assert.Equal(t, 201, w.Code)
	response := objects.ResponseQueueItem{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	dbconfig.DB.RemoveQueueItemByID(context.Background(), response.ID)
}

func TestFilterQueueItems(t *testing.T) {
	/* Test 1 - Invalid authentication */
	dbconfig := setupDatabase("", "", "", false)
	router := setUpAPI(&dbconfig)
	w := Init(
		"/api/queue/filter?key=",
		http.MethodGet, map[string][]string{}, nil, &dbconfig, router,
	)
	assert.Equal(t, 401, w.Code)

	/* Test 2 - Invalid queue_item_id (malformed) */
	dbUser := createDatabaseUser(&dbconfig)
	defer dbconfig.DB.RemoveUser(context.Background(), dbUser.ApiKey)
	w = Init(
		"/api/queue/filter?&api_key="+dbUser.ApiKey,
		http.MethodGet, map[string][]string{}, nil, &dbconfig, router,
	)
	assert.Equal(t, 200, w.Code)

	/* Test 3 - no results (404) */
	w = Init(
		"/api/queue/filter?type=&api_key="+dbUser.ApiKey,
		http.MethodGet, map[string][]string{}, nil, &dbconfig, router,
	)
	assert.Equal(t, 200, w.Code)

	/* Test 4 - Valid request */
	CreateDatabaseQueueItem(&dbconfig, "")
	defer ClearQueueMockData(&dbconfig)
	w = Init(
		"/api/queue/filter?type=product&api_key="+dbUser.ApiKey,
		http.MethodGet, map[string][]string{}, nil, &dbconfig, router,
	)
	assert.Equal(t, 200, w.Code)
}

func TestQueueViewCurrentItem(t *testing.T) {
	/* Test 1 - Invalid authentication */
	dbconfig := setupDatabase("", "", "", false)
	router := setUpAPI(&dbconfig)
	w := Init(
		"/api/queue/processing",
		http.MethodGet, map[string][]string{}, nil, &dbconfig, router,
	)
	assert.Equal(t, 401, w.Code)

	/* Test 2 - valid request */
	dbUser := createDatabaseUser(&dbconfig)
	defer dbconfig.DB.RemoveUser(context.Background(), dbUser.ApiKey)
	CreateDatabaseQueueItem(&dbconfig, "processing")
	defer ClearQueueMockData(&dbconfig)
	w = Init(
		"/api/queue/processing?&api_key="+dbUser.ApiKey,
		http.MethodGet, map[string][]string{}, nil, &dbconfig, router,
	)
	response := []database.QueueItem{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.NotEqual(t, len(response), 0)
	assert.Equal(t, 200, w.Code)
}

func TestQueueView(t *testing.T) {
	/* Test 1 - Invalid authentication */
	dbconfig := setupDatabase("", "", "", false)
	router := setUpAPI(&dbconfig)
	w := Init(
		"/api/queue",
		http.MethodGet, map[string][]string{}, nil, &dbconfig, router,
	)
	assert.Equal(t, 401, w.Code)

	/* Test 2 - Invalid queue_item_id (malformed) */
	CreateDatabaseQueueItem(&dbconfig, "")
	defer ClearQueueMockData(&dbconfig)
	dbUser := createDatabaseUser(&dbconfig)
	defer dbconfig.DB.RemoveUser(context.Background(), dbUser.ApiKey)
	w = Init(
		"/api/queue?&api_key="+dbUser.ApiKey,
		http.MethodGet, map[string][]string{}, nil, &dbconfig, router,
	)
	assert.Equal(t, 200, w.Code)
	response := []database.QueueItem{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.NotEqual(t, len(response), 0)
}

func TestClearQueueByID(t *testing.T) {
	/* Test 1 - Invalid authentication */
	dbconfig := setupDatabase("", "", "", false)
	router := setUpAPI(&dbconfig)
	w := Init(
		"/api/queue/id",
		http.MethodDelete, map[string][]string{}, nil, &dbconfig, router,
	)
	assert.Equal(t, 401, w.Code)

	/* Test 2 - Invalid queue_item_id (malformed) */
	dbUser := createDatabaseUser(&dbconfig)
	defer dbconfig.DB.RemoveUser(context.Background(), dbUser.ApiKey)
	w = Init(
		"/api/queue/abctest123?api_key="+dbUser.ApiKey,
		http.MethodDelete, map[string][]string{}, nil, &dbconfig, router,
	)
	assert.Equal(t, 400, w.Code)

	/* Test 3 - Invalid queue_item_id (404) */
	w = Init(
		"/api/queue/"+uuid.New().String()+"?api_key="+dbUser.ApiKey,
		http.MethodDelete, map[string][]string{}, nil, &dbconfig, router,
	)
	assert.Equal(t, 200, w.Code)

	/* Test 4 - Valid request */
	queueUUID := CreateDatabaseQueueItem(&dbconfig, "")
	defer ClearQueueMockData(&dbconfig)
	w = Init(
		"/api/queue/"+queueUUID.String()+"?api_key="+dbUser.ApiKey,
		http.MethodDelete, map[string][]string{}, nil, &dbconfig, router,
	)
	assert.Equal(t, 200, w.Code)
}

func TestClearQueueByFilter(t *testing.T) {
	/* Test 1 - Invalid authentication */
	dbconfig := setupDatabase("", "", "", false)
	router := setUpAPI(&dbconfig)
	w := Init(
		"/api/queue?key=",
		http.MethodDelete, map[string][]string{}, nil, &dbconfig, router,
	)
	assert.Equal(t, 401, w.Code)

	/* Test 2 - Invalid queue_item_id (malformed) */
	dbUser := createDatabaseUser(&dbconfig)
	defer dbconfig.DB.RemoveUser(context.Background(), dbUser.ApiKey)
	w = Init(
		"/api/queue?&api_key="+dbUser.ApiKey,
		http.MethodDelete, map[string][]string{}, nil, &dbconfig, router,
	)
	assert.Equal(t, 200, w.Code)

	/* Test 3 - no results (404) */
	w = Init(
		"/api/queue?type=&api_key="+dbUser.ApiKey,
		http.MethodDelete, map[string][]string{}, nil, &dbconfig, router,
	)
	assert.Equal(t, 200, w.Code)

	/* Test 4 - Valid request */
	CreateDatabaseQueueItem(&dbconfig, "")
	defer ClearQueueMockData(&dbconfig)
	w = Init(
		"/api/queue?type=product&api_key="+dbUser.ApiKey,
		http.MethodDelete, map[string][]string{}, nil, &dbconfig, router,
	)
	assert.Equal(t, 200, w.Code)
}

func TestQueuePopAndProcess(t *testing.T) {
	/* Not testing yet */
}

func TestProcessQueueItem(t *testing.T) {
	/* Not testing yet */
}

func TestDisplayQueueCount(t *testing.T) {
	dbconfig := setupDatabase("", "", "", false)
	// Test 1 - empty queue
	result, err := dbconfig.DisplayQueueCount()
	assert.Equal(t, err, nil)
	assert.Equal(t, result.AddProduct, 0)

	// Test 2 - 1 item in queue
	CreateDatabaseQueueItem(&dbconfig, "")
	defer ClearQueueMockData(&dbconfig)
	result, err = dbconfig.DisplayQueueCount()
	assert.Equal(t, err, nil)
	assert.Equal(t, result.AddProduct, 1)
}

func TestFailedQueueItem(t *testing.T) {
	dbconfig := setupDatabase("", "", "", false)
	queueUUID := CreateDatabaseQueueItem(&dbconfig, "")
	queueItem, _ := dbconfig.DB.GetQueueItemByID(context.Background(), queueUUID)
	defer ClearQueueMockData(&dbconfig)

	// Test 1 - valid params
	FailedQueueItem(&dbconfig, queueItem, errors.New("MOCK_ERROR_MESSAGE"))
	queueItemAfterUpdate, _ := dbconfig.DB.GetQueueItemByID(context.Background(), queueUUID)
	assert.Equal(t, queueItemAfterUpdate.Status, "failed")
}

func TestQueueHelper(t *testing.T) {
	/* Not testing yet */
}

func CreateDatabaseQueueItem(dbconfig *DbConfig, status string) uuid.UUID {
	payload := QueueItemPayload("test-case-valid-product-queue-item.json")
	if status != "" {
		payload.Status = status
	}
	queue_id, err := dbconfig.DB.CreateQueueItem(context.Background(), database.CreateQueueItemParams{
		ID:          GetMockQueueItemUUID(),
		Object:      payload.Object,
		QueueType:   payload.QueueType,
		Instruction: payload.Instruction,
		Status:      payload.Status,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	})
	if err != nil {
		log.Fatal(err)
	}
	return queue_id
}

func ClearQueueMockData(dbconfig *DbConfig) {
	dbconfig.DB.RemoveQueueItemByID(context.Background(), GetMockQueueItemUUID())
}

func GetMockQueueItemUUID() uuid.UUID {
	UUID, err := uuid.Parse(MOCK_QUEUE_ITEM_ID)
	if err != nil {
		return uuid.Nil
	}
	return UUID
}

func InitQueueMockRequests(
	requestURL,
	requestMethod string,
	additionalRequestHeaders map[string][]string,
	payload interface{},
	dbconfig *DbConfig,
) *http.Response {
	httpClient := http.Client{
		Timeout: time.Second * 10,
	}
	var buffer bytes.Buffer
	switch dataType := payload.(type) {
	case bytes.Buffer:
		buffer = dataType
	default:
		err := json.NewEncoder(&buffer).Encode(payload)
		if err != nil {
			log.Fatal(err)
		}
	}
	req, _ := http.NewRequest(requestMethod, requestURL, &buffer)
	for key, value := range additionalRequestHeaders {
		for _, sub_value := range value {
			req.Header.Add(key, sub_value)
		}
	}
	req.Header.Add("Content-Type", "application/json")
	res, err := httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	return res
}

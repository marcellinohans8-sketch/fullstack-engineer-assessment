package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"backend/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type taskResponse struct {
	Success bool `json:"success"`
	Data    struct {
		ID          uint   `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Status      string `json:"status"`
		Assignee    string `json:"assignee"`
	} `json:"data"`
}

type taskListResponse struct {
	Success bool `json:"success"`
	Data    []struct {
		ID          uint   `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Status      string `json:"status"`
		Assignee    string `json:"assignee"`
	} `json:"data"`
	Pagination struct {
		Page  int `json:"page"`
		Limit int `json:"limit"`
		Total int `json:"total"`
	} `json:"pagination"`
}

func TestUpdateTask(t *testing.T) {
	router := SetupRouter(t)
	created := createTask(t, router, map[string]interface{}{
		"title":       "Update Test",
		"description": "Before Update",
		"status":      "todo",
		"assignee":    "Hans",
	})

	update := map[string]interface{}{
		"title":       "Update Success",
		"description": "After Update",
		"status":      "done",
		"assignee":    "Nina",
	}

	w := performJSON(router, http.MethodPut, "/api/tasks/"+strconv.Itoa(int(created.Data.ID)), update)

	require.Equal(t, http.StatusOK, w.Code)

	var updated taskResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &updated))
	assert.True(t, updated.Success)
	assert.Equal(t, "Update Success", updated.Data.Title)
	assert.Equal(t, "After Update", updated.Data.Description)
	assert.Equal(t, "done", updated.Data.Status)
	assert.Equal(t, "Nina", updated.Data.Assignee)
}

func TestSearchTask(t *testing.T) {
	router := SetupRouter(t)

	createTask(t, router, map[string]interface{}{
		"title":       "Search Test",
		"description": "Testing search endpoint",
		"status":      "todo",
		"assignee":    "Hans",
	})
	createTask(t, router, map[string]interface{}{
		"title":       "Billing Work",
		"description": "Payment page",
		"status":      "done",
		"assignee":    "Nina",
	})

	w := performRequest(router, http.MethodGet, "/api/tasks?keyword=Search&status=todo&assignee=Hans&page=1&limit=5&sort=title%20asc")

	require.Equal(t, http.StatusOK, w.Code)

	var response taskListResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.True(t, response.Success)
	assert.Equal(t, 1, response.Pagination.Total)
	require.Len(t, response.Data, 1)
	assert.Equal(t, "Search Test", response.Data[0].Title)
}

func TestCacheInvalidation(t *testing.T) {
	router := SetupRouter(t)

	created := createTask(t, router, map[string]interface{}{
		"title":       "Redis Test",
		"description": "Cache Test",
		"status":      "todo",
		"assignee":    "Hans",
	})

	w := performRequest(router, http.MethodGet, "/api/tasks?page=1&limit=10")
	require.Equal(t, http.StatusOK, w.Code)
	assertRedisKeys(t, 1)

	createTask(t, router, map[string]interface{}{
		"title":       "Redis Create Invalidates",
		"description": "Cache should be cleared after create",
		"status":      "todo",
		"assignee":    "Nina",
	})
	assertRedisKeys(t, 0)

	w = performJSON(router, http.MethodPut, "/api/tasks/"+strconv.Itoa(int(created.Data.ID)), map[string]interface{}{
		"title":       "Redis Test Updated",
		"description": "Cache Updated",
		"status":      "done",
		"assignee":    "Hans",
	})
	require.Equal(t, http.StatusOK, w.Code)
	assertRedisKeys(t, 0)

	w = performRequest(router, http.MethodGet, "/api/tasks?page=1&limit=10")
	require.Equal(t, http.StatusOK, w.Code)
	assertRedisKeys(t, 1)

	w = performRequest(router, http.MethodDelete, "/api/tasks/"+strconv.Itoa(int(created.Data.ID)))
	require.Equal(t, http.StatusOK, w.Code)
	assertRedisKeys(t, 0)
}

func TestDuplicateTitleReturnsConflict(t *testing.T) {
	router := SetupRouter(t)

	payload := map[string]interface{}{
		"title":       "Duplicate Title",
		"description": "First",
		"status":      "todo",
		"assignee":    "Hans",
	}
	createTask(t, router, payload)

	w := performJSON(router, http.MethodPost, "/api/tasks", payload)

	assert.Equal(t, http.StatusConflict, w.Code)
	assert.Contains(t, w.Body.String(), "Task title already exists")
}

func TestSoftDeletedTasksAreHidden(t *testing.T) {
	router := SetupRouter(t)

	created := createTask(t, router, map[string]interface{}{
		"title":       "Deleted Task",
		"description": "Should not show",
		"status":      "todo",
		"assignee":    "Hans",
	})

	w := performRequest(router, http.MethodDelete, "/api/tasks/"+strconv.Itoa(int(created.Data.ID)))
	require.Equal(t, http.StatusOK, w.Code)

	w = performRequest(router, http.MethodGet, "/api/tasks?keyword=Deleted")
	require.Equal(t, http.StatusOK, w.Code)

	var response taskListResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Equal(t, 0, response.Pagination.Total)
	assert.Empty(t, response.Data)
}

func createTask(t *testing.T, router http.Handler, payload map[string]interface{}) taskResponse {
	t.Helper()

	w := performJSON(router, http.MethodPost, "/api/tasks", payload)
	require.Equal(t, http.StatusCreated, w.Code, w.Body.String())

	var response taskResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	require.True(t, response.Success)

	return response
}

func performJSON(router http.Handler, method string, path string, payload map[string]interface{}) *httptest.ResponseRecorder {
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(method, path, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	return w
}

func performRequest(router http.Handler, method string, path string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	return w
}

func assertRedisKeys(t *testing.T, expected int) {
	t.Helper()

	keys, err := config.RedisClient.Keys(config.Ctx, "tasks:*").Result()
	require.NoError(t, err)
	assert.Len(t, keys, expected)
}

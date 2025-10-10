package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jennaborowy/fullstack-Go-Docker/handlers"
	"github.com/jennaborowy/fullstack-Go-Docker/mocks"
	"github.com/jennaborowy/fullstack-Go-Docker/models"
	"github.com/jennaborowy/fullstack-Go-Docker/repository"
	"go.uber.org/mock/gomock"
)

var (
	validItem = models.NewItem("Item 1", time.Date(2025, 10, 7, 0, 0, 0, 0, time.UTC), "test description uno", 1)

	multipleItems = []models.Item{
		*models.NewItem("Item 2", time.Date(2025, 10, 8, 0, 0, 0, 0, time.UTC), "test description dos", 1),
		*models.NewItem("Item 3", time.Date(2025, 10, 10, 0, 0, 0, 0, time.UTC), "third test description", 1),
		*models.NewItem("Item 4", time.Date(2025, 10, 12, 0, 0, 0, 0, time.UTC), "fourth test description", 1),
	}

	validList = models.NewList("test list", multipleItems)
)

func TestGetItem(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(*mocks.MockItemRepositoryInterface)
		requestedID    string
		expectedStatus int
		checkResponse  func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "successful get",
			setupMock: func(m *mocks.MockItemRepositoryInterface) {
				m.EXPECT().
					GetByID(1).
					Return(validItem, nil).
					Times(1)
			},
			requestedID:    "1",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response models.Item
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}

				if response.Title != validItem.Title {
					t.Errorf("expected title '%s', got '%s'", validItem.Title, response.Title)
				}
			},
		},
		{
			name: "item does not exist",
			setupMock: func(m *mocks.MockItemRepositoryInterface) {
				m.EXPECT().
					GetByID(5).
					Return(nil, repository.ErrNotFound).
					Times(1)
			},
			requestedID:    "5",
			expectedStatus: http.StatusNotFound,
			checkResponse:  nil,
		},
		{
			name:           "invalid id format",
			setupMock:      func(m *mocks.MockItemRepositoryInterface) {},
			requestedID:    "invalid",
			expectedStatus: http.StatusBadRequest,
			checkResponse:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)

			ctrl := gomock.NewController(t)

			repo := mocks.NewMockItemRepositoryInterface(ctrl)
			handler := handlers.NewItemHandler(repo)

			tt.setupMock(repo)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Params = gin.Params{
				{Key: "id", Value: tt.requestedID},
			}

			c.Request = httptest.NewRequest(http.MethodGet, "/items/"+tt.requestedID, nil)

			handler.GetItem(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d. Response: %s",
					tt.expectedStatus, w.Code, w.Body.String())
			}

			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
		},
		)
	}

}

func TestGetItems(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(*mocks.MockItemRepositoryInterface)
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "successfully fetch all items",
			setupMock: func(m *mocks.MockItemRepositoryInterface) {
				m.EXPECT().
					GetAll().
					Return(multipleItems, nil).
					Times(1)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response []models.Item
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}

				if len(response) != 3 {
					t.Errorf("expected 3 items, got %d", len(response))
				}

				if response[0].Title != "Item 2" {
					t.Errorf("expected first item title 'Item 2', got '%s'", response[0].Title)
				}
			},
		},
		{
			name: "repository error",
			setupMock: func(m *mocks.MockItemRepositoryInterface) {
				m.EXPECT().
					GetAll().
					Return(nil, errors.New("database error")).
					Times(1)
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse:  nil,
		},
		{
			name: "empty items list",
			setupMock: func(m *mocks.MockItemRepositoryInterface) {
				m.EXPECT().
					GetAll().
					Return([]models.Item{}, nil).
					Times(1)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response []models.Item
				json.Unmarshal(w.Body.Bytes(), &response)

				if len(response) != 0 {
					t.Errorf("expected 0 items, got %d", len(response))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)

			// Fresh controller per iteration
			ctrl := gomock.NewController(t)

			repo := mocks.NewMockItemRepositoryInterface(ctrl)
			handler := handlers.NewItemHandler(repo)

			// Setup mock expectations
			tt.setupMock(repo)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = httptest.NewRequest(http.MethodGet, "/items", nil)

			handler.GetItems(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d. Response: %s",
					tt.expectedStatus, w.Code, w.Body.String())
			}

			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
		})
	}

}

func TestCreateItem(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		setupMock      func(*mocks.MockItemRepositoryInterface)
		expectedStatus int
	}{
		{
			name: "successful creation",
			requestBody: map[string]interface{}{
				"title":     "test",
				"content":   "hello this is a test description",
				"item_date": "2025-10-08",
				"list_id":   1,
			},
			setupMock: func(m *mocks.MockItemRepositoryInterface) {
				m.EXPECT().
					CreateItem("test", gomock.Any(), "hello this is a test description", 1).
					Return(&models.Item{ID: 1, Title: "test"}, nil).
					Times(1)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "repository error",
			requestBody: map[string]interface{}{
				"title":     "test",
				"content":   "hello this is a test description",
				"item_date": "2025-10-08",
				"list_id":   1,
			},
			setupMock: func(m *mocks.MockItemRepositoryInterface) {
				m.EXPECT().
					CreateItem(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errors.New("database error")).
					Times(1)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "invalid date format",
			requestBody: map[string]interface{}{
				"title":     "test",
				"content":   "hello this is a test description",
				"item_date": "invalid-date",
				"list_id":   1,
			},
			setupMock:      func(m *mocks.MockItemRepositoryInterface) {}, // No mock needed
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)

			// Fresh controller per iteration
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mocks.NewMockItemRepositoryInterface(ctrl)
			handler := handlers.NewItemHandler(repo)

			// Setup mock expectations
			tt.setupMock(repo)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			body, _ := json.Marshal(tt.requestBody)
			c.Request = httptest.NewRequest(http.MethodPost, "/items", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			handler.CreateItem(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d. Response: %s",
					tt.expectedStatus, w.Code, w.Body.String())
			}
		})
	}
}

func TestDeleteItem(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(m *mocks.MockItemRepositoryInterface)
		id             string
		expectedStatus int
	}{
		{
			name: "successful delete (item exits)",
			setupMock: func(m *mocks.MockItemRepositoryInterface) {
				m.EXPECT().
					DeleteItemByID(1).
					Return(nil).
					Times(1)
			},
			id:             "1",
			expectedStatus: http.StatusNoContent,
		},
		{
			name: "repository error",
			setupMock: func(m *mocks.MockItemRepositoryInterface) {
				m.EXPECT().
					DeleteItemByID(20).
					Return(errors.New("database error")).
					Times(1)
			},
			id:             "20",
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)

			ctrl := gomock.NewController(t)

			repo := mocks.NewMockItemRepositoryInterface(ctrl)
			handler := handlers.NewItemHandler(repo)

			tt.setupMock(repo)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Params = gin.Params{
				{Key: "id", Value: tt.id},
			}

			http.NewRequest("DELETE", "/items/"+tt.id, nil)

			handler.DeleteItem(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d. Response: %s",
					tt.expectedStatus, w.Code, w.Body.String())
			}

		})
	}
}

func TestUpdateItem(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(m *mocks.MockItemRepositoryInterface)
		id             string
		requestBody    map[string]interface{}
		expectedStatus int
		checkResponse  func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "succesfully update item",
			setupMock: func(m *mocks.MockItemRepositoryInterface) {
				m.EXPECT().
					GetByID(1).
					Return(validItem, nil).
					Times(1)

				m.EXPECT().
					UpdateItem(1, "new title", gomock.Any(), "new content").
					Return(nil).
					Times(1)
			},
			id: "1",
			requestBody: map[string]interface{}{
				"title":     "new title",
				"content":   "new content",
				"item_date": "2025-10-23",
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response models.Item
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}

				if response.Title != "new title" {
					t.Errorf("expected first item title 'new title', got '%s'", response.Title)
				}
			},
		},
		{
			name: "repository error",
			setupMock: func(m *mocks.MockItemRepositoryInterface) {
				m.EXPECT().
					GetByID(1).
					Return(validItem, nil).
					Times(1)

				m.EXPECT().
					UpdateItem(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("database error")).
					Times(1)
			},
			id: "1",
			requestBody: map[string]interface{}{
				"title":     "new title",
				"content":   "new content",
				"item_date": "2025-10-23",
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse:  nil,
		},
		{
			name: "item not found on GetByID",
			setupMock: func(m *mocks.MockItemRepositoryInterface) {
				// GetByID fails - UpdateItem is never called
				m.EXPECT().
					GetByID(999).
					Return(nil, repository.ErrNotFound).
					Times(1)
			},
			id: "999",
			requestBody: map[string]interface{}{
				"title":     "new title",
				"item_date": "2025-10-23",
				"content":   "new content",
			},
			expectedStatus: http.StatusNotFound,
			checkResponse:  nil,
		},
		{
			name: "invalid ID format",
			setupMock: func(m *mocks.MockItemRepositoryInterface) {
				// No expectations - should fail before reaching repo
			},
			id: "invalid",
			requestBody: map[string]interface{}{
				"title":     "new title",
				"item_date": "2025-10-23",
				"content":   "new content",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse:  nil,
		},
		{
			name: "invalid date format",
			setupMock: func(m *mocks.MockItemRepositoryInterface) {
				// No expectations - should fail before reaching repo
			},
			id: "1",
			requestBody: map[string]interface{}{
				"title":     "new title",
				"item_date": "invalid-date",
				"content":   "new content",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse:  nil,
		},
		{
			name: "missing required fields",
			setupMock: func(m *mocks.MockItemRepositoryInterface) {
				// No expectations - should fail at binding
			},
			id: "1",
			requestBody: map[string]interface{}{
				"title": "new title",
				// Missing item_date and content
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)

			ctrl := gomock.NewController(t)

			repo := mocks.NewMockItemRepositoryInterface(ctrl)
			handler := handlers.NewItemHandler(repo)

			tt.setupMock(repo)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Params = gin.Params{
				{Key: "id", Value: tt.id},
			}

			body, _ := json.Marshal(tt.requestBody)
			c.Request = httptest.NewRequest(http.MethodPut, "/items/"+tt.id, bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			handler.UpdateItem(c)

			if tt.expectedStatus != w.Code {
				t.Errorf("expected status %d, got %d. Response: %s",
					tt.expectedStatus, w.Code, w.Body.String())
			}

			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
		})
	}
}

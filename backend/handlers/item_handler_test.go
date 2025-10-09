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

				if response[0].Title != "Item 1" {
					t.Errorf("expected first item title 'Item 1', got '%s'", response[0].Title)
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

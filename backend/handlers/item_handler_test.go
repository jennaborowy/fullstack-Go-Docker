package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jennaborowy/fullstack-Go-Docker/handlers"
	"github.com/jennaborowy/fullstack-Go-Docker/mocks"
	"github.com/jennaborowy/fullstack-Go-Docker/models"
	"go.uber.org/mock/gomock"
)

// var (
// 	item1 = &models.Item{
// 		ID:        1,
// 		Title:     "test",
// 		Date:      time.Date(2025, 10, 8, 0, 0, 0, 0, time.Local),
// 		Content:   "hello this is a test description",
// 		ListID:    1,
// 		CreatedAt: time.Now(),
// 		UpdatedAt: time.Now(),
// 	}
// )

// var (
//     validItem = &models.Item{
//         ID:     "item-1",
//         Title:   "Test Item 1",
// 		Date:
// 		Content:
//         ListID: "list-123",
// 		CreatedAt:
// 		UpdatedAt:
//     }

//     validList = &models.List{
//         ID:   "list-456",
//         Name: "Test List",
// 		Date:
//     }

//     multipleItems = []*models.Item{
//         {ID: "1", Title: "Item 1", Date: , ListID: "list-456", CreatedAt:, UpdatedAt: },
//         {ID: "2", Title: "Item 2", Date: , ListID: "list-456", CreatedAt: , UpdatedAt: },
//         {ID: "3", Title: "Item 3", ListID: "list-456"},
//     }
// )

// func TestGetItems(t *testing.T) {
// 	ctrl := gomock.NewController(t)

// 	mock := mocks.NewMockItemRepositoryInterface(ctrl)

// 	// params in Return should match the return signature of method being mocked

// 	mock.EXPECT().
// 		GetAll().
// 		Return(nil)

// 	handler := handlers.NewItemHandler(mock)

// }

// func TestCreateItem(t *testing.T) {
// 	ctrl := gomock.NewController(t)

// 	repo := mocks.NewMockItemRepositoryInterface(ctrl)
// 	handler := handlers.NewItemHandler(repo)

// 	t.Run("successful creation", func(t *testing.T) {
// 		gin.SetMode(gin.TestMode)

// 		testTime := time.Date(2025, 10, 8, 12, 0, 0, 0, time.UTC)

// 		date := time.Date(2025, 10, 8, 0, 0, 0, 0, time.UTC)

// 		requestBody := map[string]interface{}{
// 			"title":     "test",
// 			"content":   "hello this is a test description",
// 			"item_date": "2025-10-08", // Send as string in YYYY-MM-DD format
// 			"list_id":   1,
// 		}

// 		expectedItem := &models.Item{
// 			ID:        1,
// 			Title:     "test",
// 			Content:   "hello this is a test description",
// 			Date:      date,
// 			ListID:    1,
// 			CreatedAt: testTime, // Use fixed time
// 			UpdatedAt: testTime, // Use fixed time
// 		}

// 		repo.EXPECT().
// 			CreateItem("test", date, "hello this is a test description", 1).
// 			Return(expectedItem, nil).
// 			Times(1)

// 		// if not using gin, use this to call handler.CreateItem(w, req)
// 		// body, _ := json.Marshal(item)
// 		// req := httptest.NewRequest(http.MethodPost, "/items", bytes.NewBuffer(body))
// 		// w := httptest.NewRecorder()

// 		// create context with request header
// 		w := httptest.NewRecorder()
// 		c, _ := gin.CreateTestContext(w)

// 		// prepare request body
// 		body, _ := json.Marshal(requestBody)
// 		c.Request = httptest.NewRequest(http.MethodPost, "/items", bytes.NewBuffer(body))
// 		c.Request.Header.Set("Content-Type", "application/json")

// 		// call handler
// 		handler.CreateItem(c)

// 		// assert response
// 		if w.Code != http.StatusCreated {
// 			t.Errorf("expected status %d, got %d. Response: %s", http.StatusCreated, w.Code, w.Body.String())
// 		}

// 		var response models.Item
// 		json.Unmarshal(w.Body.Bytes(), &response)

// 		if response.Title != "test" {
// 			t.Errorf("expected title 'test', got '%s'", response.Title)
// 		}
// 		if response.ListID != 1 {
// 			t.Errorf("expected list_id 1, got %d", response.ListID)
// 		}
// 	})

// 	t.Run("repository error", func(t *testing.T) {
// 		gin.SetMode(gin.TestMode)

// 		requestBody := map[string]interface{}{
// 			"title":     "test",
// 			"content":   "hello this is a test description",
// 			"item_date": "2025-10-08",
// 			"list_id":   1,
// 		}

// 		repo.EXPECT().
// 			CreateItem("test", gomock.Any(), "hello this is a test description", 1).
// 			Return(nil, errors.New("database connection failed")).
// 			Times(1)

// 		w := httptest.NewRecorder()
// 		c, _ := gin.CreateTestContext(w)

// 		body, _ := json.Marshal(requestBody)
// 		c.Request = httptest.NewRequest(http.MethodPost, "/items", bytes.NewBuffer(body))
// 		c.Request.Header.Set("Content-Type", "applications/json")

// 		handler.CreateItem(c)

// 		if w.Code != http.StatusInternalServerError {
// 			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
// 		}

// 	})

// 	t.Run("invalid date format", func(t *testing.T) {
// 		gin.SetMode(gin.TestMode)

// 		// intentionally wrong date format
// 		requestBody := map[string]interface{}{
// 			"title":     "test",
// 			"content":   "hello this is a test description",
// 			"item_date": "2025-10-08T00:00:00Z",
// 			"list_id":   1,
// 		}

// 		// No repo mock needed - binding/parsing fails first
// 		// repo.EXPECT() is NOT called

// 		w := httptest.NewRecorder()
// 		c, _ := gin.CreateTestContext(w)

// 		body, _ := json.Marshal(requestBody)
// 		c.Request = httptest.NewRequest(http.MethodPost, "/items", bytes.NewBuffer(body))
// 		c.Request.Header.Set("Content-Type", "application/json")

// 		handler.CreateItem(c)

// 		if w.Code != http.StatusBadRequest {
// 			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
// 		}
// 	})

// 	t.Run("missing required fields", func(t *testing.T) {
// 		gin.SetMode(gin.TestMode)

// 		// Missing title - should return 400
// 		requestBody := map[string]interface{}{
// 			"content":   "hello this is a test description",
// 			"item_date": "2025-10-08",
// 			"list_id":   1,
// 		}

// 		w := httptest.NewRecorder()
// 		c, _ := gin.CreateTestContext(w)

// 		body, _ := json.Marshal(requestBody)
// 		c.Request = httptest.NewRequest(http.MethodPost, "/items", bytes.NewBuffer(body))
// 		c.Request.Header.Set("Content-Type", "application/json")

// 		handler.CreateItem(c)

// 		if w.Code != http.StatusBadRequest {
// 			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
// 		}
// 	})
// }

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

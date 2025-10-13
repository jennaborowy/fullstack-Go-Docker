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
	"github.com/jennaborowy/fullstack-Go-Docker/repository"
	"go.uber.org/mock/gomock"
)

// can use validItem, validList, and multipleItems from item_handler_test, since in the same package
var (
	emptyList          = models.NewList("no items", []models.Item{})
	multipleLists      = []models.List{*emptyList, *validList}
	multipleEmptyLists = []models.List{*emptyList, *emptyList, *emptyList}
	noLists            = []models.List{}
)

// todo: get list
func TestGetList(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(m *mocks.MockListRepositoryInterface)
		id             string
		expectedStatus int
		checkResponse  func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "successfully get list",
			setupMock: func(m *mocks.MockListRepositoryInterface) {
				m.EXPECT().
					GetList(1).
					Return(validList, nil).
					Times(1)
			},
			id:             "1",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp models.List
				err := json.Unmarshal(w.Body.Bytes(), &resp)

				if err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}

				if resp.Title != validList.Title {
					t.Errorf("expected title '%s', got '%s'", validList.Title, resp.Title)
				}
			},
		},
		{
			name: "list does not exist",
			setupMock: func(m *mocks.MockListRepositoryInterface) {
				m.EXPECT().
					GetList(5).
					Return(nil, repository.ErrNotFound).
					Times(1)
			},
			id:             "5",
			expectedStatus: http.StatusNotFound,
			checkResponse:  nil,
		},
		{
			name:           "invalid id format",
			setupMock:      func(m *mocks.MockListRepositoryInterface) {},
			id:             "invalid",
			expectedStatus: http.StatusBadRequest,
			checkResponse:  nil,
		},
		{
			name: "repository error",
			setupMock: func(m *mocks.MockListRepositoryInterface) {
				m.EXPECT().
					GetList(gomock.Any()).
					Return(nil, errors.New("database error")).
					Times(1)
			},
			id:             "1",
			expectedStatus: http.StatusInternalServerError,
			checkResponse:  nil,
		},
		{
			name: "no items in list, should return",
			setupMock: func(m *mocks.MockListRepositoryInterface) {
				m.EXPECT().
					GetList(3).
					Return(emptyList, nil).
					Times(1)
			},
			id:             "3",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp models.List
				err := json.Unmarshal(w.Body.Bytes(), &resp)

				if err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}

				if resp.Title != emptyList.Title {
					t.Errorf("expected title '%s', got '%s'", emptyList.Title, resp.Title)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)

			ctrl := gomock.NewController(t)

			repo := mocks.NewMockListRepositoryInterface(ctrl)
			handler := handlers.NewListHandler(repo)

			tt.setupMock(repo)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Params = gin.Params{
				{Key: "id", Value: tt.id},
			}

			c.Request = httptest.NewRequest("GET", "/lists/"+tt.id, nil)

			handler.GetList(c)

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

// todo: get lists
func TestGetLists(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(m *mocks.MockListRepositoryInterface)
		expectedStatus int
		checkResponse  func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "successfully get all lists",
			setupMock: func(m *mocks.MockListRepositoryInterface) {
				m.EXPECT().
					GetAllLists().
					Return(multipleLists, nil).
					Times(1)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp []models.List
				if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
			},
		},
		{
			name: "repository error",
			setupMock: func(m *mocks.MockListRepositoryInterface) {
				m.EXPECT().
					GetAllLists().
					Return(nil, errors.New("database error")).
					Times(1)
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse:  nil,
		},
		{
			name: "multiple empty lists",
			setupMock: func(m *mocks.MockListRepositoryInterface) {
				m.EXPECT().
					GetAllLists().
					Return(multipleEmptyLists, nil).
					Times(1)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp []models.List
				if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
			},
		},
		{
			name: "no lists have been created",
			setupMock: func(m *mocks.MockListRepositoryInterface) {
				m.EXPECT().
					GetAllLists().
					Return(noLists, nil).
					Times(1)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp []models.List
				if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)

			ctrl := gomock.NewController(t)

			repo := mocks.NewMockListRepositoryInterface(ctrl)
			handler := handlers.NewListHandler(repo)

			tt.setupMock(repo)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = httptest.NewRequest(http.MethodGet, "/lists/", nil)

			handler.GetLists(c)

			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}

		})
	}
}

func TestCreateList(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(m *mocks.MockListRepositoryInterface)
		requestBody    map[string]interface{}
		expectedStatus int
		checkResponse  func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "successfully create list",
			setupMock: func(m *mocks.MockListRepositoryInterface) {
				m.EXPECT().
					CreateList(gomock.Any()).
					DoAndReturn(func(title string) (*models.List, error) {
						return &models.List{
							ID:    1,
							Title: "My New List",
							Items: []models.Item{},
						}, nil
					}).
					Times(1)
			},
			requestBody: map[string]interface{}{
				"title": "My New List",
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response models.List
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				if response.Title != "My New List" {
					t.Errorf("expected title 'My New List', got '%s'", response.Title)
				}
				if response.ID == 0 {
					t.Error("expected ID to be set")
				}
			},
		},
		{
			name: "repository error",
			setupMock: func(m *mocks.MockListRepositoryInterface) {
				m.EXPECT().
					CreateList(gomock.Any()).
					Return(nil, errors.New("database error")).
					Times(1)
			},
			requestBody: map[string]interface{}{
				"title": "My New List",
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse:  nil,
		},
		{
			name: "empty title",
			setupMock: func(m *mocks.MockListRepositoryInterface) {
				m.EXPECT().
					CreateList(gomock.Any()).
					DoAndReturn(func(title string) (*models.List, error) {
						return &models.List{
							ID:    1,
							Title: "",
							Items: []models.Item{},
						}, nil
					}).
					Times(1)
			},
			requestBody: map[string]interface{}{
				"title": "",
			},
			expectedStatus: http.StatusCreated,
			checkResponse:  nil,
		},
		{
			name:           "invalid JSON",
			setupMock:      func(m *mocks.MockListRepositoryInterface) {},
			requestBody:    nil, // Will send malformed JSON
			expectedStatus: http.StatusBadRequest,
			checkResponse:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)

			ctrl := gomock.NewController(t)

			repo := mocks.NewMockListRepositoryInterface(ctrl)
			handler := handlers.NewListHandler(repo)

			tt.setupMock(repo)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			var body []byte
			var err error
			if tt.requestBody != nil {
				body, err = json.Marshal(tt.requestBody)
				if err != nil {
					t.Fatalf("failed to marshal request: %v", err)
				}
			} else {
				body = []byte("{invalid json}")
			}

			c.Request = httptest.NewRequest(http.MethodPost, "/lists", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			handler.CreateList(c)

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

func TestUpdateList(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(m *mocks.MockListRepositoryInterface)
		id             string
		requestBody    map[string]interface{}
		expectedStatus int
		checkResponse  func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "successfully update list",
			setupMock: func(m *mocks.MockListRepositoryInterface) {
				m.EXPECT().
					UpdateTitle(1, "Updated Title").
					Return(&models.List{
						ID:    1,
						Title: "Updated Title",
						Items: []models.Item{},
					}, nil).
					Times(1)
			},
			id: "1",
			requestBody: map[string]interface{}{
				"title": "Updated Title",
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response models.List
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				if response.Title != "Updated Title" {
					t.Errorf("expected title 'Updated Title', got '%s'", response.Title)
				}
			},
		},
		{
			name: "repository error on UpdateList",
			setupMock: func(m *mocks.MockListRepositoryInterface) {
				m.EXPECT().
					UpdateTitle(1, "Updated Title").
					Return(nil, errors.New("database error")).
					Times(1)
			},
			id: "1",
			requestBody: map[string]interface{}{
				"title": "Updated Title",
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse:  nil,
		},
		{
			name:      "invalid ID format",
			setupMock: func(m *mocks.MockListRepositoryInterface) {},
			id:        "invalid",
			requestBody: map[string]interface{}{
				"title": "Updated Title",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse:  nil,
		},
		{
			name: "empty title",
			setupMock: func(m *mocks.MockListRepositoryInterface) {
				m.EXPECT().
					UpdateTitle(1, "").
					Return(&models.List{
						ID:    1,
						Title: "",
						Items: []models.Item{},
					}, nil).
					Times(1)
			},
			id: "1",
			requestBody: map[string]interface{}{
				"title": "",
			},
			expectedStatus: http.StatusOK,
			checkResponse:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)

			ctrl := gomock.NewController(t)

			repo := mocks.NewMockListRepositoryInterface(ctrl)
			handler := handlers.NewListHandler(repo)

			tt.setupMock(repo)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Params = gin.Params{
				{Key: "id", Value: tt.id},
			}

			body, _ := json.Marshal(tt.requestBody)
			c.Request = httptest.NewRequest(http.MethodPut, "/lists/"+tt.id, bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			handler.UpdateListTitle(c)

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

func TestDeleteList(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(m *mocks.MockListRepositoryInterface)
		id             string
		expectedStatus int
	}{
		{
			name: "successfully delete list",
			setupMock: func(m *mocks.MockListRepositoryInterface) {
				m.EXPECT().
					DeleteList(1).
					Return(nil).
					Times(1)
			},
			id:             "1",
			expectedStatus: http.StatusNoContent,
		},
		// {
		//     name: "list not found",
		//     setupMock: func(m *mocks.MockListRepositoryInterface) {
		//         m.EXPECT().
		//             DeleteList(999).
		//             Return(repository.ErrNotFound).
		//             Times(1)
		//     },
		//     id:             "999",
		//     expectedStatus: http.StatusNotFound,
		//     checkResponse:  nil,
		// },
		{
			name: "repository error",
			setupMock: func(m *mocks.MockListRepositoryInterface) {
				m.EXPECT().
					DeleteList(1).
					Return(errors.New("database error")).
					Times(1)
			},
			id:             "1",
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "invalid ID format",
			setupMock:      func(m *mocks.MockListRepositoryInterface) {},
			id:             "invalid",
			expectedStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)

			ctrl := gomock.NewController(t)

			repo := mocks.NewMockListRepositoryInterface(ctrl)
			handler := handlers.NewListHandler(repo)

			tt.setupMock(repo)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Params = gin.Params{
				{Key: "id", Value: tt.id},
			}

			c.Request = httptest.NewRequest(http.MethodDelete, "/lists/"+tt.id, nil)

			handler.DeleteList(c)

			if tt.expectedStatus != w.Code {
				t.Errorf("expected status %d, got %d. Response: %s",
					tt.expectedStatus, w.Code, w.Body.String())
			}
		})
	}
}

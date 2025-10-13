package handlers_test

import (
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

package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/jennaborowy/fullstack-Go-Docker/handlers"
	"github.com/jennaborowy/fullstack-Go-Docker/mocks"
	"go.uber.org/mock/gomock"
)

func TestGetItems(t *testing.T) {
	ctrl := gomock.NewController(t)

	mock := mocks.NewMockItemRepositoryInterface(ctrl)

	mock.EXPECT().
		GetAll().
		Return()

	handler := handlers.NewItemHandler(mock)

	router := gin.Default()
	router.GET("/api/items", handler.GetItems)

	req, _ := http.NewRequest("GET", "/api/items", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

}

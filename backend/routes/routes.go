package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/jennaborowy/fullstack-Go-Docker/handlers"
	"github.com/jennaborowy/fullstack-Go-Docker/middleware"
	"github.com/jennaborowy/fullstack-Go-Docker/repository"
)

func SetupRoutes(db *sql.DB) *gin.Engine {
	// create a new gin engine
	router := gin.Default()

	router.Use(middleware.CORSMiddleware())

	// create repositories and handlers
	itemRepo := repository.NewItemRepository(db)
	itemHandler := handlers.NewItemHandler(itemRepo)

	listRepo := repository.NewListRepository(db)
	listHandler := handlers.NewListHandler(listRepo)

	// define routes that can be used
	router.GET("/api/items", itemHandler.GetItems)
	router.GET("/api/items/:id", itemHandler.GetItem)
	// router.GET("/api/items/:list_id", itemHandler.GetItemFromList)
	router.POST("/api/items", itemHandler.CreateItem)
	router.DELETE("/api/items/:id", itemHandler.DeleteItem)
	router.PUT("/api/items/:id", itemHandler.UpdateItem)

	router.GET("/api/lists", listHandler.GetLists)
	router.GET("/api/lists/:id", listHandler.GetList)
	router.POST("/api/lists", listHandler.CreateList)
	router.DELETE("/api/lists/:id", listHandler.DeleteList)
	router.PUT("/api/lists/:id", listHandler.UpdateListTitle)

	// for testing
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "API is running 🚀",
		})
	})

	return router

}

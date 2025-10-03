package main

import (
	"log"

	"github.com/jennaborowy/fullstack-Go-Docker/config"
	"github.com/jennaborowy/fullstack-Go-Docker/database"
	"github.com/jennaborowy/fullstack-Go-Docker/routes"
)

func main() {
	// load config from .env file
	cfg := config.Load()
	log.Println("Connecting to DB with URL:", cfg.DatabaseURL)

	// Connect to database
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	log.Println("Database connected successfully")

	// create a new gin engine
	r := routes.SetupRoutes(db)

	//define routes
	// r.GET("/", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"message": "Hello World!",
	// 	})
	// })

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}

}

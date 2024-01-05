package main

import (
	"log"
	"os"

	"github.com/PandaX185/tatsumaki-chat/pkg/controllers"
	"github.com/PandaX185/tatsumaki-chat/pkg/database"
	"github.com/PandaX185/tatsumaki-chat/pkg/repository"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file %s", err.Error())
	}

	// Initialize database
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Error initializing database %s", err.Error())
	}
	defer db.Close()

	// Create repository
	r := repository.NewRepository(db)

	// Create controllers
	userController := controllers.NewUserController(r)

	// Register routes
	router := gin.Default()
	router.POST("/register", userController.CreateUser)
	router.GET("/users/:username", userController.GetUser)

	// Run server
	router.Run(os.Getenv("API_PORT"))
}

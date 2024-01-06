package main

import (
	"log"
	"os"

	"github.com/PandaX185/tatsumaki-chat/pkg/controllers"
	"github.com/PandaX185/tatsumaki-chat/pkg/database"
	"github.com/PandaX185/tatsumaki-chat/pkg/kafka"
	"github.com/PandaX185/tatsumaki-chat/pkg/middlewares"
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

	// Initialize kafka
	kConn, err := kafka.InitKafka()
	if err != nil {
		log.Fatalf("Error initializing kafka %s", err.Error())
	}
	defer kConn.Close()

	// Create repository
	r := repository.NewRepository(db)

	router := gin.Default()

	// Create controllers
	userController := controllers.NewUserController(r)
	chatController := controllers.NewChatController(r, kConn)

	// Register middlewares
	router.Use(middlewares.Auth())

	// Register routes
	userController.SetupController(router)
	chatController.SetupController(router)

	// Run server
	router.Run(os.Getenv("API_PORT"))
}

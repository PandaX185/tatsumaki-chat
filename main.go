package main

import (
	"log"
	"net/http"
	"os"

	"github.com/PandaX185/tatsumaki-chat/config"
	"github.com/PandaX185/tatsumaki-chat/domain/chats"
	"github.com/PandaX185/tatsumaki-chat/domain/users"
	"github.com/PandaX185/tatsumaki-chat/migrations"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalln(err)
	}
}

func main() {
	logger, err := config.InitLogger()
	if err != nil {
		log.Fatalln("error initializing logger ", err)
	}

	db, err := config.InitDb()
	if err != nil {
		logger.Fatalln("error initializing database: ", err)
	}
	logger.Infoln("Database is runnning...")

	if err = migrations.Down(db); err != nil {
		logger.Errorln("error rolling back old database migrations: ", err)
	} else {
		logger.Infoln("Migrations rolled back successfully")
	}

	if err = migrations.Up(db); err != nil {
		logger.Errorln("error applying database migrations: ", err)
	} else {
		logger.Infoln("Migrations applied successfully")
	}

	mux := http.NewServeMux()
	path := os.Getenv("PORT")

	// User routes
	userHandler := users.NewHandler(users.NewService(users.NewRepository()))
	mux.HandleFunc("POST /api/users", userHandler.RegisterUser)
	mux.HandleFunc("GET /api/users/{username}", userHandler.GetUserByUsername)

	// Chat routes
	chatHandler := chats.NewHandler(chats.NewService(chats.NewRepository()))
	mux.HandleFunc("POST /api/chats", chatHandler.CreateChat)

	logger.Infof("Starting server on port %v...\n", path)
	if err := http.ListenAndServe(path, mux); err != nil {
		logger.Fatalln("Error starting the server: ", err)
	}
}

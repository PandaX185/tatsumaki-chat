package main

import (
	"log"
	"net/http"
	"os"

	"github.com/PandaX185/tatsumaki-chat/config"
	"github.com/PandaX185/tatsumaki-chat/domain/user"
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

	srv := http.Server{
		Addr: os.Getenv("PORT"),
	}

	userHandler := user.NewHandler(user.NewService(user.NewRepository()), logger)

	http.HandleFunc("POST /api/users", userHandler.RegisterUser)

	logger.Infoln("Starting server on port 8080...")
	if err := srv.ListenAndServe(); err != nil {
		logger.Error("Error starting the server: ", err)
	}
}

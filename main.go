package main

import (
	"log"
	"net/http"

	"github.com/PandaX185/tatsumaki-chat/config"
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

	_, err = config.InitDb()
	if err != nil {
		logger.Fatalln("error initializing database: ", err.Error())
	}

	srv := http.Server{
		Addr: ":8080",
	}

	if err := srv.ListenAndServe(); err != nil {
		logger.Error("Error starting the server: ", err)
	}
}

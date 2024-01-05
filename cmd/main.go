package main

import (
	"github/PandaX185/tatsumaki-chat/pkg/db"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file %s", err.Error())
	}
	router := gin.Default()
	db.InitDB()

}

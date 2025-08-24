package main

import (
	"log"
	"net/http"
	"os"

	"github.com/PandaX185/tatsumaki-chat/config"
	"github.com/PandaX185/tatsumaki-chat/domain/chats"
	"github.com/PandaX185/tatsumaki-chat/domain/messages"
	"github.com/PandaX185/tatsumaki-chat/domain/shared"
	"github.com/PandaX185/tatsumaki-chat/domain/users"
	"github.com/PandaX185/tatsumaki-chat/middlewares"
	"github.com/PandaX185/tatsumaki-chat/migrations"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
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
	mux.HandleFunc("GET /api/users/current", userHandler.GetCurrentUserData)
	mux.HandleFunc("GET /api/users/{username}", userHandler.SearchByUserName)
	mux.HandleFunc("POST /api/users/login", userHandler.Login)

	// Chat routes
	chatHandler := chats.NewHandler(chats.NewService(chats.NewRepository()))
	mux.HandleFunc("POST /api/chats", chatHandler.CreateChat)
	mux.HandleFunc("GET /api/chats", chatHandler.GetAllChats)
	mux.HandleFunc("DELETE /api/chats/{chat_id}", chatHandler.Delete)
	mux.HandleFunc("PATCH /api/chats/{chat_id}", chatHandler.Edit)
	mux.HandleFunc("GET /api/chats/{chat_id}/members", chatHandler.GetChatMembers)

	// Message routes
	messageHandler := messages.NewHandler(messages.NewService(messages.NewRepository(), shared.NewRepository()))
	mux.HandleFunc("POST /api/messages", messageHandler.SendMessage)
	mux.HandleFunc("GET /api/messages/{chat_id}", messageHandler.GetAllMessages)
	mux.Handle("GET /api/realtime/messages", middlewares.VerifyJwtFromQuery(http.HandlerFunc(messageHandler.GetMessagesRealtime)))
	mux.Handle("GET /api/realtime/messages/unread", middlewares.VerifyJwtFromQuery(http.HandlerFunc(messageHandler.GetUnreadMessagesCount)))
	mux.HandleFunc("POST /api/messages/{chat_id}/read", messageHandler.MarkAsRead)

	handler := cors.AllowAll().Handler(mux)
	logger.Infof("Starting server on port %v...\n", path)
	if err := http.ListenAndServe(path, middlewares.VerifyJwt(handler)); err != nil {
		logger.Fatalln("Error starting the server: ", err)
	}
}

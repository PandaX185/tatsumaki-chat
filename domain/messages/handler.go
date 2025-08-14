package messages

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/PandaX185/tatsumaki-chat/config"
	"github.com/PandaX185/tatsumaki-chat/domain/errors"
	"github.com/PandaX185/tatsumaki-chat/domain/errors/codes"
	"github.com/gorilla/websocket"
)

type MessageHandler struct {
	service *MessageService
}

var upgrader websocket.Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewHandler(s *MessageService) *MessageHandler {
	return &MessageHandler{
		service: s,
	}
}

func (h *MessageHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	var body Message
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		jsonErr := errors.JsonError{
			Code:    codes.BAD_REQUEST,
			Message: "Error parsing request body",
		}
		jsonErr.ReturnError(w)
		return
	}

	userId := r.Context().Value("userId")
	userIdInt, _ := strconv.ParseInt(userId.(string), 10, 32)
	body.SenderId = int(userIdInt)

	res, err := h.service.Send(body)
	if err != nil {
		jsonErr := errors.JsonError{
			Code:    codes.INTERNAL,
			Message: "Error sending the message: " + err.Error(),
		}
		jsonErr.ReturnError(w)
		return
	}

	rds := config.GetRedis()
	messageJson, _ := json.Marshal(res)
	fmt.Printf("Publishing message to channel %d: %s\n", body.ChatId, string(messageJson))
	rds.Publish(r.Context(), strconv.Itoa(body.ChatId), string(messageJson))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(codes.CREATED)
	json.NewEncoder(w).Encode(res)
}

func (h *MessageHandler) GetAllMessages(w http.ResponseWriter, r *http.Request) {
	chat_id, err := strconv.ParseInt(r.PathValue("chat_id"), 10, 64)
	if err != nil {
		jsonErr := errors.JsonError{
			Code:    codes.BAD_REQUEST,
			Message: "Provide a correct chat id",
		}
		jsonErr.ReturnError(w)
		return
	}
	userId := r.Context().Value("userId")
	userIdInt, _ := strconv.ParseInt(userId.(string), 10, 32)

	res, err := h.service.GetAll(int(chat_id), int(userIdInt))
	if err != nil {
		jsonErr := errors.JsonError{
			Code:    codes.INTERNAL,
			Message: "Error getting all messages: " + err.Error(),
		}
		jsonErr.ReturnError(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(codes.OK)
	json.NewEncoder(w).Encode(res)
}

func (h *MessageHandler) GetMessagesRealtime(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Cache-Control")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	chatId := r.PathValue("chat_id")
	rds := config.GetRedis()
	pubsub := rds.Subscribe(r.Context(), chatId)
	defer pubsub.Close()

	ch := pubsub.Channel()
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "event: connected\ndata: Connected to chat %s\n\n", chatId)
	flusher.Flush()

	fmt.Printf("Streaming messages for chat %s...\n", chatId)

	notify := r.Context().Done()
	for {
		select {
		case <-notify:
			fmt.Printf("Client disconnected from chat %s\n", chatId)
			return
		case msg, ok := <-ch:
			if !ok {
				fmt.Printf("PubSub channel closed for chat %s\n", chatId)
				return
			}
			fmt.Printf("Received message for chat %s: %v\n", chatId, msg.Payload)

			fmt.Fprintf(w, "event:msg\ndata: %s\n\n", msg.Payload)
			flusher.Flush()
		}
	}
}

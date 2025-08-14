package chats

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

type ChatHandler struct {
	service *ChatService
}

var upgrader websocket.Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewHandler(s *ChatService) *ChatHandler {
	return &ChatHandler{
		service: s,
	}
}

func (h *ChatHandler) CreateChat(w http.ResponseWriter, r *http.Request) {
	var body Chat
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		jsonErr := errors.JsonError{
			Code:    codes.BAD_REQUEST,
			Message: "Error parsing request body",
		}
		fmt.Println(jsonErr)
		return
	}

	chatOwner := r.Context().Value("userId")
	owner, _ := strconv.ParseInt(chatOwner.(string), 10, 32)
	body.ChatOwner = int(owner)

	res, err := h.service.Create(body)
	if err != nil {
		jsonErr := errors.JsonError{
			Code:    codes.INTERNAL,
			Message: "Error creating the chat: " + err.Error(),
		}
		fmt.Println(jsonErr)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *ChatHandler) GetAllChats(w http.ResponseWriter, r *http.Request) {

	chatOwner := r.Context().Value("userId")
	owner, _ := strconv.ParseInt(chatOwner.(string), 10, 32)

	res, err := h.service.GetAllChats(int(owner))
	if err != nil {
		jsonErr := errors.JsonError{
			Code:    codes.INTERNAL,
			Message: "Error getting chats: " + err.Error(),
		}
		fmt.Println(jsonErr)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *ChatHandler) GetChatsRealtime(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	rds := config.GetRedis()
	pubsub := rds.Subscribe(r.Context(), fmt.Sprintf("chats:%s", r.Context().Value("userId")))
	defer pubsub.Close()

	rc := http.NewResponseController(w)
	ch := pubsub.Channel()

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	flusher.Flush()
	notify := r.Context().Done()
	for {
		select {
		case <-notify:
			fmt.Printf("Client disconnected\n")
			return
		case msg, ok := <-ch:
			if !ok {
				fmt.Printf("PubSub channel closed\n")
				return
			}
			fmt.Fprintf(w, "event:chat\ndata: %s\n\n", msg.Payload)
			rc.Flush()
		}
	}
}

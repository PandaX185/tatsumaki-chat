package messages

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/PandaX185/tatsumaki-chat/domain/errors"
	"github.com/PandaX185/tatsumaki-chat/domain/errors/codes"
	"github.com/gorilla/websocket"
	"github.com/lib/pq"
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
	lis := pq.NewListener("messages_listener", 5*time.Second, 5*time.Second, nil)
	if err := lis.Listen("on_message_sent"); err != nil {
		jsonErr := errors.JsonError{
			Code:    codes.INTERNAL,
			Message: "Error listening to messages channel: " + err.Error(),
		}
		jsonErr.ReturnError(w)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		jsonErr := errors.JsonError{
			Code:    codes.INTERNAL,
			Message: "Error upgrading to websocket connection: " + err.Error(),
		}
		jsonErr.ReturnError(w)
		return
	}

	for msg := range lis.Notify {
		fmt.Printf("msg.Extra: %v\n", msg.Extra)
		conn.WriteMessage(1, []byte(msg.Extra))
	}
}

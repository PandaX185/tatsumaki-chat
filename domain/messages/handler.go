package messages

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/PandaX185/tatsumaki-chat/config"
	"github.com/PandaX185/tatsumaki-chat/domain/errors"
	"github.com/PandaX185/tatsumaki-chat/domain/errors/codes"
	"github.com/PandaX185/tatsumaki-chat/domain/shared"
)

type MessageHandler struct {
	service *MessageService
}

func NewHandler(s *MessageService) *MessageHandler {
	return &MessageHandler{
		service: s,
	}
}

func (h *MessageHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	var body shared.Message
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		jsonErr := errors.JsonError{
			Code:    codes.BAD_REQUEST,
			Message: "Error parsing request body",
		}
		fmt.Printf("body: %v\n", body)
		fmt.Printf("err: %v\n", err)
		jsonErr.ReturnError(w)
		return
	}

	userId, _ := strconv.Atoi(r.Context().Value("userId").(string))
	body.SenderId = userId

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
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	userId, _ := strconv.Atoi(r.Context().Value("userId").(string))

	rds := config.GetRedis()
	pubsub := rds.Subscribe(context.Background(), fmt.Sprintf("messages:%d", userId))
	defer pubsub.Close()

	ch := pubsub.Channel()
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "event: connected\ndata: Connected to user %v\n\n", userId)
	flusher.Flush()

	fmt.Printf("User %v connected\n", userId)
	notify := r.Context().Done()
	for {
		select {
		case <-notify:
			fmt.Printf("User %v disconnected\n", userId)
			return
		case msg, ok := <-ch:
			if !ok {
				fmt.Printf("PubSub channel closed for user %v\n", userId)
				return
			}
			fmt.Printf("msg: %v\n", msg.Payload)
			fmt.Fprintf(w, "event: msg\ndata: %v\n\n", msg.Payload)
			flusher.Flush()
		}
	}
}

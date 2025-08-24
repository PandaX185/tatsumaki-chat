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
		fmt.Printf("err: %v\n", err)
		return
	}

	userId, _ := strconv.Atoi(r.Context().Value("userId").(string))
	body.SenderId = userId

	res, err := h.service.Send(body)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(codes.CREATED)
	json.NewEncoder(w).Encode(res)
}

func (h *MessageHandler) GetAllMessages(w http.ResponseWriter, r *http.Request) {
	chatId, _ := strconv.Atoi(r.PathValue("chat_id"))
	userId, _ := strconv.Atoi(r.Context().Value("userId").(string))

	res, err := h.service.GetAll(chatId, userId)
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

	rds := config.GetPubSubClient()
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

	notify := r.Context().Done()
	for {
		select {
		case <-notify:
			return
		case msg, ok := <-ch:
			if !ok {
				return
			}
			fmt.Fprintf(w, "event: msg\ndata: %v\n\n", msg.Payload)
			flusher.Flush()
		}
	}
}

func (h *MessageHandler) GetUnreadMessagesCount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	userId, _ := strconv.Atoi(r.Context().Value("userId").(string))

	notify := r.Context().Done()
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	count, err := h.service.GetUnreadMessagesCount(userId)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}

	jsonCount, _ := json.Marshal(count)
	fmt.Fprintf(w, "event: unread\ndata: %s\n\n", jsonCount)
	flusher.Flush()

	unreadPubSub := config.GetPubSubClient().Subscribe(context.Background(), fmt.Sprintf("unread:%d", userId))
	defer unreadPubSub.Close()

	readPubSub := config.GetPubSubClient().Subscribe(context.Background(), fmt.Sprintf("read:%d", userId))
	defer readPubSub.Close()

	unreadCh := unreadPubSub.Channel()
	readCh := readPubSub.Channel()

	var messages []shared.UnreadMessagesCount
	for {
		select {
		case <-notify:
			return
		case msg := <-unreadCh:
			fmt.Fprintf(w, "event: unread\ndata: %v\n\n", msg.Payload)
			var unread []shared.UnreadMessagesCount
			if err := json.Unmarshal([]byte(msg.Payload), &unread); err == nil {
				messages = unread
			}
			flusher.Flush()
		case msg := <-readCh:
			var readChatId int
			tmpMessages := []shared.UnreadMessagesCount{}
			if err := json.Unmarshal([]byte(msg.Payload), &readChatId); err == nil {
				for i := range messages {
					if messages[i].ChatId != readChatId {
						tmpMessages = append(tmpMessages, messages[i])
					}
				}
				messages = tmpMessages
			}

			fmt.Fprintf(w, "event: unread\ndata: %v\n\n", messages)
			flusher.Flush()
		}
	}
}

func (h *MessageHandler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	chatId, _ := strconv.Atoi(r.PathValue("chat_id"))
	userId, _ := strconv.Atoi(r.Context().Value("userId").(string))
	if err := h.service.MarkAsRead(chatId, userId); err != nil {
		jsonErr := errors.JsonError{
			Code:    codes.INTERNAL,
			Message: "Error marking message as read: " + err.Error(),
		}
		jsonErr.ReturnError(w)
		return
	}

	w.WriteHeader(codes.NO_CONTENT)
}

package chats

import (
	"encoding/json"
	"net/http"

	"github.com/PandaX185/tatsumaki-chat/domain/errors"
)

type ChatHandler struct {
	service *ChatService
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
			Code:    errors.BAD_REQUEST,
			Message: "Error parsing request body",
		}
		jsonErr.ReturnError(w)
		return
	}

	res, err := h.service.Create(body)
	if err != nil {
		jsonErr := errors.JsonError{
			Code:    errors.INTERNAL,
			Message: "Error creating the chat: " + err.Error(),
		}
		jsonErr.ReturnError(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

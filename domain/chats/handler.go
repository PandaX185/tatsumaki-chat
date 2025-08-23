package chats

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/PandaX185/tatsumaki-chat/domain/errors"
	"github.com/PandaX185/tatsumaki-chat/domain/errors/codes"
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
	var body ChatRequest
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

	userId := r.Context().Value("userId")
	user, _ := strconv.ParseInt(userId.(string), 10, 32)

	res, err := h.service.GetAllChats(int(user))
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

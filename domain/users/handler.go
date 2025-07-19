package users

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/PandaX185/tatsumaki-chat/domain/errors"
)

type UserHandler struct {
	service *UserService
}

func NewHandler(s *UserService) *UserHandler {
	return &UserHandler{
		service: s,
	}
}

func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var body User
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		jsonErr := errors.JsonError{
			Code:    errors.BAD_REQUEST,
			Message: "Error parsing request body",
		}
		jsonErr.ReturnError(w)
		return
	}

	res, err := h.service.Save(body)
	if err != nil {
		jsonErr := errors.JsonError{
			Code:    errors.INTERNAL,
			Message: "Error saving the user: " + err.Error(),
		}
		if strings.Contains(strings.ToLower(err.Error()), "duplicate") {
			jsonErr.Code = errors.CONFLICT
			jsonErr.Message = "User already exists"
		}
		jsonErr.ReturnError(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *UserHandler) GetUserByUsername(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")

	res, err := h.service.GetByUserName(username)
	if err != nil {
		jsonErr := errors.JsonError{
			Code:    errors.BAD_REQUEST,
			Message: "Error getting the user: " + err.Error(),
		}
		jsonErr.ReturnError(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

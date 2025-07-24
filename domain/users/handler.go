package users

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/PandaX185/tatsumaki-chat/domain/errors"
	"github.com/PandaX185/tatsumaki-chat/domain/errors/codes"
	"github.com/PandaX185/tatsumaki-chat/utils"
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
			Code:    codes.BAD_REQUEST,
			Message: "Error parsing request body",
		}
		jsonErr.ReturnError(w)
		return
	}

	res, err := h.service.Save(body)
	if err != nil {
		jsonErr := errors.JsonError{
			Code:    codes.INTERNAL,
			Message: "Error saving the user: " + err.Error(),
		}
		if strings.Contains(strings.ToLower(err.Error()), "duplicate") {
			jsonErr.Code = codes.CONFLICT
			jsonErr.Message = "User already exists"
		}
		jsonErr.ReturnError(w)
		return
	}

	accessToken, err := utils.GenerateJwt(map[string]string{
		"userId":   strconv.FormatInt(int64(res.Id), 10),
		"username": res.UserName,
		"fullname": res.FullName,
	})
	if err != nil {
		jsonErr := errors.JsonError{
			Code:    codes.INTERNAL,
			Message: "Error generating JWT token: " + err.Error(),
		}
		jsonErr.ReturnError(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"access_token": accessToken,
	})
}

func (h *UserHandler) GetUserByUsername(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")

	res, err := h.service.GetByUserName(username)
	if err != nil {
		jsonErr := errors.JsonError{
			Code:    codes.BAD_REQUEST,
			Message: "Error getting the user: " + err.Error(),
		}
		jsonErr.ReturnError(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var body User
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		jsonErr := errors.JsonError{
			Code:    codes.BAD_REQUEST,
			Message: "Error parsing request body",
		}
		jsonErr.ReturnError(w)
		return
	}

	res, err := h.service.Login(body.UserName, body.Password)
	if err != nil {
		jsonErr := errors.JsonError{
			Code:    codes.UNAUTHORIZED,
			Message: "Invalid credentials: " + err.Error(),
		}
		jsonErr.ReturnError(w)
		return
	}

	accessToken, err := utils.GenerateJwt(map[string]string{
		"userId":   strconv.FormatInt(int64(res.Id), 10),
		"username": res.UserName,
		"fullname": res.FullName,
	})
	if err != nil {
		jsonErr := errors.JsonError{
			Code:    codes.INTERNAL,
			Message: "Error generating JWT token: " + err.Error(),
		}
		jsonErr.ReturnError(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"access_token": accessToken,
	})
}

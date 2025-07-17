package user

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type UserHandler struct {
	service *UserService
	logger  *zap.SugaredLogger
}

func NewHandler(s *UserService, logger *zap.SugaredLogger) *UserHandler {
	return &UserHandler{
		service: s,
		logger:  logger,
	}
}

func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var body User
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		h.logger.Errorln(err)
		http.Error(w, "Error parsing request body", 400)
		return
	}

	res, err := h.service.Save(body)
	if err != nil {
		h.logger.Errorln(err)
		http.Error(w, "Error saving the user: "+err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		h.logger.Errorln(err)
		http.Error(w, "Error return json response: "+err.Error(), 500)
		return
	}
}

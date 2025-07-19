package errors

import (
	"encoding/json"
	"net/http"
)

type JsonError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

const (
	BAD_REQUEST = 400
	NOT_FOUND   = 404
	CONFLICT    = 409
	INTERNAL    = 500
)

func (err JsonError) ReturnError(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.Code)
	json.NewEncoder(w).Encode(err)
}

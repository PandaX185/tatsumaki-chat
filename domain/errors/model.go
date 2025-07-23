package errors

import (
	"encoding/json"
	"net/http"
)

type JsonError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (err JsonError) ReturnError(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.Code)
	json.NewEncoder(w).Encode(err)
}

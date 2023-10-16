package server

import (
	"encoding/json"
	"net/http"

	"github.com/masonictemple4/masonictemple4.app/internal/dtos"
)

func ErrorHandler(w http.ResponseWriter, r *http.Request, status int, err string) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dtos.ErrorReturn{
		Code:    status,
		Message: err,
	})
}

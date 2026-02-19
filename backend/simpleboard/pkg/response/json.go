package response

import (
	"encoding/json"
	"net/http"
)

const (
	CodeInvalidInput = "INVALID_INPUT"
	CodeNotFound     = "NOT_FOUND"
	CodeUnauthorized = "UNAUTHORIZED"
	CodeForbidden    = "FORBIDDEN"
	CodeConflict     = "CONFLICT"
	CodeGameFull     = "GAME_FULL"
	CodeGameOver     = "GAME_OVER"
	CodeInternal     = "INTERNAL_ERROR"
)

type errorBody struct {
	Error string `json:"error"`
	Code  string `json:"code"`
}

func JSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload != nil {
		json.NewEncoder(w).Encode(payload)
	}
}

func Error(w http.ResponseWriter, status int, code string, message string) {
	JSON(w, status, errorBody{
		Error: message,
		Code:  code,
	})
}

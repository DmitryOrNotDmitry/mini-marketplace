package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Формирует и отправляет ответ с ошибкой в формате JSON
func MakeErrorResponse(w http.ResponseWriter, err error, statusCode int) {
	type ErrorMessage struct {
		Message string
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errResponse := &ErrorMessage{Message: err.Error()}
	if errE := json.NewEncoder(w).Encode(errResponse); errE != nil {
		fmt.Println(errE)
		return
	}
}

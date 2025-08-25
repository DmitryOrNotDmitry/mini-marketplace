package handler

import (
	"errors"
	"net/http"
	"strconv"
)

func (s *Server) ClearCartHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(r.PathValue("user_id"), 10, 64)
	if err != nil {
		MakeErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	if userID < 1 {
		MakeErrorResponse(w, errors.New("user_id must be positive"), http.StatusBadRequest)
		return
	}

	_, err = s.cartService.ClearCart(r.Context(), userID)
	if err != nil {
		MakeErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

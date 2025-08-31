package handler

import (
	"net/http"
)

// ClearCartHandler обрабатывает HTTP-запрос на очистку корзины пользователя.
func (s *Server) ClearCartHandler(w http.ResponseWriter, r *http.Request) {
	var userID int64
	errs := NewRequestValidator(r).
		ParseUserID(&userID).
		Errors()
	if errs != nil {
		MakeErrorResponseByErrs(w, errs)
		return
	}

	err := s.cartService.ClearCart(r.Context(), userID)
	if err != nil {
		MakeErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

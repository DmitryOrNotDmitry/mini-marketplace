package handler

import (
	"net/http"
	"route256/cart/internal/domain"
	"route256/cart/internal/handler/validate"
	"strconv"
)

// Обрабатывает HTTP-запрос на очистку корзины пользователя
func (s *Server) ClearCartHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(r.PathValue("user_id"), 10, 64)
	if err != nil || validate.UserID(userID) != nil {
		MakeErrorResponse(w, domain.ErrUserIDNotValid, http.StatusBadRequest)
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

package handler

import (
	"net/http"
)

// DeleteCartItemHandler обрабатывает HTTP-запрос на удаление товара из корзины пользователя.
func (s *Server) DeleteCartItemHandler(w http.ResponseWriter, r *http.Request) {
	var userID int64
	var skuID int64
	errs := NewRequestValidator(r).
		ParseUserID(&userID).
		ParseSkuID(&skuID).
		Errors()
	if errs != nil {
		MakeErrorResponseByErrs(w, errs)
		return
	}

	err := s.cartService.DeleteCartItem(r.Context(), userID, skuID)
	if err != nil {
		MakeErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

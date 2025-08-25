package handler

import (
	"net/http"
	"route256/cart/internal/domain"
	"strconv"
)

func (s *Server) DeleteCartItemHandler(w http.ResponseWriter, r *http.Request) {
	skuId, err := strconv.ParseInt(r.PathValue("sku_id"), 10, 64)
	if err != nil || skuId < 1 {
		MakeErrorResponse(w, domain.ErrSKUNotValid, http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(r.PathValue("user_id"), 10, 64)
	if err != nil || userID < 1 {
		MakeErrorResponse(w, domain.ErrUserIdNotValid, http.StatusBadRequest)
		return
	}

	_, err = s.cartService.DeleteCartItem(r.Context(), userID, skuId)
	if err != nil {
		MakeErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

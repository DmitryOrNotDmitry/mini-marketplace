package handler

import (
	"net/http"
	"route256/cart/internal/domain"
	"route256/cart/internal/handler/validate"
	"strconv"
)

func (s *Server) DeleteCartItemHandler(w http.ResponseWriter, r *http.Request) {
	skuID, err := strconv.ParseInt(r.PathValue("sku_id"), 10, 64)
	if err != nil || validate.SkuID(skuID) != nil {
		MakeErrorResponse(w, domain.ErrSKUNotValid, http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(r.PathValue("user_id"), 10, 64)
	if err != nil || validate.UserID(userID) != nil {
		MakeErrorResponse(w, domain.ErrUserIDNotValid, http.StatusBadRequest)
		return
	}

	_, err = s.cartService.DeleteCartItem(r.Context(), userID, skuID)
	if err != nil {
		MakeErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

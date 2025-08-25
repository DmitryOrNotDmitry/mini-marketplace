package handler

import (
	"errors"
	"net/http"
	"strconv"
)

func (s *Server) DeleteCartItemHandler(w http.ResponseWriter, r *http.Request) {
	skuId, err := strconv.ParseInt(r.PathValue("sku_id"), 10, 64)
	if err != nil {
		MakeErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(r.PathValue("user_id"), 10, 64)
	if err != nil {
		MakeErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	if skuId < 1 || userID < 1 {
		MakeErrorResponse(w, errors.New("sku, user_id must be positive"), http.StatusBadRequest)
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

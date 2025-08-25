package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"route256/cart/internal/domain"
	"strconv"
)

type AddCartItemRequest struct {
	UserID int64  `json:"user_id"`
	Sku    int64  `json:"sku_id"`
	Count  uint32 `json:"count"`
}

type AddCartItemResponse struct {
	Sku   int64  `json:"sku"`
	Name  string `json:"name"`
	Price uint32 `json:"price"`
	Count uint32 `json:"count"`
}

func (s *Server) AddCartItemHandler(w http.ResponseWriter, r *http.Request) {
	var request AddCartItemRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		MakeErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	var err error
	request.Sku, err = strconv.ParseInt(r.PathValue("sku_id"), 10, 64)
	if err != nil {
		MakeErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	request.UserID, err = strconv.ParseInt(r.PathValue("user_id"), 10, 64)
	if err != nil {
		MakeErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	if request.Sku < 1 || request.UserID < 1 || request.Count < 1 {
		MakeErrorResponse(w, errors.New("sku, user_id and count must be positive"), http.StatusBadRequest)
		return
	}

	cartItem := &domain.CartItem{
		Sku:   request.Sku,
		Count: request.Count,
	}

	addedCartItem, err := s.cartService.AddCartItem(r.Context(), request.UserID, cartItem)
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			MakeErrorResponse(w, err, http.StatusBadRequest)
			return
		}

		MakeErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	response := &AddCartItemResponse{
		Sku:   addedCartItem.Sku,
		Name:  addedCartItem.Name,
		Price: addedCartItem.Price,
		Count: addedCartItem.Count,
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		MakeErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
}

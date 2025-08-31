package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"route256/cart/internal/domain"
)

type AddCartItemRequest struct {
	Count uint32 `json:"count" validate:"required,gt=0"`
}

type AddCartItemResponse struct {
	Sku   int64  `json:"sku"`
	Name  string `json:"name"`
	Price uint32 `json:"price"`
	Count uint32 `json:"count"`
}

// AddCartItemHandler обрабатывает HTTP-запрос на добавление товара в корзину пользователя.
func (s *Server) AddCartItemHandler(w http.ResponseWriter, r *http.Request) {
	fieldErrors := map[string]error{
		"Count": domain.ErrCountNotValid,
	}

	var userID int64
	var skuID int64
	var request AddCartItemRequest
	errs := NewRequestValidator(r).
		ParseUserID(&userID).
		ParseSkuID(&skuID).
		ParseStruct(&request, fieldErrors).
		Errors()
	if errs != nil {
		MakeErrorResponseByErrs(w, errs)
		return
	}

	cartItem := &domain.CartItem{
		Sku:   skuID,
		Count: request.Count,
	}

	addedCartItem, err := s.cartService.AddCartItem(r.Context(), userID, cartItem)
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			MakeErrorResponse(w, domain.ErrProductNotFound, http.StatusPreconditionFailed)
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
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		MakeErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
}

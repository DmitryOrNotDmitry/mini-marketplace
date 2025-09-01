package handler

import (
	"encoding/json"
	"net/http"
	"route256/cart/internal/domain"
)

type CartResponse struct {
	Items      []CartItemResponse `json:"items"`
	TotalPrice uint32             `json:"total_price"`
}

type CartItemResponse struct {
	Sku   int64  `json:"sku"`
	Name  string `json:"name"`
	Price uint32 `json:"price"`
	Count uint32 `json:"count"`
}

// GetCartHandler обрабатывает HTTP-запрос на получение содержимого корзины пользователя.
func (s *Server) GetCartHandler(w http.ResponseWriter, r *http.Request) {
	var userID int64
	errs := NewRequestValidator(r).
		ParseUserID(&userID).
		Errors()
	if errs != nil {
		MakeErrorResponseByErrs(w, errs)
		return
	}

	cart, err := s.cartService.GetCart(r.Context(), userID)
	if err != nil {
		MakeErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	if len(cart.Items) == 0 {
		MakeErrorResponse(w, domain.ErrCartNotFound, http.StatusNotFound)
		return
	}

	response := &CartResponse{
		Items:      make([]CartItemResponse, 0, len(cart.Items)),
		TotalPrice: cart.TotalPrice,
	}

	for _, item := range cart.Items {
		response.Items = append(response.Items, CartItemResponse{
			Sku:   item.Sku,
			Name:  item.Name,
			Count: item.Count,
			Price: item.Price,
		})
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		MakeErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
}

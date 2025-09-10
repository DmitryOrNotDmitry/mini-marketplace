package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"route256/cart/internal/domain"
	"route256/cart/pkg/logger"
)

type CheckoutCartesponse struct {
	OrderID int64 `json:"order_id"`
}

type OrderCheckouter interface {
	OrderCreate(ctx context.Context, userID int64, cart *domain.Cart) (int64, error)
}

// CheckoutCartHandler оформляет заказ по товарам из корзины.
func (s *Server) CheckoutCartHandler(w http.ResponseWriter, r *http.Request) {
	var userID int64
	errs := NewRequestValidator(r).
		ParseUserID(&userID).
		Errors()
	if errs != nil {
		MakeErrorResponseByErrs(w, errs)
		return
	}

	ctx := r.Context()
	cart, err := s.cartService.GetCart(ctx, userID)
	if err != nil {
		MakeErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	if len(cart.Items) == 0 {
		MakeErrorResponse(w, domain.ErrCartNotFound, http.StatusNotFound)
		return
	}

	orderID, err := s.orderCheckouter.OrderCreate(ctx, userID, cart)
	if err != nil {
		MakeErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	err = s.cartService.ClearCart(ctx, userID)
	if err != nil {
		logger.Error(fmt.Sprintf("cartService.ClearCart: %s", err))
	}

	response := &CheckoutCartesponse{
		OrderID: orderID,
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		MakeErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
}

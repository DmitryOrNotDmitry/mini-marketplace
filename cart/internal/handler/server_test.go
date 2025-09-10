package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"route256/cart/internal/domain"
	mock "route256/cart/mocks"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

type testComponentS struct {
	cartServMock       *mock.CartServiceMock
	orderCheckServMock *mock.OrderCheckouterMock
	server             *Server
}

func newTestComponentS(t *testing.T) testComponentS {
	mc := minimock.NewController(t)
	cartServMock := mock.NewCartServiceMock(mc)
	orderCheckServMock := mock.NewOrderCheckouterMock(mc)
	server := NewServer(cartServMock, orderCheckServMock)

	return testComponentS{
		cartServMock:       cartServMock,
		orderCheckServMock: orderCheckServMock,
		server:             server,
	}
}

func TestService(t *testing.T) {
	t.Parallel()

	t.Run("add cart item success", func(t *testing.T) {
		t.Parallel()

		tc := newTestComponentS(t)
		itemReq := AddCartItemRequest{
			Count: 1,
		}
		userID := int64(1)
		skuID := int64(1)

		itemOut := &domain.CartItem{
			Sku:   skuID,
			Name:  "Name 1",
			Price: 100,
			Count: itemReq.Count,
		}

		tc.cartServMock.AddCartItemMock.Return(itemOut, nil)

		res := tc.addCartItem(t, userID, skuID, itemReq)
		require.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run("add cart item with unexisiting product sku", func(t *testing.T) {
		t.Parallel()

		tc := newTestComponentS(t)
		itemReq := AddCartItemRequest{
			Count: 1,
		}
		userID := int64(1)
		skuID := int64(1)

		tc.cartServMock.AddCartItemMock.Return(nil, domain.ErrProductNotFound)

		res := tc.addCartItem(t, userID, skuID, itemReq)
		require.Equal(t, http.StatusPreconditionFailed, res.StatusCode)
	})

	t.Run("add cart item with non-positive userID", func(t *testing.T) {
		t.Parallel()

		tc := newTestComponentS(t)
		itemReq := AddCartItemRequest{
			Count: 1,
		}
		userID := int64(-1)
		skuID := int64(1)

		res := tc.addCartItem(t, userID, skuID, itemReq)
		require.Equal(t, http.StatusBadRequest, res.StatusCode)
	})

	t.Run("add cart item with non-positive skuID", func(t *testing.T) {
		t.Parallel()

		tc := newTestComponentS(t)
		itemReq := AddCartItemRequest{
			Count: 1,
		}
		userID := int64(1)
		skuID := int64(-1)

		res := tc.addCartItem(t, userID, skuID, itemReq)
		require.Equal(t, http.StatusBadRequest, res.StatusCode)
	})

	t.Run("add cart item with zero count", func(t *testing.T) {
		t.Parallel()

		tc := newTestComponentS(t)
		itemReq := AddCartItemRequest{
			Count: 0,
		}
		userID := int64(1)
		skuID := int64(1)

		res := tc.addCartItem(t, userID, skuID, itemReq)
		require.Equal(t, http.StatusBadRequest, res.StatusCode)
	})

	t.Run("get cart success", func(t *testing.T) {
		t.Parallel()

		tc := newTestComponentS(t)
		itemReq := AddCartItemRequest{
			Count: 1,
		}
		userID := int64(1)
		skuID := int64(1)

		itemOut := &domain.CartItem{
			Sku:   skuID,
			Name:  "Name 1",
			Price: 100,
			Count: itemReq.Count,
		}
		cartOut := &domain.Cart{
			Items:      []*domain.CartItem{itemOut},
			TotalPrice: itemOut.Count * itemOut.Price,
		}

		tc.cartServMock.GetCartMock.Return(cartOut, nil)

		cart := tc.getCart(t, userID)
		require.Len(t, cart.Items, 1)

		assert.Equal(t, *cartOut.Items[0], *cart.Items[0])
		assert.Equal(t, cartOut.TotalPrice, cart.TotalPrice)
	})

	t.Run("delete cart item success", func(t *testing.T) {
		t.Parallel()

		tc := newTestComponentS(t)
		userID := int64(1)
		skuID := int64(1)

		tc.cartServMock.DeleteCartItemMock.Return(nil)

		res := tc.deleteCartItem(t, userID, skuID)
		require.Equal(t, http.StatusNoContent, res.StatusCode)
	})

	t.Run("delete cart success", func(t *testing.T) {
		t.Parallel()

		tc := newTestComponentS(t)
		userID := int64(1)

		tc.cartServMock.ClearCartMock.Return(nil)

		res := tc.deleteCart(t, userID)
		require.Equal(t, http.StatusNoContent, res.StatusCode)
	})

	t.Run("checkout cart success", func(t *testing.T) {
		t.Parallel()

		tc := newTestComponentS(t)

		userID := int64(1)
		expectedOrderID := int64(1)
		cart := &domain.Cart{Items: []*domain.CartItem{
			&domain.CartItem{Sku: 1, Count: 10},
		}}

		tc.cartServMock.GetCartMock.When(minimock.AnyContext, userID).Then(cart, nil)
		tc.orderCheckServMock.OrderCreateMock.Return(expectedOrderID, nil)
		tc.cartServMock.ClearCartMock.When(minimock.AnyContext, userID).Then(nil)

		orderID, res := tc.checkoutOrder(t, userID)
		require.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, expectedOrderID, orderID)
	})

	t.Run("checkout cart failed: empty cart", func(t *testing.T) {
		t.Parallel()

		tc := newTestComponentS(t)

		userID := int64(1)
		cart := &domain.Cart{Items: []*domain.CartItem{}}

		tc.cartServMock.GetCartMock.When(minimock.AnyContext, userID).Then(cart, nil)

		_, res := tc.checkoutOrder(t, userID)
		require.Equal(t, http.StatusNotFound, res.StatusCode)
	})

	t.Run("checkout cart failed: order checkouter error", func(t *testing.T) {
		t.Parallel()

		tc := newTestComponentS(t)

		userID := int64(1)
		cart := &domain.Cart{Items: []*domain.CartItem{
			&domain.CartItem{Sku: 1, Count: 10},
		}}

		tc.cartServMock.GetCartMock.When(minimock.AnyContext, userID).Then(cart, nil)
		tc.orderCheckServMock.OrderCreateMock.Return(0, errors.New("error"))

		_, res := tc.checkoutOrder(t, userID)
		require.Equal(t, http.StatusInternalServerError, res.StatusCode)
	})
}

func (tc testComponentS) checkoutOrder(t *testing.T, userID int64) (int64, *http.Response) {
	t.Helper()

	reader := bytes.NewReader([]byte{})
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/checkout/%d", userID), reader)
	req.SetPathValue("user_id", fmt.Sprint(userID))
	w := httptest.NewRecorder()

	tc.server.CheckoutCartHandler(w, req)

	res := w.Result()
	defer res.Body.Close()

	orderIDRes := &CheckoutCartesponse{}
	json.NewDecoder(res.Body).Decode(orderIDRes)

	return orderIDRes.OrderID, res
}

func (tc testComponentS) addCartItem(t *testing.T, userID, skuID int64, item AddCartItemRequest) *http.Response {
	t.Helper()

	body, err := json.Marshal(item)
	require.NoError(t, err)
	reader := bytes.NewReader(body)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/user/%d/cart/%d", userID, skuID), reader)
	req.SetPathValue("user_id", fmt.Sprint(userID))
	req.SetPathValue("sku_id", fmt.Sprint(skuID))
	w := httptest.NewRecorder()

	tc.server.AddCartItemHandler(w, req)

	res := w.Result()
	defer res.Body.Close()

	return res
}

func (tc testComponentS) deleteCartItem(t *testing.T, userID, skuID int64) *http.Response {
	t.Helper()

	reader := bytes.NewReader([]byte{})
	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/user/%d/cart/%d", userID, skuID), reader)
	req.SetPathValue("user_id", fmt.Sprint(userID))
	req.SetPathValue("sku_id", fmt.Sprint(skuID))
	w := httptest.NewRecorder()

	tc.server.DeleteCartItemHandler(w, req)

	res := w.Result()
	defer res.Body.Close()

	return res
}

func (tc testComponentS) deleteCart(t *testing.T, userID int64) *http.Response {
	t.Helper()

	reader := bytes.NewReader([]byte{})
	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/user/%d/cart", userID), reader)
	req.SetPathValue("user_id", fmt.Sprint(userID))
	w := httptest.NewRecorder()

	tc.server.ClearCartHandler(w, req)

	res := w.Result()
	defer res.Body.Close()

	return res
}

func (tc testComponentS) getCart(t *testing.T, userID int64) *domain.Cart {
	t.Helper()

	reader := bytes.NewReader(nil)
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/user/%d/cart", userID), reader)
	req.SetPathValue("user_id", fmt.Sprint(userID))
	w := httptest.NewRecorder()

	tc.server.GetCartHandler(w, req)

	res := w.Result()
	require.Equal(t, http.StatusOK, res.StatusCode)
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	cartResponce := &CartResponse{}
	err := decoder.Decode(cartResponce)
	require.NoError(t, err)

	cart := &domain.Cart{
		Items:      make([]*domain.CartItem, 0, len(cartResponce.Items)),
		TotalPrice: cartResponce.TotalPrice,
	}
	for _, itemResp := range cartResponce.Items {
		cart.Items = append(cart.Items, &domain.CartItem{
			Sku:   itemResp.Sku,
			Name:  itemResp.Name,
			Count: itemResp.Count,
			Price: itemResp.Price,
		})
	}

	return cart
}

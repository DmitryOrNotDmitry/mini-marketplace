//go:build api

package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"route256/cart/internal/handler"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

func (s *Server) addCartItem(t provider.StepCtx, userID, skuID int64, itemRequest *handler.AddCartItemRequest) *http.Response {
	body, err := json.Marshal(itemRequest)
	t.Require().NoError(err, "marshal request")
	reader := bytes.NewReader(body)

	resp, err := http.Post(fmt.Sprintf("%v/user/%d/cart/%d", s.Host, userID, skuID), "application/json", reader)
	t.Require().NoError(err, "http post")
	defer resp.Body.Close()

	t.Require().Equal(http.StatusOK, resp.StatusCode, "must be StatusOK")

	return resp
}

func (s *Server) deleteCartItem(t provider.StepCtx, userID, skuID int64) *http.Response {
	reader := bytes.NewReader([]byte{})
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%v/user/%d/cart/%d", s.Host, userID, skuID), reader)
	t.Require().NoError(err, "create DELETE request")

	resp, err := http.DefaultClient.Do(req)
	t.Require().NoError(err, "send DELETE request")
	defer resp.Body.Close()

	t.Require().Equal(http.StatusNoContent, resp.StatusCode, "must be StatusNoContent")

	return resp
}

func (s *Server) deleteCart(t provider.StepCtx, userID int64) *http.Response {
	reader := bytes.NewReader([]byte{})
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%v/user/%d/cart", s.Host, userID), reader)
	t.Require().NoError(err, "create DELETE request")

	resp, err := http.DefaultClient.Do(req)
	t.Require().NoError(err, "send DELETE request")
	defer resp.Body.Close()

	t.Require().Equal(http.StatusNoContent, resp.StatusCode, "must be StatusNoContent")

	return resp
}

func (s *Server) getCart(t provider.StepCtx, userID int64, expectedStatus int) *http.Response {
	resp, err := http.Get(fmt.Sprintf("%v/user/%d/cart", s.Host, userID))
	t.Require().NoError(err, "http get")

	t.Require().Equal(expectedStatus, resp.StatusCode, fmt.Sprintf("must be %d", expectedStatus))

	return resp
}

func (s *Server) decodeCartResponse(t provider.StepCtx, resp *http.Response) *handler.CartResponse {
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	cartResponse := &handler.CartResponse{}
	err := decoder.Decode(cartResponse)
	t.Require().NoError(err, "decode get responce")

	return cartResponse
}

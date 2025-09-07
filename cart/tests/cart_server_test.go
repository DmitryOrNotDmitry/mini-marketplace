//go:build api

package tests

import (
	"net/http"
	"testing"

	"route256/cart/internal/handler"

	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type Server struct {
	suite.Suite

	Host string
}

func TestServer(t *testing.T) {
	t.Parallel()

	suite.RunSuite(t, new(Server))
}

func (s *Server) BeforeAll(t provider.T) {
	s.Host = "http://localhost:8080"
	t.Logf("host is %v", s.Host)
}

func (s *Server) BeforeEach(t provider.T) {
	t.Feature("Cart Service")
	t.Tags("cart", "backend", "go")
	t.Owner("Dima Cuznetsov")
	t.Labels(
		&allure.Label{Name: "platform", Value: "backed"},
	)
}

func (s *Server) TestDeleteCartItem(t provider.T) {
	t.Title("Добавляем товар и удаляем его из корзины")

	itemReq := &handler.AddCartItemRequest{
		Count: 1,
	}
	userID := int64(1)
	skuID := int64(1076963)

	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Добавляем товар", func(t provider.StepCtx) {
			s.addCartItem(t, userID, skuID, itemReq)
		})
	})

	t.WithNewStep("Удаляем товар", func(t provider.StepCtx) {
		s.deleteCartItem(t, userID, skuID)
	})

	t.WithNewStep("Получаем пустую корзину", func(t provider.StepCtx) {
		s.getCart(t, userID, http.StatusNotFound)
	})
}

func (s *Server) TestGetCart(t provider.T) {
	t.Title("Получаем корзину с двумя товарами")

	itemReq1 := &handler.AddCartItemRequest{
		Count: 1,
	}
	itemReq2 := &handler.AddCartItemRequest{
		Count: 1,
	}
	userID := int64(1)
	skuID1 := int64(1076963)
	skuID2 := int64(1625903)

	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Добавляем 3 товара", func(t provider.StepCtx) {
			s.addCartItem(t, userID, skuID1, itemReq1)
			s.addCartItem(t, userID, skuID1, itemReq1)
			s.addCartItem(t, userID, skuID2, itemReq2)
		})
	})

	t.WithNewStep("Получаем корзину", func(t provider.StepCtx) {
		resp := s.getCart(t, userID, http.StatusOK)
		cartResponse := s.decodeCartResponse(t, resp)

		t.Require().Len(cartResponse.Items, 2)
		t.Logf("cartResponse: %v", cartResponse)

		t.Assert().EqualValues(2, cartResponse.Items[0].Count)
		t.Assert().EqualValues(1, cartResponse.Items[1].Count)
	})

	t.WithNewStep("Очищаем корзину", func(t provider.StepCtx) {
		s.deleteCart(t, userID)
	})
}

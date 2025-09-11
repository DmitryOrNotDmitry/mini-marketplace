package tests

import (
	"context"
	"net/http"
	"testing"

	"route256/tests/helpers/clients"

	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type Server struct {
	suite.Suite

	cartClient *clients.CartClient
}

func TestServer(t *testing.T) {
	t.Parallel()

	suite.RunSuite(t, new(Server))
}

func (s *Server) BeforeAll(t provider.T) {
	url := "http://localhost:8080"
	s.cartClient = clients.NewClient(url)
	t.Logf("url is %v", url)
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

	count := uint32(1)
	userID := int64(1)
	skuID := int64(30816475)
	ctx := context.Background()

	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Добавляем товар", func(t provider.StepCtx) {
			resStatus := s.cartClient.AddCartItem(ctx, t, userID, skuID, count)
			t.Require().Equal(http.StatusOK, resStatus, "не совпадает статус код")
		})
	})

	t.WithNewStep("Удаляем товар", func(t provider.StepCtx) {
		resStatus := s.cartClient.DeleteCartItem(ctx, t, userID, skuID)
		t.Require().Equal(http.StatusNoContent, resStatus, "не совпадает статус код")
	})

	t.WithNewStep("Получаем пустую корзину", func(t provider.StepCtx) {
		_, resStatus := s.cartClient.GetCart(ctx, t, userID)
		t.Require().Equal(http.StatusNotFound, resStatus, "не совпадает статус код")
	})
}

func (s *Server) TestGetCart(t provider.T) {
	t.Title("Получаем корзину с двумя товарами")

	count1 := uint32(1)
	count2 := uint32(1)
	userID := int64(2)
	skuID1 := int64(4465995)
	skuID2 := int64(30816475)
	ctx := context.Background()

	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Добавляем 3 товара", func(t provider.StepCtx) {
			resStatus := s.cartClient.AddCartItem(ctx, t, userID, skuID1, count1)
			t.Require().Equal(http.StatusOK, resStatus, "не совпадает статус код")

			resStatus = s.cartClient.AddCartItem(ctx, t, userID, skuID1, count1)
			t.Require().Equal(http.StatusOK, resStatus, "не совпадает статус код")

			resStatus = s.cartClient.AddCartItem(ctx, t, userID, skuID2, count2)
			t.Require().Equal(http.StatusOK, resStatus, "не совпадает статус код")
		})
	})

	t.WithNewStep("Получаем корзину", func(t provider.StepCtx) {
		cart, resStatus := s.cartClient.GetCart(ctx, t, userID)
		t.Require().Equal(http.StatusOK, resStatus, "не совпадает статус код")

		t.Require().Len(cart.Items, 2)
		t.Logf("cartResponse: %v", cart)

		t.Assert().EqualValues(2, cart.Items[0].Count)
		t.Assert().EqualValues(1, cart.Items[1].Count)
	})

	t.WithNewStep("Очищаем корзину", func(t provider.StepCtx) {
		resStatus := s.cartClient.ClearCart(ctx, t, userID)
		t.Require().Equal(http.StatusNoContent, resStatus, "не совпадает статус код")
	})
}

//go:build api

package tests

import (
	"context"
	"net/http"
	"sync"
	"testing"

	"route256/tests/helpers/clients"

	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type LomsSuite struct {
	suite.Suite

	cartClient *clients.CartClient
	lomsClient *clients.LomsClient
}

func TestLomsSuite(t *testing.T) {
	t.Parallel()

	suite.RunSuite(t, new(LomsSuite))
}

func (cs *LomsSuite) BeforeAll(t provider.T) {
	cartUrl := "http://localhost:8080"
	cs.cartClient = clients.NewCartClient(cartUrl)
	t.Logf("cart url is %v", cartUrl)

	lomsUrl := "http://localhost:8084"
	cs.lomsClient = clients.NewLomsClient(lomsUrl)
	t.Logf("loms url is %v", lomsUrl)
}

func (cs *LomsSuite) BeforeEach(t provider.T) {
	t.Feature("Loms+Cart Services")
	t.Tags("cart", "backend", "go")
	t.Owner("Dima Cuznetsov")
	t.Labels(
		&allure.Label{Name: "platform", Value: "backed"},
	)
}

func (cs *LomsSuite) TestOrderInfoAndPay(t provider.T) {
	t.Title("Оформляем заказ успешно и оплачиваем его")

	count := uint32(1)
	userID := int64(3)
	skuID := int64(30816475)
	ctx := context.Background()

	var orderID int64

	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Добавляем товар", func(t provider.StepCtx) {
			resStatus := cs.cartClient.AddCartItem(ctx, t, userID, skuID, count)
			t.Require().Equal(http.StatusOK, resStatus, "не совпадает статус код")
		})
	})

	t.WithNewStep("Оформляем заказ", func(t provider.StepCtx) {
		var resStatus int
		orderID, resStatus = cs.cartClient.Checkout(ctx, t, userID)
		t.Require().Equal(http.StatusOK, resStatus, "не совпадает статус код")
	})

	t.WithNewStep("Проверяем, что заказ создан", func(t provider.StepCtx) {
		order, resStatus := cs.lomsClient.OrderInfo(ctx, t, orderID)
		t.Require().Equal(http.StatusOK, resStatus, "не совпадает статус код")
		t.Require().Len(order.Items, 1, "в заказе должен быть 1 item")
		t.Assert().Equal(skuID, order.Items[0].SkuID)
	})

	t.WithNewStep("Получаем пустую корзину", func(t provider.StepCtx) {
		_, resStatus := cs.cartClient.GetCart(ctx, t, userID)
		t.Require().Equal(http.StatusNotFound, resStatus, "не совпадает статус код")
	})

	t.WithNewStep("Оплачиваем заказ", func(t provider.StepCtx) {
		resStatus := cs.lomsClient.PayOrder(ctx, t, orderID)
		t.Require().Equal(http.StatusOK, resStatus, "не совпадает статус код")
	})

	t.WithNewStep("Проверяем статус оплаченного заказа", func(t provider.StepCtx) {
		order, resStatus := cs.lomsClient.OrderInfo(ctx, t, orderID)
		t.Require().Equal(http.StatusOK, resStatus, "не совпадает статус код")
		t.Assert().EqualValues("paid", order.Status)
	})
}

func (cs *LomsSuite) TestOrderCancel(t provider.T) {
	t.Title("Оформляем заказ успешно и отменяем его")

	count := uint32(1)
	userID := int64(4)
	skuID := int64(30816475)
	ctx := context.Background()

	var orderID int64

	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Добавляем товар", func(t provider.StepCtx) {
			resStatus := cs.cartClient.AddCartItem(ctx, t, userID, skuID, count)
			t.Require().Equal(http.StatusOK, resStatus, "не совпадает статус код")
		})
	})

	t.WithNewStep("Оформляем заказ", func(t provider.StepCtx) {
		var resStatus int
		orderID, resStatus = cs.cartClient.Checkout(ctx, t, userID)
		t.Require().Equal(http.StatusOK, resStatus, "не совпадает статус код")
	})

	t.WithNewStep("Отменяем заказ", func(t provider.StepCtx) {
		resStatus := cs.lomsClient.CancelOrder(ctx, t, orderID)
		t.Require().Equal(http.StatusOK, resStatus, "не совпадает статус код")
	})

	t.WithNewStep("Проверяем статус отмененного заказа", func(t provider.StepCtx) {
		order, resStatus := cs.lomsClient.OrderInfo(ctx, t, orderID)
		t.Require().Equal(http.StatusOK, resStatus, "не совпадает статус код")
		t.Assert().EqualValues("cancelled", order.Status)
	})

	t.WithNewStep("Пытаемся оплатить отмененный заказ", func(t provider.StepCtx) {
		resStatus := cs.lomsClient.PayOrder(ctx, t, orderID)
		t.Require().Equal(http.StatusBadRequest, resStatus, "не совпадает статус код")
	})

	t.WithNewStep("Отменяем уже отмененный заказ (успешно)", func(t provider.StepCtx) {
		resStatus := cs.lomsClient.CancelOrder(ctx, t, orderID)
		t.Require().Equal(http.StatusOK, resStatus, "не совпадает статус код")
	})
}

func (cs *LomsSuite) TestCreateManyOrders(t provider.T) {
	t.Title("Оформляем заказ успешно 100 раз параллельно и проверяем стоки")

	userID := int64(5)
	skuID1 := int64(3618852)
	skuID2 := int64(4288068)
	unexistedSkuID := int64(4288069)
	ctx := context.Background()

	order := &clients.Order{
		UserID: userID,
		Items: []*clients.OrderItem{
			&clients.OrderItem{SkuID: skuID1, Count: 1},
			&clients.OrderItem{SkuID: skuID2, Count: 1},
		},
	}
	failOrder := &clients.Order{
		UserID: userID,
		Items: []*clients.OrderItem{
			&clients.OrderItem{SkuID: skuID1, Count: 1},
			&clients.OrderItem{SkuID: unexistedSkuID, Count: 1},
		},
	}

	var skuID1stock uint32
	var skuID2stock uint32

	t.WithNewStep("Проверяем стоки", func(t provider.StepCtx) {
		var resStatus int
		skuID1stock, resStatus = cs.lomsClient.StockInfo(ctx, t, skuID1)
		t.Require().Equal(http.StatusOK, resStatus, "не совпадает статус код")

		skuID2stock, resStatus = cs.lomsClient.StockInfo(ctx, t, skuID2)
		t.Require().Equal(http.StatusOK, resStatus, "не совпадает статус код")
	})

	t.WithNewStep("Добавляем 100 заказов", func(t provider.StepCtx) {
		var finish sync.WaitGroup
		finish.Add(100)

		for i := 0; i < 50; i++ {
			go func() {
				defer finish.Done()

				var resStatus int
				_, resStatus = cs.lomsClient.CreateOrder(ctx, t, order)
				t.Require().Equal(http.StatusOK, resStatus, "не совпадает статус код")
			}()

			go func() {
				defer finish.Done()

				var resStatus int
				_, resStatus = cs.lomsClient.CreateOrder(ctx, t, failOrder)
				t.Require().Equal(http.StatusBadRequest, resStatus, "не совпадает статус код")
			}()
		}

		finish.Wait()
	})

	t.WithNewStep("Проверяем стоки после создания заказов", func(t provider.StepCtx) {
		var resStatus int
		curSkuID1stock, resStatus := cs.lomsClient.StockInfo(ctx, t, skuID1)
		t.Require().Equal(http.StatusOK, resStatus, "не совпадает статус код")
		t.Require().EqualValues(skuID1stock-50, curSkuID1stock)

		curSkuID2stock, resStatus := cs.lomsClient.StockInfo(ctx, t, skuID2)
		t.Require().Equal(http.StatusOK, resStatus, "не совпадает статус код")
		t.Require().EqualValues(skuID2stock-50, curSkuID2stock)
	})

}

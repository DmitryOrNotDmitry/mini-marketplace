package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

type CartClient struct {
	baseUrl string
	cl      *http.Client
}

func NewClient(baseUrl string) *CartClient {
	return &CartClient{
		baseUrl: baseUrl,
		cl:      http.DefaultClient,
	}
}

func (c *CartClient) AddCartItem(ctx context.Context, t provider.StepCtx, userID, skuID int64, count uint32) int {
	body, err := json.Marshal(AddCartItemRequest{
		Count: count,
	})
	t.Require().NoError(err, "сериализация тела запроса")
	t.WithNewAttachment("request payload", allure.JSON, body)

	url := fmt.Sprintf("%v/user/%d/cart/%d", c.baseUrl, userID, skuID)
	r, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	t.Require().NoError(err, "создание запроса")

	r.Header.Add("Content-Type", "application/json")

	res, err := c.cl.Do(r)
	t.Require().NoError(err, "выполнение запроса")
	defer res.Body.Close()

	return res.StatusCode
}

func (c *CartClient) DeleteCartItem(ctx context.Context, t provider.StepCtx, userID, skuID int64) int {
	url := fmt.Sprintf("%v/user/%d/cart/%d", c.baseUrl, userID, skuID)
	r, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, bytes.NewBuffer([]byte{}))
	t.Require().NoError(err, "создание запроса")

	res, err := c.cl.Do(r)
	t.Require().NoError(err, "выполнение запроса")
	defer res.Body.Close()

	return res.StatusCode
}

func (c *CartClient) ClearCart(ctx context.Context, t provider.StepCtx, userID int64) int {
	url := fmt.Sprintf("%v/user/%d/cart", c.baseUrl, userID)
	r, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, bytes.NewBuffer([]byte{}))
	t.Require().NoError(err, "создание запроса")

	res, err := c.cl.Do(r)
	t.Require().NoError(err, "выполнение запроса")
	defer res.Body.Close()

	return res.StatusCode
}

func (c *CartClient) GetCart(ctx context.Context, t provider.StepCtx, userID int64) (*Cart, int) {
	url := fmt.Sprintf("%s/user/%d/cart", c.baseUrl, userID)
	r, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	t.Require().NoError(err, "создание запроса")

	res, err := c.cl.Do(r)
	t.Require().NoError(err, "выполнение запроса")
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	t.Require().NoError(err, "считывание ответа")
	t.WithNewAttachment("response body", allure.JSON, body)

	if res.StatusCode != http.StatusOK {
		return nil, res.StatusCode
	}

	cartResp := &GetCartResponse{}
	err = json.Unmarshal(body, cartResp)
	t.Require().NoError(err, "парсинг ответа")

	cart := &Cart{
		Items:      make([]*CartItem, 0, len(cartResp.Items)),
		TotalPrice: cartResp.TotalPrice,
	}
	for _, item := range cartResp.Items {
		cart.Items = append(cart.Items, &CartItem{
			SkuID: item.SkuID,
			Count: item.Count,
			Name:  item.Name,
			Price: item.Price,
		})
	}

	return cart, res.StatusCode
}

func (c *CartClient) Checkout(ctx context.Context, t provider.StepCtx, userID int64) (orderID int64, statusCode int) {
	data, err := json.Marshal(CheckoutOrderRequest{
		UserID: userID,
	})
	t.Require().NoError(err, "сериализация тела запроса")
	t.WithNewAttachment("request payload", allure.JSON, data)

	url := fmt.Sprintf("%s/checkout/%d", c.baseUrl, userID)
	r, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	t.Require().NoError(err, "создание запроса")

	r.Header.Add("Content-Type", "application/json")

	res, err := c.cl.Do(r)
	t.Require().NoError(err, "выполнение запроса")
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return 0, res.StatusCode
	}

	body, err := io.ReadAll(res.Body)
	t.Require().NoError(err, "не удалось считать ответ")
	t.WithNewAttachment("response body", allure.JSON, body)

	checkoutResp := &CheckoutOrderResponse{}
	err = json.Unmarshal(body, checkoutResp)
	t.Require().NoError(err, "парсинг ответа")

	return checkoutResp.OrderID, res.StatusCode
}

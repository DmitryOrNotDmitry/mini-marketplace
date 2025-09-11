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

type LomsClient struct {
	baseUrl string
	cl      *http.Client
}

func NewLomsClient(baseUrl string) *LomsClient {
	return &LomsClient{
		baseUrl: baseUrl,
		cl:      http.DefaultClient,
	}
}

func (c *LomsClient) CreateOrder(ctx context.Context, t provider.StepCtx, order *Order) (int64, int) {
	req := OrderCreateRequest{
		UserID: order.UserID,
		Items:  make([]*OrderCreateItemRequest, 0, len(order.Items)),
	}
	for _, item := range order.Items {
		req.Items = append(req.Items, &OrderCreateItemRequest{
			SkuID: item.SkuID,
			Count: item.Count,
		})
	}

	body, err := json.Marshal(req)
	t.Require().NoError(err, "сериализация тела запроса")
	t.WithNewAttachment("request payload", allure.JSON, body)

	url := fmt.Sprintf("%v/order/create", c.baseUrl)
	r, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	t.Require().NoError(err, "создание запроса")

	r.Header.Add("Content-Type", "application/json")

	res, err := c.cl.Do(r)
	t.Require().NoError(err, "выполнение запроса")
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return -1, res.StatusCode
	}

	body, err = io.ReadAll(res.Body)
	t.Require().NoError(err, "считывание ответа")
	t.WithNewAttachment("response body", allure.JSON, body)

	respBody := &OrderCreateResponse{}
	err = json.Unmarshal(body, respBody)
	t.Require().NoError(err, "парсинг ответа")

	return respBody.OrderID, res.StatusCode
}

func (c *LomsClient) PayOrder(ctx context.Context, t provider.StepCtx, orderID int64) int {
	req := OrderPayRequest{
		OrderID: orderID,
	}

	body, err := json.Marshal(req)
	t.Require().NoError(err, "сериализация тела запроса")
	t.WithNewAttachment("request payload", allure.JSON, body)

	url := fmt.Sprintf("%v/order/pay", c.baseUrl)
	r, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	t.Require().NoError(err, "создание запроса")

	r.Header.Add("Content-Type", "application/json")

	res, err := c.cl.Do(r)
	t.Require().NoError(err, "выполнение запроса")
	defer res.Body.Close()

	return res.StatusCode
}

func (c *LomsClient) CancelOrder(ctx context.Context, t provider.StepCtx, orderID int64) int {
	req := OrderCancelRequest{
		OrderID: orderID,
	}

	body, err := json.Marshal(req)
	t.Require().NoError(err, "сериализация тела запроса")
	t.WithNewAttachment("request payload", allure.JSON, body)

	url := fmt.Sprintf("%v/order/cancel", c.baseUrl)
	r, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	t.Require().NoError(err, "создание запроса")

	r.Header.Add("Content-Type", "application/json")

	res, err := c.cl.Do(r)
	t.Require().NoError(err, "выполнение запроса")
	defer res.Body.Close()

	return res.StatusCode
}

func (c *LomsClient) OrderInfo(ctx context.Context, t provider.StepCtx, orderID int64) (*Order, int) {
	url := fmt.Sprintf("%v/order/info?orderId=%d", c.baseUrl, orderID)
	r, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	t.Require().NoError(err, "создание запроса")

	r.Header.Add("Content-Type", "application/json")

	res, err := c.cl.Do(r)
	t.Require().NoError(err, "выполнение запроса")
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, res.StatusCode
	}

	body, err := io.ReadAll(res.Body)
	t.Require().NoError(err, "считывание ответа")
	t.WithNewAttachment("response body", allure.JSON, body)

	respBody := &OrderInfoResponse{}
	t.Log(string(body))
	err = json.Unmarshal(body, respBody)
	t.Require().NoError(err, "парсинг ответа")

	order := &Order{
		UserID: respBody.UserID,
		Status: string(respBody.Status),
		Items:  make([]*OrderItem, 0, len(respBody.Items)),
	}
	for _, item := range respBody.Items {
		order.Items = append(order.Items, &OrderItem{
			SkuID: item.SkuID,
			Count: item.Count,
		})
	}

	return order, res.StatusCode
}

func (c *LomsClient) StockInfo(ctx context.Context, t provider.StepCtx, skuID int64) (uint32, int) {
	req := StockInfoRequest{
		SkuID: skuID,
	}

	body, err := json.Marshal(req)
	t.Require().NoError(err, "сериализация тела запроса")
	t.WithNewAttachment("request payload", allure.JSON, body)

	url := fmt.Sprintf("%v/stock/info?sku=%d", c.baseUrl, skuID)
	r, err := http.NewRequestWithContext(ctx, http.MethodGet, url, bytes.NewBuffer(body))
	t.Require().NoError(err, "создание запроса")

	r.Header.Add("Content-Type", "application/json")

	res, err := c.cl.Do(r)
	t.Require().NoError(err, "выполнение запроса")
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return 0, res.StatusCode
	}

	body, err = io.ReadAll(res.Body)
	t.Require().NoError(err, "считывание ответа")
	t.WithNewAttachment("response body", allure.JSON, body)

	respBody := &StockInfoResponse{}
	err = json.Unmarshal(body, respBody)
	t.Require().NoError(err, "парсинг ответа")

	return respBody.Count, res.StatusCode
}

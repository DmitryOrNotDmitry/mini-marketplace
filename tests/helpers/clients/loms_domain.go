package clients

type Order struct {
	OrderID int64
	UserID  int64
	Status  string
	Items   []*OrderItem
}

type OrderItem struct {
	SkuID int64
	Count uint32
}

type OrderCreateRequest struct {
	UserID int64                     `json:"userId,string"`
	Items  []*OrderCreateItemRequest `json:"items"`
}

type OrderCreateResponse struct {
	OrderID int64 `json:"orderId,string"`
}

type OrderCreateItemRequest struct {
	SkuID int64  `json:"sku,string"`
	Count uint32 `json:"count"`
}

type OrderPayRequest struct {
	OrderID int64 `json:"orderId,string"`
}

type OrderCancelRequest struct {
	OrderID int64 `json:"orderId,string"`
}

type OrderInfoRequest struct {
	OrderID int64 `json:"orderId,string"`
}

type OrderInfoResponse struct {
	UserID int64                    `json:"userId,string"`
	Status string                   `json:"status"`
	Items  []*OrderInfoItemResponse `json:"items"`
}

type OrderInfoItemResponse struct {
	SkuID int64  `json:"sku,string"`
	Count uint32 `json:"count"`
}

type StockInfoRequest struct {
	SkuID int64 `json:"sku,string"`
}

type StockInfoResponse struct {
	Count uint32 `json:"count"`
}

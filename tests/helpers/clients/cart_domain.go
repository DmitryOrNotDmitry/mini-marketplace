package clients

type AddCartItemRequest struct {
	Count uint32 `json:"count"`
}

type GetCartResponse struct {
	Items      []*GetCartItemResponse `json:"items"`
	TotalPrice uint32                 `json:"total_price"`
}

type GetCartItemResponse struct {
	SkuID int64  `json:"sku_id"`
	Count uint32 `json:"count"`
	Price uint32 `json:"price"`
	Name  string `json:"name"`
}

type CheckoutOrderRequest struct {
	UserID int64 `json:"user_id"`
}

type CheckoutOrderResponse struct {
	OrderID int64 `json:"order_id"`
}

type Cart struct {
	Items      []*CartItem
	TotalPrice uint32
}

type CartItem struct {
	SkuID int64
	Count uint32
	Price uint32
	Name  string
}

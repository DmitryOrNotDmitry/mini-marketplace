package domain

// OrderItem хранит данные о товаре в заказе.
type OrderItem struct {
	SkuID int64
	Count uint32
}

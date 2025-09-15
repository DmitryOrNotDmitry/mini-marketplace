package domain

// Cart хранит данные о корзине.
type Cart struct {
	Items      []*CartItem
	TotalPrice uint32
}

package domain

// CartItem хранит данные о товаре в корзине.
type CartItem struct {
	Sku   int64
	Name  string
	Count uint32
	Price uint32
}

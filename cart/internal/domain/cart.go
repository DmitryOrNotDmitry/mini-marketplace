package domain

type Cart struct {
	Items      []*CartItem
	TotalPrice uint32
}

package domain

type Cart struct {
	Items      []*CartItem
	TotalPrice int64
}

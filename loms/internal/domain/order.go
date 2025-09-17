package domain

const (
	New             Status = "new"
	AwaitingPayment Status = "awaiting payment"
	Failed          Status = "failed"
	Paid            Status = "paid"
	Cancelled       Status = "cancelled"
)

// Order хранит данные о заказе пользователя.
type Order struct {
	OrderID int64
	UserID  int64
	Status  Status
	Items   []*OrderItem
}

// Status тип для статуса заказа.
type Status string

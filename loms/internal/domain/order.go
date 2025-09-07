package domain

type Order struct {
	OrderID int64
	UserID  int64
	Status  Status
	Items   []*OrderItem
}

type Status string

const (
	New             Status = "new"
	AwaitingPayment Status = "awaiting payment"
	Failed          Status = "failed"
	Payed           Status = "payed"
	Cancelled       Status = "cancelled"
)

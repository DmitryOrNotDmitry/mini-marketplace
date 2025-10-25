package domain

import "time"

// Comment отражает сущность комментария на товар от пользователя
type Comment struct {
	ID        int64
	UserID    int64
	Sku       int64
	Content   string
	CreatedAt time.Time
}

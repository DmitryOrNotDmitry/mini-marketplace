package domain

// EventStatus описывает статус события в outbox
type EventStatus string

const (
	// EventNew - событие создано, но не обработано
	EventNew EventStatus = "new"

	// Complete - событие обработано
	Complete EventStatus = "complete"

	// Dead - при обработке события возникла ошибка
	Dead EventStatus = "dead"
)

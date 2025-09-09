package domain

import "errors"

var ErrCanNotReserveItem = errors.New("недостаточно товара для резервирования")
var ErrItemStockNotExist = errors.New("в стоке нет такого товара")

var ErrOrderNotExist = errors.New("заказа с таким ID не существует")
var ErrEmptyOrderItems = errors.New("список товаров не должен быть пустым")
var ErrPayWithInvalidOrderStatus = errors.New("оплата заказа в невалидном статусе невозможна")
var ErrCancelWithInvalidOrderStatus = errors.New("невозможно отменить неудавшийся или оплаченный заказ")

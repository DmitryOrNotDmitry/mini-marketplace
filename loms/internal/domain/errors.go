package domain

import "errors"

var ErrCanNotReserveItem = errors.New("недостаточно товара для резервирования")
var ErrItemStockNotExist = errors.New("в стоке нет такого товара")

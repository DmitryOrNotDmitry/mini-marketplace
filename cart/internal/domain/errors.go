package domain

import "errors"

var ErrCartNotFound = errors.New("у пользователя пустая корзина")
var ErrProductNotFound = errors.New("SKU не существует")

var ErrSKUNotValid = errors.New("SKU должен быть натуральным числом (больше нуля)")
var ErrUserIDNotValid = errors.New("идентификатор пользователя должен быть натуральным числом (больше нуля)")
var ErrCountNotValid = errors.New("количество должно быть натуральным числом (больше нуля)")

var ErrOutOfStock = errors.New("невозможно добавить товара по количеству больше, чем есть в стоках")

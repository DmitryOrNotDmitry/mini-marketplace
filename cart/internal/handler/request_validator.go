package handler

import (
	"encoding/json"
	"net/http"
	"route256/cart/internal/domain"
	"route256/cart/internal/handler/validate"
	"strconv"
)

type RequestValidator struct {
	errs      []error
	validator *validate.ValidatorAdapter
	r         *http.Request
}

// NewRequestValidator конструктор для InputValidator.
func NewRequestValidator(r *http.Request) *RequestValidator {
	return &RequestValidator{
		errs:      make([]error, 0),
		validator: validate.NewValidatorAdapter(),
		r:         r,
	}
}

// ParseUserID парсит из запроса и валидирует userID.
func (iv *RequestValidator) ParseUserID(userID *int64) *RequestValidator {
	var err error
	*userID, err = strconv.ParseInt(iv.r.PathValue("user_id"), 10, 64)
	if err != nil || iv.validator.UserID(*userID) != nil {
		iv.errs = append(iv.errs, domain.ErrUserIDNotValid)
	}
	return iv
}

// ParseSkuID парсит из запроса и валидирует skuID.
func (iv *RequestValidator) ParseSkuID(skuID *int64) *RequestValidator {
	var err error
	*skuID, err = strconv.ParseInt(iv.r.PathValue("sku_id"), 10, 64)
	if err != nil || iv.validator.UserID(*skuID) != nil {
		iv.errs = append(iv.errs, domain.ErrSKUNotValid)
	}
	return iv
}

// ParseStruct парсит из тела запроса (json) в структуру и валидирует её.
// s - указатель на структуру.
func (iv *RequestValidator) ParseStruct(s any, fieldErrors map[string]error) *RequestValidator {
	if err := json.NewDecoder(iv.r.Body).Decode(s); err != nil {
		iv.errs = append(iv.errs, err)
	}

	err := iv.validator.Struct(s, fieldErrors)
	if err != nil {
		iv.errs = append(iv.errs, err)
	}
	return iv
}

// Errors возвращает слайс ошибок, если они были. Иначе nil.
func (iv *RequestValidator) Errors() []error {
	if len(iv.errs) == 0 {
		return nil
	}
	return iv.errs
}

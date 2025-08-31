package validate

import "github.com/go-playground/validator/v10"

type ValidatorAdapter struct {
	validate *validator.Validate
}

// NewValidatorAdapter конструктор для ValidatorAdapter
func NewValidatorAdapter() *ValidatorAdapter {
	return &ValidatorAdapter{
		validate: validator.New(),
	}
}

// Struct выполняет валидацию структуры с использованием ошибок для полей
func (va *ValidatorAdapter) Struct(s any, fieldErrors map[string]error) error {
	if err := va.validate.Struct(s); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, e := range validationErrors {
				if customErr, found := fieldErrors[e.Field()]; found {
					return customErr
				}
			}
		}
		return err
	}
	return nil
}

func validateVar(va *ValidatorAdapter, s any, tag string) error {
	return va.validate.Var(s, tag)
}

// UserID валидирует идентификатор пользователя.
func (va *ValidatorAdapter) UserID(userID int64) error {
	return validateVar(va, userID, "required,gt=0")
}

// SkuID валидирует SKU товара.
func (va *ValidatorAdapter) SkuID(skuID int64) error {
	return validateVar(va, skuID, "required,gt=0")
}

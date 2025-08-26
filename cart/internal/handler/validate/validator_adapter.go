package validate

import "github.com/go-playground/validator/v10"

var validate = validator.New()

// Выполняет валидацию структуры с использованием ошибок для полей
func Struct[T any](s T, fieldErrors map[string]error) error {
	if err := validate.Struct(s); err != nil {
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

func validateVar[T any](s T, tag string) error {
	return validate.Var(s, tag)
}

// Валидирует идентификатор пользователя
func UserID(userID int64) error {
	return validateVar(userID, "required,gt=0")
}

// Валидирует SKU товара
func SkuID(skuID int64) error {
	return validateVar(skuID, "required,gt=0")
}

package validation

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	validator2 "github.com/go-playground/validator/v10"

	"awesomeProject5/subscriberOrder/internal/models"
)

var validator *validator2.Validate

func Validate(data []byte) (*models.Order, error) {

	orderReq := &models.Order{}

	validator = validator2.New()
	if err := validator.RegisterValidation("currency", validateCurrency); err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(data, &orderReq); err != nil {
		return nil, fmt.Errorf("unmarshalling error: %v", err)
	}

	if err := validator.Struct(orderReq); err != nil {
		if _, ok := err.(*validator2.InvalidValidationError); ok {
			return nil, fmt.Errorf("validation error: %v", err)

		}

		for _, err := range err.(validator2.ValidationErrors) {
			return nil, fmt.Errorf("validation error in field '%s': %s", err.StructField(), err.Tag())
		}
	} else {
		return orderReq, nil
	}

	return nil, nil

}

func validateCurrency(fl validator2.FieldLevel) bool {
	currency := fl.Field().String()
	return strings.EqualFold(currency, "usd") || strings.EqualFold(currency, "rub")
}

package api

import (
	"github.com/fsobh/simplebank/util"
	"github.com/go-playground/validator/v10"
)

// this custom validator MUST get registered with Gin to be used (check server.go file)
var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {

	//FieldLevel.Field() <-- get value of Field (its a reflection value, so we added the .Interface() to it)
	//.(string) is to covert the value to string

	//will return the currency as a string and an ok boolean value if successful
	if currency, ok := fieldLevel.Field().Interface().(string); ok {

		return util.IsSupportedCurrency(currency) //check if currency is supported
	}

	return false
}

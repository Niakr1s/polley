package helpers

import "github.com/go-playground/validator/v10"

// Validate is a global validate variable.
var Validate *validator.Validate

func init() {
	Validate = validator.New()
}

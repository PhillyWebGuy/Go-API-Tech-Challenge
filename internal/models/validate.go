package models

import "github.com/go-playground/validator/v10"

// Validate is the instance of the validator used to validate struct fields.
var Validate *validator.Validate

// init initializes the Validate variable with a new validator instance.
// This function is automatically called when the package is imported.
func init() {
	Validate = validator.New()
}

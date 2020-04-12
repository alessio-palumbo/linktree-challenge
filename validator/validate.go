package validator

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	validator "gopkg.in/go-playground/validator.v9"
)

const (
	validationRequiredField   = "is required"
	validationInvalidField    = "is invalid"
	validationRequiredWithout = "is required in absence of"
	validationMaxLength       = "is longer than"
	validationMaxSize         = "is greater than"

	lkDateFormat = "Jan 02 2006"
)

// CustomValidator is a custom payload validator
type CustomValidator struct {
	validator *validator.Validate
}

// New returns a new instance of CustomValidator
func New() *CustomValidator {
	return &CustomValidator{validator: validator.New()}
}

// Validate applies the validation rules specified in the payload `validate tag` and returns an error
func (cv *CustomValidator) Validate(i interface{}) error {
	cv.registerCustomValidations()

	err := cv.validator.Struct(i)
	if err == nil {
		return nil
	}

	trans := make([]string, len(err.(validator.ValidationErrors)))
	for i, vErr := range err.(validator.ValidationErrors) {
		trans[i] = formatTranslation(vErr)
	}

	return fmt.Errorf("validation errors: %s", strings.Join(trans, ", "))
}

func (cv *CustomValidator) registerCustomValidations() {
	cv.validator.RegisterValidation("lkDate", validateLkDate)
}

func formatTranslation(vErr validator.FieldError) string {
	var field = vErr.StructField()
	var tag = vErr.Tag()

	switch tag {
	case "required":
		return translate(field, validationRequiredField)
	case "required_without":
		return translate(field, validationRequiredWithout, vErr.Param())
	case "max":
		if vErr.Kind() == reflect.String {
			return translate(field, validationMaxLength, vErr.Param(), "characters")
		}
		return translate(field, validationMaxSize, vErr.Param())
	}

	return translate(field, validationInvalidField)
}

func translate(field string, validationType string, params ...string) string {

	trans := fmt.Sprintf("%s %s", field, validationType)
	for _, v := range params {
		trans += " " + v
	}

	return trans
}

// TODO this could be a more general date format
func validateLkDate(fl validator.FieldLevel) bool {
	_, err := time.Parse(lkDateFormat, fl.Field().String())
	return err == nil
}

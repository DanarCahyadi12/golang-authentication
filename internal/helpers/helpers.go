package helpers

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

func GetFirstValidationErrorsAndConvert(validationErrors error) string {
	var field string
	var tag string
	var param string
	if errs := validationErrors.(validator.ValidationErrors); errs != nil {
		for _, err := range errs {
			field = err.Field()
			tag = err.Tag()
			param = err.Param()
			break
		}
	}

	message := fmt.Sprintf("%s must be %s %s", field, tag, param)
	return strings.TrimSpace(message)
}

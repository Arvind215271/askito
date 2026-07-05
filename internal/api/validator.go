package api

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func Validate(i interface{}) error {

	err := validate.Struct(i)

	if err == nil {
		return nil
	}

	appErr := Err.Common.Validation()

	for _, fe := range err.(validator.ValidationErrors) {

		field := fe.Field()

		var message string

		switch fe.Tag() {

		case "required":
			message = fmt.Sprintf(
				"%s is required",
				field,
			)

		case "email":
			message = fmt.Sprintf(
				"%s must be a valid email",
				field,
			)

		case "min":
			message = fmt.Sprintf(
				"%s must be at least %s characters",
				field,
				fe.Param(),
			)

		case "max":
			message = fmt.Sprintf(
				"%s must be at most %s characters",
				field,
				fe.Param(),
			)

		case "len":
			message = fmt.Sprintf(
				"%s must be exactly %s characters",
				field,
				fe.Param(),
			)

		case "gte":
			message = fmt.Sprintf(
				"%s must be greater than or equal to %s",
				field,
				fe.Param(),
			)

		case "lte":
			message = fmt.Sprintf(
				"%s must be less than or equal to %s",
				field,
				fe.Param(),
			)

		case "gt":
			message = fmt.Sprintf(
				"%s must be greater than %s",
				field,
				fe.Param(),
			)

		case "lt":
			message = fmt.Sprintf(
				"%s must be less than %s",
				field,
				fe.Param(),
			)

		case "hexcolor":
			message = fmt.Sprintf(
				"%s must be a valid hex color",
				field,
			)

		default:
			message = fmt.Sprintf(
				"%s is invalid",
				field,
			)
		}

		appErr.AddField(field, message)
	}

	return appErr
}

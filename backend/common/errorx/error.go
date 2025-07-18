package errorx

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/status"
)

func GetErrorMsg(err error) string {
	if err == nil {
		return ""
	}
	// rpc error
	if s, ok := status.FromError(err); ok {
		return s.Message()
	}
	// validator error
	if errs, ok := err.(validator.ValidationErrors); ok && len(errs) > 0 {
		firstErr := errs[0]
		return getValidatorErrMsg(firstErr)
	}

	return err.Error()
}

func getValidatorErrMsg(f validator.FieldError) string {
	if f == nil {
		return ""
	}

	switch f.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", f.Field())
	case "email":
		return fmt.Sprintf("%s is invalid", f.Field())
	case "min":
		return fmt.Sprintf("%s should be at least %s", f.Field(), f.Param())
	case "max":
		return fmt.Sprintf("%s should be at most %s", f.Field(), f.Param())
	case "len":
		return fmt.Sprintf("%s length must be %s", f.Field(), f.Param())
	case "eqfield":
		return fmt.Sprintf("%s must be equal to %s", f.Field(), f.Param())
	case "nefield":
		return fmt.Sprintf("%s cannot be equal to %s", f.Field(), f.Param())
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", f.Field(), f.Param())
	case "lt":
		return fmt.Sprintf("%s must be less than %s", f.Field(), f.Param())
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", f.Field(), f.Param())
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", f.Field(), f.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of %s", f.Field(), f.Param())
	case "notoneof":
		return fmt.Sprintf("%s cannot be one of %s", f.Field(), f.Param())
	case "datetime":
		return fmt.Sprintf("%s is an invalid datetime", f.Field())
	case "uuid":
		return fmt.Sprintf("%s is an invalid UUID", f.Field())
	case "url":
		return fmt.Sprintf("%s is an invalid URL", f.Field())
	case "json":
		return fmt.Sprintf("%s is invalid JSON", f.Field())
	case "xml":
		return fmt.Sprintf("%s is invalid XML", f.Field())
	case "base64":
		return fmt.Sprintf("%s is invalid Base64", f.Field())
	case "contains":
		return fmt.Sprintf("%s must contain %s", f.Field(), f.Param())
	case "excludes":
		return fmt.Sprintf("%s cannot contain %s", f.Field(), f.Param())
	case "starts":
		return fmt.Sprintf("%s must start with %s", f.Field(), f.Param())
	case "ends":
		return fmt.Sprintf("%s must end with %s", f.Field(), f.Param())
	case "is":
		return fmt.Sprintf("%s must be %s", f.Field(), f.Param())
	case "ascii":
		return fmt.Sprintf("%s must be ASCII characters", f.Field())
	case "printascii":
		return fmt.Sprintf("%s must be printable ASCII characters", f.Field())
	case "multibyte":
		return fmt.Sprintf("%s must be multibyte characters", f.Field())
	case "printmultibyte":
		return fmt.Sprintf("%s must be printable multibyte characters", f.Field())
	case "lowercase":
		return fmt.Sprintf("%s must be lowercase characters", f.Field())
	case "uppercase":
		return fmt.Sprintf("%s must be uppercase characters", f.Field())
	case "titlecase":
		return fmt.Sprintf("%s must be title-case characters", f.Field())
	case "number":
		return fmt.Sprintf("%s must be a number", f.Field())
	case "alpha":
		return fmt.Sprintf("%s must be alphabetic characters", f.Field())
	case "alphanum":
		return fmt.Sprintf("%s must be alphanumeric characters", f.Field())
	case "hexadecimal":
		return fmt.Sprintf("%s must be a hexadecimal number", f.Field())
	case "hexcolor":
		return fmt.Sprintf("%s must be a hexadecimal color", f.Field())
	case "rgb":
		return fmt.Sprintf("%s must be an RGB color", f.Field())
	case "rgba":
		return fmt.Sprintf("%s must be an RGBA color", f.Field())
	case "hsl":
		return fmt.Sprintf("%s must be an HSL color", f.Field())
	case "hsla":
		return fmt.Sprintf("%s must be an HSLA color", f.Field())
	default:
		return f.Error()
	}

}

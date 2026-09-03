package constants

import (
	"fmt"
	"slices"
)

type ApiMethod string

const ApiMethodName = "api_method_name"

const (
	ApiMethodGet    ApiMethod = "GET"
	ApiMethodPost   ApiMethod = "POST"
	ApiMethodPut    ApiMethod = "PUT"
	ApiMethodDelete ApiMethod = "DELETE"
	ApiMethodPatch  ApiMethod = "PATCH"
	ApiMethodHead   ApiMethod = "HEAD"
)

func (m ApiMethod) Validate() error {

	ranges := []ApiMethod{
		ApiMethodGet,
		ApiMethodPost,
		ApiMethodPut,
		ApiMethodDelete,
		ApiMethodPatch,
		ApiMethodHead,
	}
	return validateConstant(m, ranges, ApiMethodName)

	// switch m {
	// case ApiMethodGet,
	// 	ApiMethodPost,
	// 	ApiMethodPut,
	// 	ApiMethodDelete,
	// 	ApiMethodPatch,
	// 	ApiMethodHead:
	// 	return true
	// default:
	// 	return false
	//
	// }
}

func validateConstant[t comparable](target t, ranges []t, constName string) error {
	if slices.Contains(ranges, target) {
			return nil
		}
	return fmt.Errorf("invalid %s: %v", constName, target)
}

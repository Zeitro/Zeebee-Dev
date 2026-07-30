package constants

import "fmt"

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
	for _, t2 := range ranges {
		if target == t2 {
			return nil
		}
	}
	return fmt.Errorf("invalid %s: %v", constName, target)
}

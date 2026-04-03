package constants

type ApiMethod string

const (
	ApiMethodGet    ApiMethod = "GET"
	ApiMethodPost   ApiMethod = "POST"
	ApiMethodPut    ApiMethod = "PUT"
	ApiMethodDelete ApiMethod = "DELETE"
	ApiMethodPatch  ApiMethod = "PATCH"
	ApiMethodHead   ApiMethod = "HEAD"
)

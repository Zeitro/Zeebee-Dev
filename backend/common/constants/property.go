package constants

import "zeebee-backend/common/enums"

// type:Single Family,Townhouse,Multi Family,Condo
type PropertyType string

const (
	PropertyTypeSingleFamily PropertyType = "Single Family"
	PropertyTypeTownhouse    PropertyType = "Townhouse"
	PropertyTypeMultiFamily  PropertyType = "Multi Family"
	PropertyTypeCondo        PropertyType = "Condo"
)

var PropertyTypeName = "property_type"

func init() {
	// register property type
	enums.Register(PropertyTypeName, []enums.EnumOpt{
		{
			Label: string(PropertyTypeSingleFamily),
			Value: string(PropertyTypeSingleFamily),
		},
		{
			Label: string(PropertyTypeTownhouse),
			Value: string(PropertyTypeTownhouse),
		},
		{
			Label: string(PropertyTypeMultiFamily),
			Value: string(PropertyTypeMultiFamily),
		},
		{
			Label: string(PropertyTypeCondo),
			Value: string(PropertyTypeCondo),
		},
	})
}

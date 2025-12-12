package constants

import "zeebee-backend/common/enums"

// type:'Single Family Residence' '2 unit' '3unit' '4 unit' '4+ Multi' 'Townhouse' 'PUD/Condominium' 'Mixed Use' 'Commercial'
type PropertyType string

const (
	PropertyTypeSingleFamilyResidence PropertyType = "Single Family Residence"
	PropertyTypeTwoUnit               PropertyType = "2 unit"
	PropertyTypeThreeUnit             PropertyType = "3 unit"
	PropertyTypeFourUnit              PropertyType = "4 unit"
	PropertyTypeFourPlusMulti         PropertyType = "4+ Multi"
	PropertyTypeTownhouse             PropertyType = "Townhouse"
	PropertyTypePUDCondominium        PropertyType = "PUD/Condominium"
	PropertyTypeMixedUse              PropertyType = "Mixed Use"
	PropertyTypeCommercial            PropertyType = "Commercial"
)

var PropertyTypeName = "property_type"

func init() {
	// register property type
	enums.Register(PropertyTypeName, []enums.EnumOpt{
		{
			Label: string(PropertyTypeSingleFamilyResidence),
			Value: string(PropertyTypeSingleFamilyResidence),
		},
		{
			Label: string(PropertyTypeTwoUnit),
			Value: string(PropertyTypeTwoUnit),
		},
		{
			Label: string(PropertyTypeThreeUnit),
			Value: string(PropertyTypeThreeUnit),
		},
		{
			Label: string(PropertyTypeFourUnit),
			Value: string(PropertyTypeFourUnit),
		},
		{
			Label: string(PropertyTypeFourPlusMulti),
			Value: string(PropertyTypeFourPlusMulti),
		},
		{
			Label: string(PropertyTypeTownhouse),
			Value: string(PropertyTypeTownhouse),
		},
		{
			Label: string(PropertyTypePUDCondominium),
			Value: string(PropertyTypePUDCondominium),
		},
		{
			Label: string(PropertyTypeMixedUse),
			Value: string(PropertyTypeMixedUse),
		},
		{
			Label: string(PropertyTypeCommercial),
			Value: string(PropertyTypeCommercial),
		},
	})
}

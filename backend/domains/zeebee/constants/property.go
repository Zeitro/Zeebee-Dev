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
	PropertyTypeRawLand               PropertyType = "Raw Land"
	PropertyTypeEntitledLand          PropertyType = "Entitled Land"
)

var PropertyTypeName = "property_type"

type BorrowerPropertyStatus string

const (
	BorrowerPropertyStatusActive             BorrowerPropertyStatus = "Active"
	BorrowerPropertyStatusPurchasePending    BorrowerPropertyStatus = "Purchase Pending"
	BorrowerPropertyStatusReFinancePending   BorrowerPropertyStatus = "Re-finance Pending"
	BorrowerPropertyStatusChangeTitlePending BorrowerPropertyStatus = "Change Title Pending"
	BorrowerPropertyStatusSalePending        BorrowerPropertyStatus = "Sale Pending"
	BorrowerPropertyStatusAssignmentPending  BorrowerPropertyStatus = "Assignment Pending"
	BorrowerPropertyStatusSold               BorrowerPropertyStatus = "Sold"
)

var BorrowerPropertyStatusName = "borrower_property_status"

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
		{
			Label: string(PropertyTypeRawLand),
			Value: string(PropertyTypeRawLand),
		},
		{
			Label: string(PropertyTypeEntitledLand),
			Value: string(PropertyTypeEntitledLand),
		},
	})

	// register borrower property status
	enums.Register(BorrowerPropertyStatusName, []enums.EnumOpt{
		{
			Label: string(BorrowerPropertyStatusActive),
			Value: string(BorrowerPropertyStatusActive),
		},
		{
			Label: string(BorrowerPropertyStatusPurchasePending),
			Value: string(BorrowerPropertyStatusPurchasePending),
		},
		{
			Label: string(BorrowerPropertyStatusReFinancePending),
			Value: string(BorrowerPropertyStatusReFinancePending),
		},
		{
			Label: string(BorrowerPropertyStatusChangeTitlePending),
			Value: string(BorrowerPropertyStatusChangeTitlePending),
		},
		{
			Label: string(BorrowerPropertyStatusSalePending),
			Value: string(BorrowerPropertyStatusSalePending),
		},
		{
			Label: string(BorrowerPropertyStatusAssignmentPending),
			Value: string(BorrowerPropertyStatusAssignmentPending),
		},
		{
			Label: string(BorrowerPropertyStatusSold),
			Value: string(BorrowerPropertyStatusSold),
		},
	})
}

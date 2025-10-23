package validate_rules

import (
	"github.com/duke-git/lancet/v2/convertor"
	"github.com/go-playground/validator/v10"

	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/common/constants"
)

var PropertyUnitTag = "invalid_unit_for_property_type"
var PropertyUnitMsgTpl = "{0} property type requires unit value of {1}"

// LoanPreApprovalRequestProperty
func LoanPreApprovalRequestProperty(sl validator.StructLevel) {
	property, ok := sl.Current().Interface().(types.LoanPreApprovalRequestProperty)
	if !ok {
		sl.ReportError(property, "property", "Property", "property_type", "invalid property")
		return
	}
	switch property.PropertyType {
	case string(constants.PropertyTypeSingleFamilyResidence):
		if property.Unit != 1 {
			sl.ReportError(property.Unit, "unit", "Unit", PropertyUnitTag, convertor.ToString(property.Unit))
		}
	case string(constants.PropertyTypeTwoUnit):
		if property.Unit != 2 {
			sl.ReportError(property.Unit, "unit", "Unit", PropertyUnitTag, string(property.Unit))
		}
	case string(constants.PropertyTypeThreeUnit):
		if property.Unit != 3 {
			sl.ReportError(property.Unit, "unit", "Unit", PropertyUnitTag, convertor.ToString(property.Unit))
		}
	case string(constants.PropertyTypeFourUnit):
		if property.Unit != 4 {
			sl.ReportError(property.Unit, "unit", "Unit", PropertyUnitTag, convertor.ToString(property.Unit))
		}
	case string(constants.PropertyTypeFourPlusMulti):
		if property.Unit < 4 {
			sl.ReportError(property.Unit, "unit", "Unit", PropertyUnitTag, convertor.ToString(property.Unit))
		}
	default:
		if property.Unit < 1 {
			sl.ReportError(property.Unit, "unit", "Unit", PropertyUnitTag, " >= 1")
		}
	}
}

// LoanApplyPropertyEntry
func LoanApplyPropertyEntry(sl validator.StructLevel) {
	property, ok := sl.Current().Interface().(types.LoanApplyPropertyEntry)
	if !ok {
		sl.ReportError(property, "property", "Property", "property_type", "invalid property")
		return
	}
	switch property.PropertyType {
	case string(constants.PropertyTypeSingleFamilyResidence):
		if property.Unit != 1 {
			sl.ReportError(property.Unit, "unit", "Unit", PropertyUnitTag, convertor.ToString(property.Unit))
		}
	case string(constants.PropertyTypeTwoUnit):
		if property.Unit != 2 {
			sl.ReportError(property.Unit, "unit", "Unit", PropertyUnitTag, convertor.ToString(property.Unit))
		}
	case string(constants.PropertyTypeThreeUnit):
		if property.Unit != 3 {
			sl.ReportError(property.Unit, "unit", "Unit", PropertyUnitTag, convertor.ToString(property.Unit))
		}
	case string(constants.PropertyTypeFourUnit):
		if property.Unit != 4 {
			sl.ReportError(property.Unit, "unit", "Unit", PropertyUnitTag, convertor.ToString(property.Unit))
		}
	case string(constants.PropertyTypeFourPlusMulti):
		if property.Unit < 4 {
			sl.ReportError(property.Unit, "unit", "Unit", PropertyUnitTag, convertor.ToString(property.Unit))
		}
	default:
		if property.Unit < 1 {
			sl.ReportError(property.Unit, "unit", "Unit", PropertyUnitTag, " >= 1")
		}
	}
}

// LoanPropertyEntry
func LoanPropertyEntry(sl validator.StructLevel) {
	property, ok := sl.Current().Interface().(types.LoanApplyPropertyEntry)
	if !ok {
		sl.ReportError(property, "property", "Property", "property_type", "invalid property")
		return
	}
	switch property.PropertyType {
	case string(constants.PropertyTypeSingleFamilyResidence):
		if property.Unit != 1 {
			sl.ReportError(property.Unit, "unit", "Unit", PropertyUnitTag, convertor.ToString(property.Unit))
		}
	case string(constants.PropertyTypeTwoUnit):
		if property.Unit != 2 {
			sl.ReportError(property.Unit, "unit", "Unit", PropertyUnitTag, convertor.ToString(property.Unit))
		}
	case string(constants.PropertyTypeThreeUnit):
		if property.Unit != 3 {
			sl.ReportError(property.Unit, "unit", "Unit", PropertyUnitTag, convertor.ToString(property.Unit))
		}
	case string(constants.PropertyTypeFourUnit):
		if property.Unit != 4 {
			sl.ReportError(property.Unit, "unit", "Unit", PropertyUnitTag, convertor.ToString(property.Unit))
		}
	case string(constants.PropertyTypeFourPlusMulti):
		if property.Unit < 4 {
			sl.ReportError(property.Unit, "unit", "Unit", PropertyUnitTag, convertor.ToString(property.Unit))
		}
	default:
		if property.Unit < 1 {
			sl.ReportError(property.Unit, "unit", "Unit", PropertyUnitTag, " >= 1")
		}
	}
}

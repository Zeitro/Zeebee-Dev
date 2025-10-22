package validate_rules

import (
	"github.com/go-playground/validator/v10"

	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/common/constants"
)

func LoanPreApprovalRequestProperty(sl validator.StructLevel) {
	property, ok := sl.Current().Interface().(types.LoanPreApprovalRequestProperty)
	if !ok {
		sl.ReportError(property, "property", "Property", "property_type", "xx")
		return
	}
	switch property.PropertyType {
	case string(constants.PropertyTypeSingleFamilyResidence):
		if property.Unit != 1 {
			sl.ReportError(property.Unit, "unit", "Unit", string(constants.PropertyTypeSingleFamilyResidence), "1")
		}
	case string(constants.PropertyTypeTwoUnit):
		if property.Unit != 2 {
			sl.ReportError(property.Unit, "unit", "Unit", string(constants.PropertyTypeTwoUnit), "2")
		}
	case string(constants.PropertyTypeThreeUnit):
		if property.Unit != 3 {
			sl.ReportError(property.Unit, "unit", "Unit", string(constants.PropertyTypeThreeUnit), "3")
		}
	case string(constants.PropertyTypeFourUnit):
		if property.Unit != 4 {
			sl.ReportError(property.Unit, "unit", "Unit", string(constants.PropertyTypeFourUnit), "4")
		}
	case string(constants.PropertyTypeFourPlusMulti):
		if property.Unit < 4 {
			sl.ReportError(property.Unit, "unit", "Unit", string(constants.PropertyTypeFourPlusMulti), "4+")
		}
	default:
		if property.Unit < 1 {
			sl.ReportError(property.Unit, "unit", "Unit", "property_type", "")
		}

	}
}

func LoanApplyPropertyEntry(sl validator.StructLevel) {
	property, ok := sl.Current().Interface().(types.LoanApplyPropertyEntry)
	if !ok {
		sl.ReportError(property, "property", "Property", "property_type", "xx")
		return
	}
	switch property.PropertyType {
	case string(constants.PropertyTypeSingleFamilyResidence):
		if property.Unit != 1 {
			sl.ReportError(property.Unit, "unit", "Unit", string(constants.PropertyTypeSingleFamilyResidence), "1")
		}
	case string(constants.PropertyTypeTwoUnit):
		if property.Unit != 2 {
			sl.ReportError(property.Unit, "unit", "Unit", string(constants.PropertyTypeTwoUnit), "2")
		}
	case string(constants.PropertyTypeThreeUnit):
		if property.Unit != 3 {
			sl.ReportError(property.Unit, "unit", "Unit", string(constants.PropertyTypeThreeUnit), "3")
		}
	case string(constants.PropertyTypeFourUnit):
		if property.Unit != 4 {
			sl.ReportError(property.Unit, "unit", "Unit", string(constants.PropertyTypeFourUnit), "4")
		}
	case string(constants.PropertyTypeFourPlusMulti):
		if property.Unit < 4 {
			sl.ReportError(property.Unit, "unit", "Unit", string(constants.PropertyTypeFourPlusMulti), "4+")
		}
	default:
		if property.Unit < 1 {
			sl.ReportError(property.Unit, "unit", "Unit", "property_type", "")
		}

	}
}

// LoanPropertyEntry
func LoanPropertyEntry(sl validator.StructLevel) {
	property, ok := sl.Current().Interface().(types.LoanApplyPropertyEntry)
	if !ok {
		sl.ReportError(property, "property", "Property", "property_type", "xx")
		return
	}
	switch property.PropertyType {
	case string(constants.PropertyTypeSingleFamilyResidence):
		if property.Unit != 1 {
			sl.ReportError(property.Unit, "unit", "Unit", string(constants.PropertyTypeSingleFamilyResidence), "1")
		}
	case string(constants.PropertyTypeTwoUnit):
		if property.Unit != 2 {
			sl.ReportError(property.Unit, "unit", "Unit", string(constants.PropertyTypeTwoUnit), "2")
		}
	case string(constants.PropertyTypeThreeUnit):
		if property.Unit != 3 {
			sl.ReportError(property.Unit, "unit", "Unit", string(constants.PropertyTypeThreeUnit), "3")
		}
	case string(constants.PropertyTypeFourUnit):
		if property.Unit != 4 {
			sl.ReportError(property.Unit, "unit", "Unit", string(constants.PropertyTypeFourUnit), "4")
		}
	case string(constants.PropertyTypeFourPlusMulti):
		if property.Unit < 4 {
			sl.ReportError(property.Unit, "unit", "Unit", string(constants.PropertyTypeFourPlusMulti), "4+")
		}
	default:
		if property.Unit < 1 {
			sl.ReportError(property.Unit, "unit", "Unit", "property_type", "")
		}

	}
}

package validate_rules

import (
	"github.com/duke-git/lancet/v2/convertor"
	"github.com/go-playground/validator/v10"

	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/common/constants"
	"zeebee-backend/common/timex"
)

var PropertyUnitTag = "invalid_unit_for_property_type"
var PropertyUnitMsgTpl = "{0} property type requires unit value of {1}"

// LoanPreApprovalRequestProperty
func LoanPreApprovalRequestProperty(sl validator.StructLevel) {
	property, ok := sl.Current().Interface().(types.LoanPreApprovalRequestProperty)
	if !ok {
		sl.ReportError(property, "property", "Property", "property_type", "Invalid Property")
		return
	}
	switch property.PropertyType {
	case string(constants.PropertyTypeSingleFamilyResidence):
		if property.Unit != 1 {
			sl.ReportError(property.Unit, property.PropertyType, "Unit", PropertyUnitTag, convertor.ToString(property.Unit))
		}
	case string(constants.PropertyTypeTwoUnit):
		if property.Unit != 2 {
			sl.ReportError(property.Unit, property.PropertyType, "Unit", PropertyUnitTag, string(property.Unit))
		}
	case string(constants.PropertyTypeThreeUnit):
		if property.Unit != 3 {
			sl.ReportError(property.Unit, property.PropertyType, "Unit", PropertyUnitTag, convertor.ToString(property.Unit))
		}
	case string(constants.PropertyTypeFourUnit):
		if property.Unit != 4 {
			sl.ReportError(property.Unit, property.PropertyType, "Unit", PropertyUnitTag, convertor.ToString(property.Unit))
		}
	case string(constants.PropertyTypeFourPlusMulti):
		if property.Unit < 4 {
			sl.ReportError(property.Unit, property.PropertyType, "Unit", PropertyUnitTag, convertor.ToString(property.Unit))
		}
	default:
		if property.Unit < 1 {
			sl.ReportError(property.Unit, property.PropertyType, "Unit", PropertyUnitTag, " >= 1")
		}
	}
}

// LoanPropertyEntry
func LoanPropertyEntry(sl validator.StructLevel) {
	val := sl.Current().Interface()
	property, ok := val.(types.LoanPropertyEntry)
	if !ok {
		sl.ReportError(property, "property", "Property", "property_type", "Invalid Property")
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

func PropertyCreateRequest(sl validator.StructLevel) {
	val := sl.Current().Interface()
	property, ok := val.(types.PropertyCreateRequest)
	if !ok {
		sl.ReportError(property, "property", "Property", "property_type", "Invalid Property")
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

func PropertyUpdateRequest(sl validator.StructLevel) {
	val := sl.Current().Interface()
	property, ok := val.(types.PropertyUpdateRequest)
	if !ok {
		sl.ReportError(property, "property", "Property", "property_type", "Invalid Property")
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

// LoanPropertyApplyEntry
func LoanPropertyApplyEntry(sl validator.StructLevel) {
	var loanPurpose string
	var loanType string
	val := sl.Current().Interface()
	parent := sl.Parent().Interface()
	property, ok := val.(types.LoanPropertyApplyEntry)
	if !ok {
		sl.ReportError(property, "property", "Property", "property_type", "Invalid Loan Apply Property")
		return
	}
	switch req := parent.(type) {
	case types.LoanApplyRequest:
		// req := parent.(types.LoanApplyRequest)
		loanPurpose = req.LoanPurpose
		loanType = req.LoanType

	default:
		sl.ReportError(property, "property", "Property", "property_type", "Invalid Loan Apply Property")
		return
	}

	if loanPurpose != string(constants.LoanPurposePurchase) {
		if len(property.ExistingLenderCompany) == 0 {
			sl.ReportError(property.ExistingLenderCompany, "existing lender company", "ExistingLenderCompany", "required", "property.ExistingLenderCompany")
			return
		}
		if property.ExistingRate < 0 {
			sl.ReportError(property.ExistingRate, "existing rate", "ExistingRate", "gt", convertor.ToString(property.ExistingRate))
			return
		}
		if property.ExistingRate > 100 {
			sl.ReportError(property.ExistingRate, "existing rate", "ExistingRate", "lt", convertor.ToString(property.ExistingRate))
			return
		}

		if property.EstPropertyValueAmount <= 0 {
			sl.ReportError(property.EstPropertyValueAmount, "est property value amount", "EstPropertyValueAmount", "gt", convertor.ToString(property.EstPropertyValueAmount))
			return
		}

		_, err := timex.FormateTimeToStrWithUsStyle(property.ExistingMaturityDate)
		if err != nil {
			sl.ReportError(property.ExistingMaturityDate, "existing maturity date", "ExistingMaturityDate", "datetime", property.ExistingMaturityDate)
			return
		}

		_, err = timex.FormateTimeToStrWithUsStyle(property.PurchaseDate)
		if err != nil {
			sl.ReportError(property.ExistingMaturityDate, "purchase date", "PurchaseDate", "datetime", property.ExistingMaturityDate)
			return
		}

	}

	if loanPurpose == string(constants.LoanPurposePurchase) {
		_, err := timex.FormateTimeToStrWithUsStyle(property.ScheduledCloseEscrowDate)
		if err != nil {
			sl.ReportError(property.ScheduledCloseEscrowDate, "existing maturity date", "ExistingMaturityDate", "datetime", property.ScheduledCloseEscrowDate)
			return
		}
	}

	if loanType == string(constants.LoanTypeRehab) || loanType == string(constants.LoanTypeConstruction) {
		if property.Budget == nil {
			sl.ReportError(property.Budget, "budget", "Budget", "required", "Property Budget")
			return
		}
	}

	if loanType == string(constants.LoanTypeRental) {
		if property.EstYearlyRentalIncomeAmount <= 0 {
			sl.ReportError(property.EstYearlyRentalIncomeAmount, "est yearly rental income", "EstYearlyRentalIncome", "gt", convertor.ToString(property.EstYearlyRentalIncomeAmount))
			return
		}
		if property.YearlyLiabilityExpenseAmount <= 0 {
			sl.ReportError(property.YearlyLiabilityExpenseAmount, "yearly liability expense", "YearlyLiabilityExpense", "gt", convertor.ToString(property.YearlyLiabilityExpenseAmount))
			return
		}
		if property.Dscr <= 0 {
			sl.ReportError(property.Dscr, "dscr", "Dscr", "gt", convertor.ToString(property.Dscr))
			return
		}
	}

}

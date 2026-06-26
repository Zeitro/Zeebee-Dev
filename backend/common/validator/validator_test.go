package validator

import (
	"testing"

	"github.com/duke-git/lancet/v2/convertor"
	"github.com/go-playground/validator/v10"

	"zeebee-backend/domains/zeebee/constants"
)

var unitTag = "unit"
var unitErrorMsgTpl = "{0} property type requires unit value of {1}"

// TestAddressEntry represents a test address structure
type TestAddressEntry struct {
	State    string `json:"state,optional" validate:"required"`
	City     string `json:"city,optional" validate:"required"`
	Address1 string `json:"address1,optional" validate:"required"`
	Address2 string `json:"address2,optional"`
	ZipCode  string `json:"zip_code,optional"`
}

// TestLoanPreApprovalRequestProperty represents a test loan pre-approval request property structure
type TestLoanPreApprovalRequestProperty struct {
	Address       *TestAddressEntry `json:"address" validate:"required"`
	PropertyType  string            `json:"property_type" validate:"required,oneof='Single Family Residence' '2 unit' '3 unit' '4 unit' '4+ Multi' 'Townhouse' 'PUD/Condominium' 'Mixed Use' 'Commercial'"`
	IsPermitReady bool              `json:"is_permit_ready"`
	Unit          int32             `json:"unit" validate:"gte=1,lte=100"`
}

// LoanPreApprovalRequestPropertyRule validates the loan pre-approval request property
func LoanPreApprovalRequestPropertyRule(sl validator.StructLevel) {
	property, ok := sl.Current().Interface().(TestLoanPreApprovalRequestProperty)
	if !ok {
		sl.ReportError(property, "property", "Property", unitTag, "")
		return
	}
	switch property.PropertyType {
	case string(constants.PropertyTypeSingleFamilyResidence):
		if property.Unit != 1 {
			sl.ReportError(property.Unit, property.PropertyType, "Unit", unitTag, convertor.ToString(property.Unit))
		}
	case string(constants.PropertyTypeTwoUnit):
		if property.Unit != 2 {
			sl.ReportError(property.Unit, property.PropertyType, "Unit", unitTag, convertor.ToString(property.Unit))
		}
	case string(constants.PropertyTypeThreeUnit):
		if property.Unit != 3 {
			sl.ReportError(property.Unit, property.PropertyType, "Unit", unitTag, convertor.ToString(property.Unit))
		}
	case string(constants.PropertyTypeFourUnit):
		if property.Unit != 4 {
			sl.ReportError(property.Unit, property.PropertyType, "Unit", unitTag, convertor.ToString(property.Unit))
		}
	case string(constants.PropertyTypeFourPlusMulti):
		if property.Unit < 4 {
			sl.ReportError(property.Unit, property.PropertyType, "Unit", unitTag, convertor.ToString(property.Unit))
		}
	default:
		if property.Unit < 1 {
			sl.ReportError(property.Unit, property.PropertyType, "Unit", unitTag, string(property.Unit))
		}
	}
}

// TestValidator_LoanPreApprovalRequestPropertyValidation tests the LoanPreApprovalRequestPropertyRule validation
func TestValidator_LoanPreApprovalRequestPropertyValidation(t *testing.T) {
	// Create validator instance
	v, err := NewValidator(LocalEn)
	if err != nil {
		t.Fatalf("fail to create validator: %v", err)
	}

	// Register LoanPreApprovalRequestProperty rule
	v.RegisterStructRule(LoanPreApprovalRequestPropertyRule, TestLoanPreApprovalRequestProperty{})

	// Create test cases covering all rules
	testCases := []struct {
		name        string
		property    TestLoanPreApprovalRequestProperty
		expectError bool
	}{
		// Test correct values for all PropertyType types
		{
			name: "SingleFamilyResidence_CorrectUnit",
			property: TestLoanPreApprovalRequestProperty{
				Address: &TestAddressEntry{
					State:    "CA",
					City:     "San Francisco",
					Address1: "123 Main St",
				},
				PropertyType: string(constants.PropertyTypeSingleFamilyResidence),
				Unit:         1,
			},
			expectError: false,
		},
		{
			name: "TwoUnit_CorrectUnit",
			property: TestLoanPreApprovalRequestProperty{
				Address: &TestAddressEntry{
					State:    "CA",
					City:     "San Francisco",
					Address1: "123 Main St",
				},
				PropertyType: string(constants.PropertyTypeTwoUnit),
				Unit:         2,
			},
			expectError: false,
		},
		{
			name: "ThreeUnit_CorrectUnit",
			property: TestLoanPreApprovalRequestProperty{
				Address: &TestAddressEntry{
					State:    "CA",
					City:     "San Francisco",
					Address1: "123 Main St",
				},
				PropertyType: string(constants.PropertyTypeThreeUnit),
				Unit:         3,
			},
			expectError: false,
		},
		{
			name: "FourUnit_CorrectUnit",
			property: TestLoanPreApprovalRequestProperty{
				Address: &TestAddressEntry{
					State:    "CA",
					City:     "San Francisco",
					Address1: "123 Main St",
				},
				PropertyType: string(constants.PropertyTypeFourUnit),
				Unit:         4,
			},
			expectError: false,
		},
		{
			name: "FourPlusMulti_CorrectUnit",
			property: TestLoanPreApprovalRequestProperty{
				Address: &TestAddressEntry{
					State:    "CA",
					City:     "San Francisco",
					Address1: "123 Main St",
				},
				PropertyType: string(constants.PropertyTypeFourPlusMulti),
				Unit:         5,
			},
			expectError: false,
		},
		{
			name: "FourPlusMulti_EdgeCaseUnit4",
			property: TestLoanPreApprovalRequestProperty{
				Address: &TestAddressEntry{
					State:    "CA",
					City:     "San Francisco",
					Address1: "123 Main St",
				},
				PropertyType: string(constants.PropertyTypeFourPlusMulti),
				Unit:         4,
			},
			expectError: false,
		},
		{
			name: "Townhouse_CorrectUnit",
			property: TestLoanPreApprovalRequestProperty{
				Address: &TestAddressEntry{
					State:    "CA",
					City:     "San Francisco",
					Address1: "123 Main St",
				},
				PropertyType: string(constants.PropertyTypeTownhouse),
				Unit:         1,
			},
			expectError: false,
		},
		{
			name: "PUDCondominium_CorrectUnit",
			property: TestLoanPreApprovalRequestProperty{
				Address: &TestAddressEntry{
					State:    "CA",
					City:     "San Francisco",
					Address1: "123 Main St",
				},
				PropertyType: string(constants.PropertyTypePUDCondominium),
				Unit:         1,
			},
			expectError: false,
		},
		{
			name: "MixedUse_CorrectUnit",
			property: TestLoanPreApprovalRequestProperty{
				Address: &TestAddressEntry{
					State:    "CA",
					City:     "San Francisco",
					Address1: "123 Main St",
				},
				PropertyType: string(constants.PropertyTypeMixedUse),
				Unit:         1,
			},
			expectError: false,
		},
		{
			name: "Commercial_CorrectUnit",
			property: TestLoanPreApprovalRequestProperty{
				Address: &TestAddressEntry{
					State:    "CA",
					City:     "San Francisco",
					Address1: "123 Main St",
				},
				PropertyType: string(constants.PropertyTypeCommercial),
				Unit:         1,
			},
			expectError: false,
		},
		// Test incorrect Unit values
		{
			name: "SingleFamilyResidence_WrongUnit",
			property: TestLoanPreApprovalRequestProperty{
				Address: &TestAddressEntry{
					State:    "CA",
					City:     "San Francisco",
					Address1: "123 Main St",
				},
				PropertyType: string(constants.PropertyTypeSingleFamilyResidence),
				Unit:         2,
			},
			expectError: true,
		},
		{
			name: "TwoUnit_WrongUnit",
			property: TestLoanPreApprovalRequestProperty{
				Address: &TestAddressEntry{
					State:    "CA",
					City:     "San Francisco",
					Address1: "123 Main St",
				},
				PropertyType: string(constants.PropertyTypeTwoUnit),
				Unit:         3,
			},
			expectError: true,
		},
		{
			name: "ThreeUnit_WrongUnit",
			property: TestLoanPreApprovalRequestProperty{
				Address: &TestAddressEntry{
					State:    "CA",
					City:     "San Francisco",
					Address1: "123 Main St",
				},
				PropertyType: string(constants.PropertyTypeThreeUnit),
				Unit:         2,
			},
			expectError: true,
		},
		{
			name: "FourUnit_WrongUnit",
			property: TestLoanPreApprovalRequestProperty{
				Address: &TestAddressEntry{
					State:    "CA",
					City:     "San Francisco",
					Address1: "123 Main St",
				},
				PropertyType: string(constants.PropertyTypeFourUnit),
				Unit:         3,
			},
			expectError: true,
		},
		{
			name: "FourPlusMulti_WrongUnit",
			property: TestLoanPreApprovalRequestProperty{
				Address: &TestAddressEntry{
					State:    "CA",
					City:     "San Francisco",
					Address1: "123 Main St",
				},
				PropertyType: string(constants.PropertyTypeFourPlusMulti),
				Unit:         3,
			},
			expectError: true,
		},
		// Test boundary cases for unknown PropertyType
		{
			name: "UnknownPropertyType_NegativeUnit",
			property: TestLoanPreApprovalRequestProperty{
				Address: &TestAddressEntry{
					State:    "CA",
					City:     "San Francisco",
					Address1: "123 Main St",
				},
				PropertyType: "UnknownType",
				Unit:         -1,
			},
			expectError: true,
		},
		{
			name: "UnknownPropertyType_PositiveUnit",
			property: TestLoanPreApprovalRequestProperty{
				Address: &TestAddressEntry{
					State:    "CA",
					City:     "San Francisco",
					Address1: "123 Main St",
				},
				PropertyType: "UnknownType",
				Unit:         5,
			},
			expectError: true, // Should be an error because PropertyType is not in the allowed range
		},
		// Test missing required fields
		{
			name: "MissingAddress",
			property: TestLoanPreApprovalRequestProperty{
				Address:      nil,
				PropertyType: string(constants.PropertyTypeSingleFamilyResidence),
				Unit:         1,
			},
			expectError: true,
		},
		{
			name: "MissingPropertyType",
			property: TestLoanPreApprovalRequestProperty{
				Address: &TestAddressEntry{
					State:    "CA",
					City:     "San Francisco",
					Address1: "123 Main St",
				},
				PropertyType: "",
				Unit:         1,
			},
			expectError: true,
		},
		{
			name: "InvalidPropertyType",
			property: TestLoanPreApprovalRequestProperty{
				Address: &TestAddressEntry{
					State:    "CA",
					City:     "San Francisco",
					Address1: "123 Main St",
				},
				PropertyType: "InvalidType",
				Unit:         1,
			},
			expectError: true,
		},
		// Test Unit out of range
		{
			name: "UnitTooSmall",
			property: TestLoanPreApprovalRequestProperty{
				Address: &TestAddressEntry{
					State:    "CA",
					City:     "San Francisco",
					Address1: "123 Main St",
				},
				PropertyType: string(constants.PropertyTypeSingleFamilyResidence),
				Unit:         0,
			},
			expectError: true,
		},
		{
			name: "UnitTooLarge",
			property: TestLoanPreApprovalRequestProperty{
				Address: &TestAddressEntry{
					State:    "CA",
					City:     "San Francisco",
					Address1: "123 Main St",
				},
				PropertyType: string(constants.PropertyTypeSingleFamilyResidence),
				Unit:         101,
			},
			expectError: true,
		},
		// Test Address field validation
		{
			name: "MissingAddressFields",
			property: TestLoanPreApprovalRequestProperty{
				Address: &TestAddressEntry{
					State:    "",
					City:     "",
					Address1: "",
				},
				PropertyType: string(constants.PropertyTypeSingleFamilyResidence),
				Unit:         1,
			},
			expectError: true,
		},
	}

	// Execute test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := v.Validate(tc.property)

			if tc.expectError && err == nil {
				t.Errorf("Expected error but got nil")
			}

			if !tc.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

// TestPerson represents a person with embedded property for testing
type TestPerson struct {
	Name     string                             `json:"name" validate:"required"`
	Property TestLoanPreApprovalRequestProperty `json:"property" validate:"required"`
}

// TestValidator_EmbeddedPropertyValidation tests if embedded TestLoanPreApprovalRequestProperty applies LoanPreApprovalRequestPropertyRule correctly
func TestValidator_EmbeddedPropertyValidation(t *testing.T) {
	// Create validator instance
	v, err := NewValidator(LocalEn)
	if err != nil {
		t.Fatalf("fail to create validator: %v", err)
	}

	// Register LoanPreApprovalRequestProperty rule
	v.RegisterStructRule(LoanPreApprovalRequestPropertyRule, TestLoanPreApprovalRequestProperty{})

	// Create test cases
	testCases := []struct {
		name          string
		person        TestPerson
		expectError   bool
		errorContains string
	}{
		{
			name: "ValidPersonWithCorrectProperty",
			person: TestPerson{
				Name: "John Doe",
				Property: TestLoanPreApprovalRequestProperty{
					Address: &TestAddressEntry{
						State:    "CA",
						City:     "San Francisco",
						Address1: "123 Main St",
					},
					PropertyType: string(constants.PropertyTypeSingleFamilyResidence),
					Unit:         1,
				},
			},
			expectError: false,
		},
		{
			name: "ValidPersonWithCorrectProperty2",
			person: TestPerson{
				Name: "Jane Smith",
				Property: TestLoanPreApprovalRequestProperty{
					Address: &TestAddressEntry{
						State:    "NY",
						City:     "New York",
						Address1: "456 Broadway",
					},
					PropertyType: string(constants.PropertyTypeFourUnit),
					Unit:         4,
				},
			},
			expectError: false,
		},
		{
			name: "InvalidPerson_EmptyName",
			person: TestPerson{
				Name: "",
				Property: TestLoanPreApprovalRequestProperty{
					Address: &TestAddressEntry{
						State:    "CA",
						City:     "San Francisco",
						Address1: "123 Main St",
					},
					PropertyType: string(constants.PropertyTypeSingleFamilyResidence),
					Unit:         1,
				},
			},
			expectError:   true,
			errorContains: "required",
		},
		{
			name: "InvalidPerson_MissingProperty",
			person: TestPerson{
				Name:     "John Doe",
				Property: TestLoanPreApprovalRequestProperty{},
			},
			expectError: true,
		},
		{
			name: "InvalidPerson_WrongPropertyUnit",
			person: TestPerson{
				Name: "John Doe",
				Property: TestLoanPreApprovalRequestProperty{
					Address: &TestAddressEntry{
						State:    "CA",
						City:     "San Francisco",
						Address1: "123 Main St",
					},
					PropertyType: string(constants.PropertyTypeTwoUnit),
					Unit:         3, // Wrong Unit value, should be 2
				},
			},
			expectError: true,
		},
		{
			name: "InvalidPerson_MissingAddress",
			person: TestPerson{
				Name: "John Doe",
				Property: TestLoanPreApprovalRequestProperty{
					Address:      nil,
					PropertyType: string(constants.PropertyTypeSingleFamilyResidence),
					Unit:         1,
				},
			},
			expectError: true,
		},
		{
			name: "InvalidPerson_InvalidPropertyType",
			person: TestPerson{
				Name: "John Doe",
				Property: TestLoanPreApprovalRequestProperty{
					Address: &TestAddressEntry{
						State:    "CA",
						City:     "San Francisco",
						Address1: "123 Main St",
					},
					PropertyType: "InvalidType",
					Unit:         1,
				},
			},
			expectError: true,
		},
		{
			name: "InvalidPerson_MissingAddressFields",
			person: TestPerson{
				Name: "John Doe",
				Property: TestLoanPreApprovalRequestProperty{
					Address: &TestAddressEntry{
						State:    "",
						City:     "",
						Address1: "",
					},
					PropertyType: string(constants.PropertyTypeSingleFamilyResidence),
					Unit:         1,
				},
			},
			expectError: true,
		},
	}

	// Execute test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := v.Validate(tc.person)

			if tc.expectError && err == nil {
				t.Errorf("Expected error but got nil")
			}

			if !tc.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestValidator_ValidateRuleMsg(t *testing.T) {
	// Create validator instance
	v, err := NewValidator(LocalEn)
	if err != nil {
		t.Fatalf("fail to create validator: %v", err)
	}

	// Register LoanPreApprovalRequestProperty rule
	v.RegisterStructRule(LoanPreApprovalRequestPropertyRule, TestLoanPreApprovalRequestProperty{})
	_ = v.RegisterTansMsg(unitTag, unitErrorMsgTpl)
	noName := TestPerson{
		Name: "",
		Property: TestLoanPreApprovalRequestProperty{
			Address: &TestAddressEntry{
				State:    "CA",
				City:     "San Francisco",
				Address1: "123 Main St",
			},
			PropertyType: string(constants.PropertyTypeSingleFamilyResidence),
			Unit:         1,
		},
	}
	err = v.Validate(noName)
	t.Logf("no name validate error msg:%v\n", err)

	noPropertyType := TestPerson{
		Name: "John Doe",
		Property: TestLoanPreApprovalRequestProperty{
			Address: &TestAddressEntry{
				State:    "CA",
				City:     "San Francisco",
				Address1: "123 Main St",
			},
			PropertyType: "",
			Unit:         1,
		},
	}
	err = v.Validate(noPropertyType)
	t.Logf("no property validate error msg:%v\n", err)

	invalidUnitType := TestPerson{
		Name: "John Doe",
		Property: TestLoanPreApprovalRequestProperty{
			Address: &TestAddressEntry{
				State:    "CA",
				City:     "San Francisco",
				Address1: "123 Main St",
			},
			PropertyType: string(constants.PropertyTypeSingleFamilyResidence),
			Unit:         2,
		},
	}
	err = v.Validate(invalidUnitType)
	t.Logf("invalid uint validate error msg:%v\n", err)
}

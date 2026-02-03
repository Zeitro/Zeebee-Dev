package data

import (
	"database/sql/driver"
	"encoding/json"

	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
)

// TemplateVariableGroup represents a group of variables
type TemplateVariableGroup struct {
	Name      string              `json:"name"`
	Variables []*TemplateVariable `json:"variables"`
}

// Scan implements sql.Scanner interface for TemplateVariableGroup
func (tvg *TemplateVariableGroup) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errorx.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, tvg)
}

// Value implements driver.Valuer interface for TemplateVariableGroup
func (tvg TemplateVariableGroup) Value() (driver.Value, error) {
	return json.Marshal(tvg)
}

// TemplateVariable represents a variable definition in the template schema
type TemplateVariable struct {
	Label      string                           `json:"label"`
	Type       constants.TemplateVariableType   `json:"type"` // String,Array,Object
	Key        string                           `json:"key"`
	IsEditable bool                             `json:"is_editable"`
	IsRequired bool                             `json:"is_required"`
	Source     constants.TemplateVariableSource `json:"source"` // Fixed,Dynamic,Template
	FixedValue string                           `json:"fixed_value"`
	Format     constants.TemplateVariableFormat `json:"format,omitempty"`     // Email,Phone,Currency
	Properties map[string]*TemplateVariable     `json:"properties,omitempty"` // for object type
	Items      []*TemplateVariable              `json:"items,omitempty"`      // for array type
}

// Scan implements sql.Scanner interface for TemplateVariable
func (tv *TemplateVariable) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errorx.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, tv)
}

// Value implements driver.Valuer interface for TemplateVariable
func (tv TemplateVariable) Value() (driver.Value, error) {
	return json.Marshal(tv)
}

// Validate checks if the TemplateVariable has valid configuration
func (tv *TemplateVariable) Validate() error {
	validTypes := map[constants.TemplateVariableType]bool{
		constants.TemplateVariableTypeString: true,
		constants.TemplateVariableTypeArray:  true,
		constants.TemplateVariableTypeObject: true,
	}
	if !validTypes[tv.Type] {
		return errorx.Errorf("invalid variable type: %v", tv.Type)
	}

	// If IsRequired is true and Source is Fixed or Template, FixedValue must have a value
	if tv.IsRequired && (tv.Source == constants.TemplateVariableSourceFixed || tv.Source == constants.TemplateVariableSourceTemplate) {
		if tv.FixedValue == "" {
			return errorx.Errorf("FixedValue is required when IsRequired is true and Source is %v", tv.Source)
		}
	}

	// Validate nested structures based on type
	if tv.Type == constants.TemplateVariableTypeObject && tv.Properties != nil {
		for _, prop := range tv.Properties {
			if err := prop.Validate(); err != nil {
				return err
			}
		}
	} else if tv.Type == constants.TemplateVariableTypeArray && tv.Items != nil {
		for _, item := range tv.Items {
			if err := item.Validate(); err != nil {
				return err
			}
		}
	}

	return nil
}

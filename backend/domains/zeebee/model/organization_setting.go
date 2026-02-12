package model

// OrganizationSetting represents an organization's customized template
type OrganizationSetting struct {
	Base
	OrganizationId int64                  `json:"organization_id" gorm:"column:organization_id;default:0;index"`
	Purpose        string                 `json:"purpose" gorm:"column:purpose;type:varchar(100);not null;default:''"`
	Enabled        bool                   `json:"enabled" gorm:"column:enabled;default:false"`
	Setting        map[string]interface{} `json:"setting" gorm:"column:setting;type:jsonb;not null;default:'{}';serializer:json"`

	// Relationships
	Organization *Organization `json:"organization,omitempty" gorm:"foreignKey:OrganizationId;references:Id"`
}

// TableName specifies the table name for OrganizationTemplate
func (OrganizationSetting) TableName() string {
	return "organization_settings"
}

// RemoveRelationships removes relationship data to avoid circular references
func (ot *OrganizationSetting) RemoveRelationships() {
	// ot.BaseTemplate = nil
	ot.Organization = nil
}

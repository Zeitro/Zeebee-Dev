package model

type LenderEmployee struct {
	Base

	UserId                 int64  `gorm:"column:user_id;not null;default:0" json:"user_id"`
	DefaultLenderCompanyId int64  `gorm:"column:default_lender_company_id;not null;default:0" json:"default_lender_company_id"`
	FirstName              string `gorm:"column:first_name;size:100;not null;default:''" json:"first_name"`
	LastName               string `gorm:"column:last_name;size:100;not null;default:''" json:"last_name"`
	FullName               string `gorm:"column:full_name;size:100;not null;default:''" json:"full_name"`
	//  relationships
	User *User `json:"user" gorm:"foreignKey:UserId;references:Id"`
	// LenderCompany   *LenderCompany   `json:"lender_company" gorm:"foreignKey:LenderCompanyId;references:Id"`
	LenderCompanies []*LenderCompany `json:"lender_companies" gorm:"many2many:lender_company_and_employees;joinForeignKey:lender_employee_id;joinReferences:lender_company_id"`
}

func (m *LenderEmployee) TableName() string {
	return "lender_employees"
}

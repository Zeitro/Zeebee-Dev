package dto

type LoanBorrowerPageRequestDto struct {
	BasePageRequestDto
	Filter LoanBorrowerPageRequestFilter `json:"filter"`
	Sorter BasePageSorter                `json:"sorter"`
}

type LoanBorrowerPageRequestFilter struct {
	UserId         int64 `json:"user_id,omitempty"`
	OrganizationId int64 `json:"organization_id"`
	// BorrowerOrganizationId int64    `json:"borrower_organization_id"`
	Stage         string   `json:"stage,omitempty"`
	Keyword       string   `json:"keyword,omitempty"`
	KeywordType   string   `json:"keyword_type,omitempty"`
	LoanPurpose   []string `json:"loan_purpose,omitempty"`
	LoanType      []string `json:"loan_type,omitempty"`
	Lender        []string `json:"lender,omitempty"`
	PendingStatus []string `json:"pending_status,omitempty"`
}

type LoanLenderPageRequestDto struct {
	BasePageRequestDto
	Filter LoanLenderPageRequestFilter `json:"filter,omitempty"`
	Sorter BasePageSorter              `json:"sorter"`
}

type LoanLenderPageRequestFilter struct {
	OrganizationId int64    `json:"organization_id"`
	Stage          string   `json:"stage,omitempty"`
	Keyword        string   `json:"keyword,omitempty"`
	KeywordType    string   `json:"keyword_type,omitempty"`
	LoanPurpose    []string `json:"loan_purpose,omitempty"`
	LoanType       []string `json:"loan_type,omitempty"`
	PendingStatus  []string `json:"pending_status,omitempty"`
}

type LoanAdminPageRequestDto struct {
	BasePageRequestDto
	Filter LoanAdminPageRequestFilter `json:"filter,omitempty"`
	Sorter BasePageSorter             `json:"sorter"`
}

type LoanAdminPageRequestFilter struct {
	// OrganizationId int64    `json:"organization_id"`
	Stage         string   `json:"stage,omitempty"`
	Keyword       string   `json:"keyword,omitempty"`
	KeywordType   string   `json:"keyword_type,omitempty"`
	LoanPurpose   []string `json:"loan_purpose,omitempty"`
	LoanType      []string `json:"loan_type,omitempty"`
	PendingStatus []string `json:"pending_status,omitempty"`
	Lender        []string `json:"lender,omitempty"`
}

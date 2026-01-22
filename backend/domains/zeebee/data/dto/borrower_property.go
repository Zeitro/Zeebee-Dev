package dto

type BorrowerPropertyPageRequestDto struct {
	BasePageRequestDto
	Filter BorrowerPropertyPageRequestFilter `json:"filter"`
	Sorter BasePageSorter                    `json:"sorter"`
}

type BorrowerPropertyPageRequestFilter struct {
	UserId         int64    `json:"user_id"`
	OrganizationId int64    `json:"organization_id"`
	Keyword        string   `json:"keyword,omitempty"`
	Applicant      []string `json:"applicant,omitempty"`
}

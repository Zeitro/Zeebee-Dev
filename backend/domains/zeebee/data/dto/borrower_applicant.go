package dto

type BorrowerApplicantBorrowerPageRequestDto struct {
	BasePageRequestDto
	Filter BorrowerApplicantBorrowerPageRequestFilter `json:"filter"`
	Sorter BasePageSorter
}

type BorrowerApplicantBorrowerPageRequestFilter struct {
	UserId         int64  `json:"user_id"`
	OrganizationId int64  `json:"organization_id"`
	Keyword        string `json:"keyword"`
}

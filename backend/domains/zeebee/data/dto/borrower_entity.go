package dto

type BorrowerEntityPageRequestDto struct {
	BasePageRequestDto
	Filter BorrowerEntityPageRequestFilter `json:"filter"`
	Sorter BasePageSorter                  `json:"sorter"`
}

type BorrowerEntityPageRequestFilter struct {
	UserId         int64  `json:"user_id"`
	OrganizationId int64  `json:"organization_id"`
	Keyword        string `json:"keyword"`
}

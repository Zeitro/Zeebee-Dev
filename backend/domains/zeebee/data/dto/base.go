package dto

type BasePageRequestDto struct {
	Page  int64 `json:"page"`
	Limit int64 `json:"limit"`
}

type BasePageFilter struct {
	Keyword string `json:"keyword"`
}

type BasePageSorter struct {
	Column string `json:"column"`
	Sort   string `json:"sort"`
}

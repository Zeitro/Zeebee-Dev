package dto

type ProjectManagementTaskTypeSnapshotDto struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

type ProjectManagementTaskTypeDto struct {
	Id          int64
	Name        string
	Description string
	Color       string
	SortOrder   int
}

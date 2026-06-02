package dto

type DocumentCreateDto struct {
	EntityType string `json:"entity_type"`
	EntityId   int64  `json:"entity_id"`
	ParentId   int64  `json:"parent_id"`
	Name       string `json:"name"`
}

type DocumentUploadFileDto struct {
	DocumentId int64  `json:"document_id"`
	FileName   string `json:"file_name"`
	FileUrl    string `json:"file_url"`
	FileSuffix string `json:"file_suffix"`
	OperatorId int64  `json:"operator_id"`
}

type DocumentDeleteFileDto struct {
	FileId     int64 `json:"file_id"`
	OperatorId int64 `json:"operator_id"`
}

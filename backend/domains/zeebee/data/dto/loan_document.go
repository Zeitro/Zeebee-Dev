package dto

import (
	"database/sql"
	"time"

	"zeebee-backend/domains/zeebee/constants"
)

type LoanDocumentEntityDto struct {
	Id             int64                          `json:"id"`
	LoanId         int64                          `json:"loan_id"`
	Category       constants.LoanDocumentCategory `json:"category"`
	FileType       constants.LoanDocumentFileType `json:"file_type"`
	Comment        string                         `json:"comment"`
	Status         constants.LoanDocumentStatus   `json:"status"`
	Files          []*FileEntityDto               `json:"files"`
	LastUploadedAt *sql.NullTime                  `json:"last_uploaded_at"`
	CreatedAt      time.Time                      `json:"created_at"`
	UpdatedAt      time.Time                      `json:"updated_at"`
}

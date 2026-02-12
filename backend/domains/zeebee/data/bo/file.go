package bo

import (
	"time"

	"zeebee-backend/domains/zeebee/constants"
)

type FileEntityBo struct {
	Id         int64                    `json:"id"`
	EntityType constants.FileEntityType `json:"entity_type"`
	EntityId   int64                    `json:"entity_id"`
	FileUrl    string                   `json:"file_url"`
	CreatedAt  time.Time                `json:"created_at"`
	UpdatedAt  time.Time                `json:"updated_at"`
}

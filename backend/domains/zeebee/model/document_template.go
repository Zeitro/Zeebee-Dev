package model

import "zeebee-backend/domains/zeebee/constants"

type DocumentTemplate struct {
	Base
	Name      string                       `gorm:"column:name;type:varchar(255);not null;default:''" json:"name"`
	Type      string                       `gorm:"column:type;type:varchar(20);not null;default:''" json:"type"`
	Structure DocumentTemplateStructureMap `gorm:"column:structure;type:json;not null;default:'{}';serializer:json" json:"structure"`
	Status    constants.CommonStatus       `gorm:"column:status;type:varchar(20);not null;default:''" json:"status"`
}

func (m DocumentTemplate) TableName() string {
	return "document_templates"
}

type DocumentTemplateStructureMap map[string][]DocumentTemplateStructure

type DocumentTemplateStructure struct {
	Name            string                      `json:"name"`
	SortOrder       int32                       `json:"sort_order"`
	IsRequired      bool                        `json:"is_required"`
	FileParseConfig *FileParseConfig            `json:"file_parse_config,omitempty"`
	// Tags classify the target node; copied into Document.metadata.tags at
	// transaction creation so manual_select selectors can scope files by tag.
	Tags     []string                    `json:"tags,omitempty"`
	Children []DocumentTemplateStructure `json:"children"`
}

// FileParseConfig declares how a document node is parsed. Rules is the v2
// multi-rule form; Parser is the legacy single-parser form normalized to one
// OnUpload rule by effectiveRules (kept for back-compat with existing rows).
type FileParseConfig struct {
	Parser string          `json:"parser,omitempty"`
	Rules  []FileParseRule `json:"rules,omitempty"`
}

// FileParseRule binds a parser to a trigger mode. Selector is required for
// ManualSelect (scopes files by tags); ignored otherwise.
type FileParseRule struct {
	Parser   string       `json:"parser"`             // Rpa | EscrowInfo
	Mode     string       `json:"mode"`               // OnUpload | ManualTree | ManualSelect
	Selector *TagSelector `json:"selector,omitempty"` // ManualSelect only
}

// TagSelector matches documents whose metadata.tags contain all listed tags (AND).
type TagSelector struct {
	Tags []string `json:"tags"`
}

// FileParseRule.Mode values.
const (
	FileParseModeOnUpload     = "OnUpload"
	FileParseModeManualTree   = "ManualTree"
	FileParseModeManualSelect = "ManualSelect"
)

// FileParseRule.Parser values.
const (
	FileParseParserRpa        = "Rpa"
	FileParseParserEscrowInfo = "EscrowInfo"
)

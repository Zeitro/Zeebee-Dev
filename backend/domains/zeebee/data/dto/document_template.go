package dto

import "zeebee-backend/domains/zeebee/constants"

type DocumentTemplateStructureMap map[string][]*DocumentTemplateStructure

type DocumentTemplateStructure struct {
	Name            string                           `json:"name"`
	SortOrder       int32                            `json:"sort_order"`
	IsRequired      bool                             `json:"is_required"`
	FileParseConfig *DocumentTemplateFileParseConfig `json:"file_parse_config,omitempty"`
	Tags            []string                         `json:"tags,omitempty"`
	Children        []*DocumentTemplateStructure     `json:"children,omitempty"`
}

type DocumentTemplateFileParseConfig struct {
	Rules []DocumentTemplateFileParseRule `json:"rules,omitempty"`
}

// FileParseRule is the lean, frontend-facing projection of model.FileParseRule.
// copierx maps it onto proto.FileParseRule by matching field names.
type DocumentTemplateFileParseRule struct {
	Parser      constants.DocumentTemplateParser  `json:"parser"`
	Source      constants.DocumentTemplateSource  `json:"source"`
	TriggeredBy constants.DocumentTemplateTrigger `json:"triggered_by"`
	Label       string                            `json:"label,omitempty"`
	Selector    *DocumentTemplateTagSelector      `json:"selector,omitempty"`
}

// TagSelector matches documents whose metadata.tags contain all listed tags.
type DocumentTemplateTagSelector struct {
	Tags []string `json:"tags,omitempty"`
}

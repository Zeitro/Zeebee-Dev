package handler

import (
	"encoding/json"

	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
)

// This file holds the rule/file-parse-config resolution the dispatcher (parent
// parse_file package) uses to decide which handler a request routes to. It is
// the only shared routing surface; per-parser collect/extract/arbitrate logic
// lives in each handler's own file (escrow.go / rpa.go).

// effectiveRules returns the rule list from a FileParseConfig. An absent or
// empty config yields no rules.
func effectiveRules(cfg *dto.DocumentTemplateFileParseConfig) []dto.DocumentTemplateFileParseRule {
	if cfg == nil {
		return nil
	}
	return cfg.Rules
}

// OnUploadRules returns only the rules that fire on file upload. Exported for
// the dispatcher to route an upload batch.
func OnUploadRules(cfg *dto.DocumentTemplateFileParseConfig) []dto.DocumentTemplateFileParseRule {
	var out []dto.DocumentTemplateFileParseRule
	for _, r := range effectiveRules(cfg) {
		if r.Source == constants.DocumentTemplateFileParseSourceSingleFile {
			out = append(out, r)
		}
	}
	return out
}

// ResolveRule returns the first rule matching parser, or ok=false. Used by the
// dispatcher (manual trigger) to pick the rule for an explicit parser.
func ResolveRule(cfg *dto.DocumentTemplateFileParseConfig, parser string) (dto.DocumentTemplateFileParseRule, bool) {
	for _, r := range effectiveRules(cfg) {
		if string(r.Parser) == parser {
			return r, true
		}
	}
	return dto.DocumentTemplateFileParseRule{}, false
}

// DocFileParseConfig reads file_parse_config from a document's metadata jsonb.
// Missing/unreadable metadata yields nil.
func DocFileParseConfig(doc *model.Document) *dto.DocumentTemplateFileParseConfig {
	if doc == nil || len(doc.Metadata) == 0 {
		return nil
	}
	var meta struct {
		FileParseConfig *dto.DocumentTemplateFileParseConfig `json:"file_parse_config"`
	}
	if err := json.Unmarshal(doc.Metadata, &meta); err != nil {
		return nil
	}
	return meta.FileParseConfig
}

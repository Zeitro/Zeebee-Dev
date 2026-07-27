package parse_file

import (
	"encoding/json"

	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
)

// Rule/file-parse-config ROUTING for the dispatcher. Only routing surface;
// per-parser logic lives in package handler. Trigger (when a rule fires) and
// Source (which files) are independent — upload-filtering is by TriggeredBy,
// never Source.

// effectiveRules returns the config's rule list; absent/empty → nil.
func effectiveRules(cfg *dto.DocumentTemplateFileParseConfig) []dto.DocumentTemplateFileParseRule {
	if cfg == nil {
		return nil
	}
	return cfg.Rules
}

// ruleMatchesTrigger is the single routing gate (FR-031): a rule runs only via
// the trigger it declares. Looks only at TriggeredBy, not Source.
func ruleMatchesTrigger(rule dto.DocumentTemplateFileParseRule, trigger Trigger) bool {
	want := constants.FileParseTriggeredTypeManual
	if trigger == TriggerOnUpload {
		want = constants.FileParseTriggeredTypeUpload
	}
	return rule.TriggeredBy == want
}

// onUploadRules returns rules that auto-fire on upload (gated by
// ruleMatchesTrigger(Upload), not Source).
func onUploadRules(cfg *dto.DocumentTemplateFileParseConfig) []dto.DocumentTemplateFileParseRule {
	var out []dto.DocumentTemplateFileParseRule
	for _, r := range effectiveRules(cfg) {
		if ruleMatchesTrigger(r, TriggerOnUpload) {
			out = append(out, r)
		}
	}
	return out
}

// resolveRule returns the first rule matching parser. Trigger-agnostic (a
// recovered job carries only its parser).
func resolveRule(cfg *dto.DocumentTemplateFileParseConfig, parser string) (dto.DocumentTemplateFileParseRule, bool) {
	for _, r := range effectiveRules(cfg) {
		if string(r.Parser) == parser {
			return r, true
		}
	}
	return dto.DocumentTemplateFileParseRule{}, false
}

// docFileParseConfig reads file_parse_config from document metadata jsonb; missing/unreadable → nil.
func docFileParseConfig(doc *model.Document) *dto.DocumentTemplateFileParseConfig {
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

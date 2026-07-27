package parse_file

import (
	"encoding/json"

	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
)

// This file holds the rule/file-parse-config ROUTING the dispatcher uses to
// decide which handler a request routes to and which rules fire on upload. It is
// the only routing surface; per-parser collect/extract/writeback logic lives in
// package handler.
//
// Trigger (when a rule fires) and Source (which files it scopes to) are
// independent dimensions: upload-filtering is by TriggeredBy, never by Source.

// effectiveRules returns the rule list from a FileParseConfig. An absent or empty
// config yields no rules.
func effectiveRules(cfg *dto.DocumentTemplateFileParseConfig) []dto.DocumentTemplateFileParseRule {
	if cfg == nil {
		return nil
	}
	return cfg.Rules
}

// ruleMatchesTrigger is the single routing gate (FR-031): a rule may run only
// via the trigger it declares. Trigger (when a rule fires) and Source (which
// files it scopes to) are orthogonal — this looks only at TriggeredBy, never at
// Source. Both dispatch paths resolve a rule by parser first, then apply this.
func ruleMatchesTrigger(rule dto.DocumentTemplateFileParseRule, trigger Trigger) bool {
	want := constants.FileParseTriggeredTypeManual
	if trigger == TriggerOnUpload {
		want = constants.FileParseTriggeredTypeUpload
	}
	return rule.TriggeredBy == want
}

// onUploadRules returns the rules that auto-fire on file upload — gated by
// ruleMatchesTrigger(Upload), not by Source. An Upload-triggered rule fires on
// upload regardless of whether its scope is SingleFile or Directory.
func onUploadRules(cfg *dto.DocumentTemplateFileParseConfig) []dto.DocumentTemplateFileParseRule {
	var out []dto.DocumentTemplateFileParseRule
	for _, r := range effectiveRules(cfg) {
		if ruleMatchesTrigger(r, TriggerOnUpload) {
			out = append(out, r)
		}
	}
	return out
}

// resolveRule returns the first rule matching parser (ok=false if none). Used by
// the manual dispatcher and by recovery. Trigger-agnostic: a recovered job only
// carries its parser, not its trigger, so this matches on parser alone.
func resolveRule(cfg *dto.DocumentTemplateFileParseConfig, parser string) (dto.DocumentTemplateFileParseRule, bool) {
	for _, r := range effectiveRules(cfg) {
		if string(r.Parser) == parser {
			return r, true
		}
	}
	return dto.DocumentTemplateFileParseRule{}, false
}

// docFileParseConfig reads file_parse_config from a document's metadata jsonb.
// Missing/unreadable metadata yields nil.
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

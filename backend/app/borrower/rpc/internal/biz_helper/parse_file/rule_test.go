package parse_file

import (
	"encoding/json"
	"testing"

	"gorm.io/datatypes"

	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
)

func mustMeta(t *testing.T, v any) datatypes.JSON {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal meta: %v", err)
	}
	return datatypes.JSON(b)
}

// TestRuleMatchesTrigger is the routing-gate test (FR-031): a rule runs only via
// the trigger it declares, independent of its Source.
func TestRuleMatchesTrigger(t *testing.T) {
	escrowUpload := dto.DocumentTemplateFileParseRule{Parser: constants.DocumentTemplateFileParseParserEscrowInfo, Source: constants.DocumentTemplateFileParseSourceSingleFile, TriggeredBy: constants.FileParseTriggeredTypeUpload}
	rpaManual := dto.DocumentTemplateFileParseRule{Parser: constants.DocumentTemplateFileParseParserRPA, Source: constants.DocumentTemplateFileParseSourceDirectory, TriggeredBy: constants.FileParseTriggeredTypeManual}
	if !ruleMatchesTrigger(escrowUpload, TriggerOnUpload) {
		t.Error("Upload rule must match on_upload trigger")
	}
	if ruleMatchesTrigger(escrowUpload, TriggerManual) {
		t.Error("Upload rule must NOT match manual trigger")
	}
	if !ruleMatchesTrigger(rpaManual, TriggerManual) {
		t.Error("Manual rule must match manual trigger")
	}
	if ruleMatchesTrigger(rpaManual, TriggerOnUpload) {
		t.Error("Manual rule must NOT match on_upload trigger")
	}
}

func TestEffectiveRules(t *testing.T) {
	tests := []struct {
		name string
		cfg  *dto.DocumentTemplateFileParseConfig
		want int
	}{
		{"nil", nil, 0},
		{"empty", &dto.DocumentTemplateFileParseConfig{}, 0},
		{"rules pass through", &dto.DocumentTemplateFileParseConfig{
			Rules: []dto.DocumentTemplateFileParseRule{{Parser: constants.DocumentTemplateFileParseParserRPA, Source: constants.DocumentTemplateFileParseSourceDirectory}}}, 1},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := effectiveRules(tc.cfg)
			if len(got) != tc.want {
				t.Fatalf("got %d rules, want %d (%+v)", len(got), tc.want, got)
			}
		})
	}
}

// TestOnUploadRules proves trigger and source are decoupled: escrow is
// Upload-triggered despite a Directory source; rpa is Manual-triggered despite a
// SingleFile source. Only the Upload rule fires.
func TestOnUploadRules(t *testing.T) {
	cfg := &dto.DocumentTemplateFileParseConfig{Rules: []dto.DocumentTemplateFileParseRule{
		{Parser: constants.DocumentTemplateFileParseParserEscrowInfo, Source: constants.DocumentTemplateFileParseSourceDirectory, TriggeredBy: constants.FileParseTriggeredTypeUpload},
		{Parser: constants.DocumentTemplateFileParseParserRPA, Source: constants.DocumentTemplateFileParseSourceSingleFile, TriggeredBy: constants.FileParseTriggeredTypeManual},
	}}
	got := onUploadRules(cfg)
	if len(got) != 1 || got[0].Parser != constants.DocumentTemplateFileParseParserEscrowInfo {
		t.Fatalf("onUploadRules = %+v, want only the Upload-triggered escrow rule", got)
	}
}

func TestDocFileParseConfig(t *testing.T) {
	d := &model.Document{Metadata: mustMeta(t, map[string]any{
		"file_parse_config": map[string]any{
			"rules": []map[string]any{{"parser": "Rpa", "triggered_by": "Manual"}},
		},
	})}
	if cfg := docFileParseConfig(d); cfg == nil || len(cfg.Rules) != 1 {
		t.Fatalf("docFileParseConfig = %+v", d)
	}
	if docFileParseConfig(&model.Document{}) != nil {
		t.Fatal("empty doc must yield nil cfg")
	}
}

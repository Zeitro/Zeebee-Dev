package handler

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

func TestEffectiveRules(t *testing.T) {
	tests := []struct {
		name string
		cfg  *dto.DocumentTemplateFileParseConfig
		want int
	}{
		{"nil", nil, 0},
		{"empty", &dto.DocumentTemplateFileParseConfig{}, 0},
		{"rules pass through", &dto.DocumentTemplateFileParseConfig{
			Rules: []dto.DocumentTemplateFileParseRule{{Parser: constants.DocumentTemplateFileParseParserRpa, Source: constants.DocumentTemplateFileParseSourceDirectory}}}, 1},
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

func TestOnUploadRules(t *testing.T) {
	cfg := &dto.DocumentTemplateFileParseConfig{Rules: []dto.DocumentTemplateFileParseRule{
		{Parser: constants.DocumentTemplateFileParseParserEscrowInfo, Source: constants.DocumentTemplateFileParseSourceSingleFile},
		{Parser: constants.DocumentTemplateFileParseParserRpa, Source: constants.DocumentTemplateFileParseSourceDirectory},
	}}
	got := OnUploadRules(cfg)
	if len(got) != 1 || got[0].Parser != constants.DocumentTemplateFileParseParserEscrowInfo {
		t.Fatalf("OnUploadRules = %+v, want only escrow SingleFile", got)
	}
}

func doc(id, parent int64, typ constants.DocumentType, tags ...string) *model.Document {
	meta, _ := json.Marshal(map[string]any{"tags": tags})
	return &model.Document{
		Base:     model.Base{Id: id},
		ParentId: parent,
		Type:     typ,
		Metadata: datatypes.JSON(meta),
	}
}

func TestCollectSubtree(t *testing.T) {
	// tree: 1(Branch) -> [2(Leaf), 3(Branch) -> [4(Leaf)]]
	docs := []*model.Document{
		doc(1, 0, constants.DocumentTypeBranch),
		doc(2, 1, constants.DocumentTypeLeaf),
		doc(3, 1, constants.DocumentTypeBranch),
		doc(4, 3, constants.DocumentTypeLeaf),
	}
	leaves := CollectSubtree(docs, 1)
	if len(leaves) != 2 {
		t.Fatalf("got %d leaves, want 2", len(leaves))
	}
	got := map[int64]bool{}
	for _, d := range leaves {
		got[d.Id] = true
	}
	if !got[2] || !got[4] {
		t.Fatalf("missing leaves; got %v", got)
	}
}

func TestCollectSubtreeAnchorMissing(t *testing.T) {
	docs := []*model.Document{doc(1, 0, constants.DocumentTypeLeaf)}
	if got := CollectSubtree(docs, 999); got != nil {
		t.Fatalf("expected nil for missing anchor, got %+v", got)
	}
}

func TestCollectBySelector(t *testing.T) {
	docs := []*model.Document{
		doc(1, 0, constants.DocumentTypeLeaf, "emd"),
		doc(2, 0, constants.DocumentTypeLeaf, "emd", "signed"),
		doc(3, 0, constants.DocumentTypeLeaf, "signed"),
		doc(4, 0, constants.DocumentTypeBranch, "emd"), // branch excluded
	}
	got := CollectBySelector(docs, &dto.DocumentTemplateTagSelector{Tags: []string{"emd"}})
	if len(got) != 2 {
		t.Fatalf("tag emd: got %d, want 2 (docs 1,2)", len(got))
	}
	got = CollectBySelector(docs, &dto.DocumentTemplateTagSelector{Tags: []string{"emd", "signed"}})
	if len(got) != 1 || got[0].Id != 2 {
		t.Fatalf("AND tags emd+signed: got %+v, want only doc 2", got)
	}
	if got := CollectBySelector(docs, &dto.DocumentTemplateTagSelector{Tags: []string{}}); got != nil {
		t.Fatalf("empty selector must match nothing, got %+v", got)
	}
	if got := CollectBySelector(docs, nil); got != nil {
		t.Fatalf("nil selector must match nothing, got %+v", got)
	}
}

func TestDocFileParseConfigAndTags(t *testing.T) {
	d := &model.Document{Metadata: mustMeta(t, map[string]any{
		"file_parse_config": map[string]any{
			"rules": []map[string]any{{"parser": "Rpa", "mode": "ManualTree"}},
		},
		"tags": []string{"emd"},
	})}
	cfg := DocFileParseConfig(d)
	if cfg == nil || len(cfg.Rules) != 1 {
		t.Fatalf("DocFileParseConfig = %+v", cfg)
	}
	tags := docTags(d)
	if len(tags) != 1 || tags[0] != "emd" {
		t.Fatalf("docTags = %+v", tags)
	}
	// empty metadata is safe
	empty := &model.Document{}
	if DocFileParseConfig(empty) != nil || docTags(empty) != nil {
		t.Fatal("empty doc must yield nil cfg/tags")
	}
}

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

func TestDocTags(t *testing.T) {
	d := &model.Document{Metadata: mustMeta(t, map[string]any{"tags": []string{"emd"}})}
	tags := docTags(d)
	if len(tags) != 1 || tags[0] != "emd" {
		t.Fatalf("docTags = %+v", tags)
	}
	if docTags(&model.Document{}) != nil {
		t.Fatal("empty doc must yield nil tags")
	}
}

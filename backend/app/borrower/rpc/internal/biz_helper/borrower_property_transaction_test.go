package biz_helper

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"

	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
)

func metaJSON(isRequired bool) datatypes.JSON {
	b, _ := json.Marshal(map[string]bool{"is_required": isRequired})
	return b
}

func TestBuildChecklistTree_AllRequiredLeaves(t *testing.T) {
	docs := []*model.Document{
		{Base: model.Base{Id: 1}, ParentId: 0, Name: "Root"},
		{Base: model.Base{Id: 2}, ParentId: 1, Name: "Contract", Metadata: metaJSON(true)},
		{Base: model.Base{Id: 3}, ParentId: 2, Name: "RPA", Metadata: metaJSON(true)},
		{Base: model.Base{Id: 4}, ParentId: 2, Name: "Counter Offer", Metadata: metaJSON(true)},
	}

	result := BuildChecklistTree(docs, nil)

	assert.Len(t, result.RequiredDocuments, 1)
	assert.Len(t, result.OptionalDocuments, 0)

	contract := result.RequiredDocuments[0]
	assert.Equal(t, "Contract", contract.Name)
	assert.Equal(t, int32(2), contract.RequiredDocumentCount)
	assert.Equal(t, int32(0), contract.OptionalDocumentCount)
	assert.Equal(t, int32(2), contract.SubDocumentCount)

	assert.Equal(t, int32(2), result.TotalDocumentCount)
	assert.Equal(t, int32(2), result.TotalRequiredDocumentCount)
	assert.Equal(t, int32(0), result.TotalOptionalDocumentCount)
}

func TestBuildChecklistTree_MixedRequiredOptional(t *testing.T) {
	docs := []*model.Document{
		{Base: model.Base{Id: 1}, ParentId: 0, Name: "Root"},
		{Base: model.Base{Id: 2}, ParentId: 1, Name: "Insurance", Metadata: metaJSON(true)},
		{Base: model.Base{Id: 3}, ParentId: 2, Name: "EOI", Metadata: metaJSON(true)},
		{Base: model.Base{Id: 4}, ParentId: 2, Name: "Full Policy", Metadata: metaJSON(true)},
		{Base: model.Base{Id: 5}, ParentId: 2, Name: "Flood Insurance", Metadata: metaJSON(false)},
	}

	result := BuildChecklistTree(docs, nil)

	assert.Len(t, result.RequiredDocuments, 1)
	insurance := result.RequiredDocuments[0]
	assert.Equal(t, int32(2), insurance.RequiredDocumentCount)
	assert.Equal(t, int32(1), insurance.OptionalDocumentCount)
	assert.Equal(t, int32(3), insurance.SubDocumentCount)

	assert.Equal(t, int32(2), result.TotalRequiredDocumentCount)
	assert.Equal(t, int32(1), result.TotalOptionalDocumentCount)
	assert.Equal(t, int32(3), result.TotalDocumentCount)
}

func TestBuildChecklistTree_TopLevelOptional(t *testing.T) {
	docs := []*model.Document{
		{Base: model.Base{Id: 1}, ParentId: 0, Name: "Root"},
		{Base: model.Base{Id: 2}, ParentId: 1, Name: "HOA", Metadata: metaJSON(false)},
		{Base: model.Base{Id: 3}, ParentId: 2, Name: "HOA Sheet", Metadata: metaJSON(false)},
		{Base: model.Base{Id: 4}, ParentId: 2, Name: "HOA Rules", Metadata: metaJSON(false)},
	}

	result := BuildChecklistTree(docs, nil)

	assert.Len(t, result.RequiredDocuments, 0)
	assert.Len(t, result.OptionalDocuments, 1)

	hoa := result.OptionalDocuments[0]
	assert.Equal(t, int32(2), hoa.OptionalDocumentCount)
	assert.Equal(t, int32(0), hoa.RequiredDocumentCount)

	assert.Equal(t, int32(0), result.TotalRequiredDocumentCount)
	assert.Equal(t, int32(2), result.TotalOptionalDocumentCount)
}

func TestBuildChecklistTree_SingleRequiredLeaf(t *testing.T) {
	docs := []*model.Document{
		{Base: model.Base{Id: 1}, ParentId: 0, Name: "Root"},
		{Base: model.Base{Id: 2}, ParentId: 1, Name: "Disclosures", Metadata: metaJSON(true)},
	}

	result := BuildChecklistTree(docs, nil)

	assert.Len(t, result.RequiredDocuments, 1)
	assert.Equal(t, int32(0), result.TotalRequiredDocumentCount,
		"leaf node itself should not count as required child")
	assert.Equal(t, int32(0), result.TotalOptionalDocumentCount)
	assert.Equal(t, int32(0), result.TotalDocumentCount)

	leaf := result.RequiredDocuments[0]
	assert.Equal(t, int32(0), leaf.RequiredDocumentCount)
	assert.Equal(t, int32(0), leaf.OptionalDocumentCount)
	assert.Equal(t, int32(0), leaf.SubDocumentCount)
}

func TestBuildChecklistTree_FullScenario(t *testing.T) {
	// Simulates the full Purchase transaction from user's expected output
	docs := []*model.Document{
		{Base: model.Base{Id: 1}, ParentId: 0, Name: "Root"},

		// Contract (required, 6 required leaves)
		{Base: model.Base{Id: 10}, ParentId: 1, Name: "Contract", Metadata: metaJSON(true)},
		{Base: model.Base{Id: 11}, ParentId: 10, Name: "Original RPA", Metadata: metaJSON(true)},
		{Base: model.Base{Id: 12}, ParentId: 10, Name: "Assignment", Metadata: metaJSON(true)},
		{Base: model.Base{Id: 13}, ParentId: 10, Name: "Counter Offer", Metadata: metaJSON(true)},
		{Base: model.Base{Id: 14}, ParentId: 10, Name: "Addendum", Metadata: metaJSON(true)},
		{Base: model.Base{Id: 15}, ParentId: 10, Name: "Seller financing", Metadata: metaJSON(true)},
		{Base: model.Base{Id: 16}, ParentId: 10, Name: "JV Agreement", Metadata: metaJSON(true)},

		// Escrow (required, 9 required leaves)
		{Base: model.Base{Id: 20}, ParentId: 1, Name: "Escrow", Metadata: metaJSON(true)},
		{Base: model.Base{Id: 21}, ParentId: 20, Name: "Wire Instructions", Metadata: metaJSON(true)},
		{Base: model.Base{Id: 22}, ParentId: 20, Name: "EMD Payment Proof", Metadata: metaJSON(true)},
		{Base: model.Base{Id: 23}, ParentId: 20, Name: "EMD Receipt", Metadata: metaJSON(true)},
		{Base: model.Base{Id: 24}, ParentId: 20, Name: "Prelim Report", Metadata: metaJSON(true)},
		{Base: model.Base{Id: 25}, ParentId: 20, Name: "NHD Report", Metadata: metaJSON(true)},
		{Base: model.Base{Id: 26}, ParentId: 20, Name: "Colored map", Metadata: metaJSON(true)},
		{Base: model.Base{Id: 27}, ParentId: 20, Name: "FinCen", Metadata: metaJSON(true)},
		{Base: model.Base{Id: 28}, ParentId: 20, Name: "Buyer Open Package", Metadata: metaJSON(true)},
		{Base: model.Base{Id: 29}, ParentId: 20, Name: "LOI", Metadata: metaJSON(true)},

		// Disclosures (required, leaf)
		{Base: model.Base{Id: 30}, ParentId: 1, Name: "Disclosures", Metadata: metaJSON(true)},

		// Loan (required, 4 required leaves)
		{Base: model.Base{Id: 40}, ParentId: 1, Name: "Loan", Metadata: metaJSON(true)},
		{Base: model.Base{Id: 41}, ParentId: 40, Name: "Application Form", Metadata: metaJSON(true)},
		{Base: model.Base{Id: 42}, ParentId: 40, Name: "Loan Document", Metadata: metaJSON(true)},
		{Base: model.Base{Id: 43}, ParentId: 40, Name: "Contact form", Metadata: metaJSON(true)},
		{Base: model.Base{Id: 44}, ParentId: 40, Name: "Appraisal Report", Metadata: metaJSON(true)},

		// Insurance (required, 3 required + 1 optional leaves)
		{Base: model.Base{Id: 50}, ParentId: 1, Name: "Insurance", Metadata: metaJSON(true)},
		{Base: model.Base{Id: 51}, ParentId: 50, Name: "EOI", Metadata: metaJSON(true)},
		{Base: model.Base{Id: 52}, ParentId: 50, Name: "Payment Receipt", Metadata: metaJSON(true)},
		{Base: model.Base{Id: 53}, ParentId: 50, Name: "Full Policy", Metadata: metaJSON(true)},
		{Base: model.Base{Id: 54}, ParentId: 50, Name: "Flood Insurance", Metadata: metaJSON(false)},

		// Closing Package (required, 5 required leaves)
		{Base: model.Base{Id: 60}, ParentId: 1, Name: "Closing Package", Metadata: metaJSON(true)},
		{Base: model.Base{Id: 61}, ParentId: 60, Name: "Buyer Estimated Statement", Metadata: metaJSON(true)},
		{Base: model.Base{Id: 62}, ParentId: 60, Name: "Final Buyer Statement", Metadata: metaJSON(true)},
		{Base: model.Base{Id: 63}, ParentId: 60, Name: "All Signed Documents", Metadata: metaJSON(true)},
		{Base: model.Base{Id: 64}, ParentId: 60, Name: "Recorded GD", Metadata: metaJSON(true)},
		{Base: model.Base{Id: 65}, ParentId: 60, Name: "Recorded DOT", Metadata: metaJSON(true)},

		// 1031 (optional, 1 optional leaf)
		{Base: model.Base{Id: 70}, ParentId: 1, Name: "1031", Metadata: metaJSON(false)},
		{Base: model.Base{Id: 71}, ParentId: 70, Name: "1031 Exchange Documents", Metadata: metaJSON(false)},

		// Lease (optional, 2 optional leaves)
		{Base: model.Base{Id: 80}, ParentId: 1, Name: "Lease", Metadata: metaJSON(false)},
		{Base: model.Base{Id: 81}, ParentId: 80, Name: "Lease Agreement", Metadata: metaJSON(false)},
		{Base: model.Base{Id: 82}, ParentId: 80, Name: "Tenant Estoppel", Metadata: metaJSON(false)},

		// HOA (optional, 1 optional leaf)
		{Base: model.Base{Id: 90}, ParentId: 1, Name: "HOA", Metadata: metaJSON(false)},
		{Base: model.Base{Id: 91}, ParentId: 90, Name: "HOA Sheet", Metadata: metaJSON(false)},

		// Solar (optional, 1 optional leaf)
		{Base: model.Base{Id: 100}, ParentId: 1, Name: "Solar", Metadata: metaJSON(false)},
		{Base: model.Base{Id: 101}, ParentId: 100, Name: "Solar Sheet", Metadata: metaJSON(false)},
	}

	result := BuildChecklistTree(docs, nil)

	assert.Len(t, result.RequiredDocuments, 6)
	assert.Len(t, result.OptionalDocuments, 4)

	assert.Equal(t, int32(27), result.TotalRequiredDocumentCount)
	assert.Equal(t, int32(6), result.TotalOptionalDocumentCount)
	assert.Equal(t, int32(33), result.TotalDocumentCount)

	// Verify Insurance node
	var insuranceNode *dto.ChecklistNode
	for _, n := range result.RequiredDocuments {
		if n.Name == "Insurance" {
			insuranceNode = n
			break
		}
	}
	assert.Equal(t, int32(3), insuranceNode.RequiredDocumentCount)
	assert.Equal(t, int32(1), insuranceNode.OptionalDocumentCount)

	// Verify leaf nodes have 0 counts
	var disclosuresNode *dto.ChecklistNode
	for _, n := range result.RequiredDocuments {
		if n.Name == "Disclosures" {
			disclosuresNode = n
			break
		}
	}
	assert.Equal(t, int32(0), disclosuresNode.RequiredDocumentCount)
	assert.Equal(t, int32(0), disclosuresNode.OptionalDocumentCount)
	assert.Equal(t, int32(0), disclosuresNode.SubDocumentCount)
}

func TestBuildChecklistTree_WithFileCounts(t *testing.T) {
	docs := []*model.Document{
		{Base: model.Base{Id: 1}, ParentId: 0, Name: "Root"},
		{Base: model.Base{Id: 2}, ParentId: 1, Name: "Contract", Metadata: metaJSON(true)},
		{Base: model.Base{Id: 3}, ParentId: 2, Name: "RPA", Metadata: metaJSON(true)},
		{Base: model.Base{Id: 4}, ParentId: 2, Name: "Counter Offer", Metadata: metaJSON(true)},
		{Base: model.Base{Id: 5}, ParentId: 2, Name: "Counter Offer2", Metadata: metaJSON(false)},
	}

	fileCounts := map[int64]int32{
		3: 2,
		5: 1,
	} // RPA has 2 files

	result := BuildChecklistTree(docs, fileCounts)

	assert.Equal(t, int32(1), result.TotalRequiredUploadedDocumentCount,
		"1 leaf (RPA) has files uploaded")

	contract := result.RequiredDocuments[0]
	assert.Equal(t, int32(1), contract.UploadedDocumentCount,
		"1 direct child has uploads")
}

func TestBuildChecklistTree_EmptyDocuments(t *testing.T) {
	result := BuildChecklistTree([]*model.Document{}, nil)

	assert.Equal(t, int32(0), result.TotalDocumentCount)
	assert.Equal(t, int32(0), result.TotalRequiredDocumentCount)
	assert.Equal(t, int32(0), result.TotalOptionalDocumentCount)
	assert.Nil(t, result.RequiredDocuments)
	assert.Nil(t, result.OptionalDocuments)
}

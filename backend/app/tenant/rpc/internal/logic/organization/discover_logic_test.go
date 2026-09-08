package organizationlogic

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/domains/zeebee/lenderdiscover"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/proto"
)

func TestDiscoverRunResultToProto(t *testing.T) {
	got, err := discoverRunResultToProto(&lenderdiscover.RunResult{
		TriggerType: "scheduled",
		Status:      "success",
		ErrorMsg:    "boom",
		Stats: lenderdiscover.RunStats{
			RequestedCount: 5, InsertedCount: 3, SkippedDuplicateCount: 1, ErrorCount: 1,
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, "scheduled", got.TriggerType)
	assert.Equal(t, "boom", got.ErrorMessage)
	assert.Equal(t, int32(5), got.Stats.RequestedCount)
}

func TestDiscoverRawLenderFromCommitProto(t *testing.T) {
	raw, err := discoverRawLenderFromCommitProto(&tenant.LenderDiscoverCommitLender{
		Name: "Acme Lending",
		LenderMetadata: &proto.OrganizationLenderMetadataEntity{
			LoanAmountRange: &proto.RangeFloat{Min: 1, Max: 2},
			LoanTermRange:   &proto.RangeInt32{Min: 7, Max: 9},
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, "Acme Lending", raw.Name)
	if assert.NotNil(t, raw.Metadata.Lender) {
		assert.Equal(t, &model.RangeFloat64{Min: 1, Max: 2}, raw.Metadata.Lender.LoanAmountRange)
		assert.Equal(t, &model.RangeInt{Min: 7, Max: 9}, raw.Metadata.Lender.LoanTermRange)
	}
}

func TestDiscoverLenderMetadataToProto_Empty(t *testing.T) {
	got, err := discoverLenderMetadataToProto(lenderdiscover.RawLender{})
	assert.NoError(t, err)
	assert.Nil(t, got)
}

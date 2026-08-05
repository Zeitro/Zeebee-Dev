package organizationlogic

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/domains/zeebee/lenderdiscover"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/proto"
)

func TestDiscoverRunResultToProto(t *testing.T) {
	// nil -> empty entity
	got, err := discoverRunResultToProto(nil)
	assert.NoError(t, err)
	assert.Equal(t, &admin.LenderDiscoverRunEntity{}, got)

	// ErrorMsg maps to ErrorMessage; stats int -> int32 via copierx.
	in := &lenderdiscover.RunResult{
		TriggerType: "scheduled",
		Status:      "success",
		ErrorMsg:    "boom",
		Stats: lenderdiscover.RunStats{
			RequestedCount: 5, InsertedCount: 3, SkippedDuplicateCount: 1, ErrorCount: 1,
		},
	}
	got, err = discoverRunResultToProto(in)
	assert.NoError(t, err)
	assert.Equal(t, "scheduled", got.TriggerType)
	assert.Equal(t, "success", got.Status)
	assert.Equal(t, "boom", got.ErrorMessage)
	if assert.NotNil(t, got.Stats) {
		assert.Equal(t, int32(5), got.Stats.RequestedCount)
		assert.Equal(t, int32(3), got.Stats.InsertedCount)
		assert.Equal(t, int32(1), got.Stats.SkippedDuplicateCount)
		assert.Equal(t, int32(1), got.Stats.ErrorCount)
	}
}

func rawLenderWith(lender *lenderdiscover.RawLenderMetadata) lenderdiscover.RawLender {
	return lenderdiscover.RawLender{Metadata: struct {
		Lender *lenderdiscover.RawLenderMetadata `json:"lender"`
	}{Lender: lender}}
}

func TestDiscoverLenderMetadataToProto(t *testing.T) {
	// Lender nil -> nil
	got, err := discoverLenderMetadataToProto(lenderdiscover.RawLender{})
	assert.NoError(t, err)
	assert.Nil(t, got)

	// all ranges nil -> nil (prior behavior)
	got, err = discoverLenderMetadataToProto(rawLenderWith(&lenderdiscover.RawLenderMetadata{}))
	assert.NoError(t, err)
	assert.Nil(t, got)

	// populated: float passthrough + int->int32, unset ranges stay nil
	got, err = discoverLenderMetadataToProto(rawLenderWith(&lenderdiscover.RawLenderMetadata{
		LoanAmountRange: &model.RangeFloat64{Min: 10000.5, Max: 50000},
		LoanTermRange:   &model.RangeInt{Min: 12, Max: 360},
	}))
	assert.NoError(t, err)
	if assert.NotNil(t, got) {
		assert.Equal(t, &proto.RangeFloat{Min: 10000.5, Max: 50000}, got.LoanAmountRange)
		assert.Equal(t, &proto.RangeInt32{Min: 12, Max: 360}, got.LoanTermRange)
		assert.Nil(t, got.LoanRateRange)
		assert.Nil(t, got.LoanPointRange)
	}
}

func TestDiscoverRawLenderFromCommitProto_MetadataRoundTrip(t *testing.T) {
	// proto flat ranges -> RawLender nested metadata; int32 -> int.
	commit := &admin.LenderDiscoverCommitLender{
		Name: "Acme Lending",
		LenderMetadata: &proto.OrganizationLenderMetadataEntity{
			LoanAmountRange: &proto.RangeFloat{Min: 1, Max: 2},
			LoanTermRange:   &proto.RangeInt32{Min: 7, Max: 9},
		},
	}
	out, err := discoverRawLenderFromCommitProto(commit)
	assert.NoError(t, err)
	assert.Equal(t, "Acme Lending", out.Name)
	if assert.NotNil(t, out.Metadata.Lender) {
		assert.Equal(t, &model.RangeFloat64{Min: 1, Max: 2}, out.Metadata.Lender.LoanAmountRange)
		assert.Equal(t, &model.RangeInt{Min: 7, Max: 9}, out.Metadata.Lender.LoanTermRange)
		assert.Nil(t, out.Metadata.Lender.LoanRateRange)
	}

	// nil metadata -> nil Lender
	commit2 := &admin.LenderDiscoverCommitLender{Name: "NoMeta"}
	out2, err := discoverRawLenderFromCommitProto(commit2)
	assert.NoError(t, err)
	assert.Nil(t, out2.Metadata.Lender)
}

package organizationlogic

import (
	"strings"

	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/lenderdiscover"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/proto"
)

func discoverConfigToProto(cfg *lenderdiscover.Config) (*admin.LenderDiscoverSettingsEntity, error) {
	return copierx.Trans[lenderdiscover.Config, admin.LenderDiscoverSettingsEntity](cfg)
}

func discoverRunResultToProto(r *lenderdiscover.RunResult) *admin.LenderDiscoverRunEntity {
	if r == nil {
		return &admin.LenderDiscoverRunEntity{}
	}
	return &admin.LenderDiscoverRunEntity{
		TriggerType:  r.TriggerType,
		Status:       r.Status,
		ErrorMessage: r.ErrorMsg,
		Stats: &admin.LenderDiscoverRunStats{
			RequestedCount:        int32(r.Stats.RequestedCount),
			InsertedCount:         int32(r.Stats.InsertedCount),
			SkippedDuplicateCount: int32(r.Stats.SkippedDuplicateCount),
			ErrorCount:            int32(r.Stats.ErrorCount),
		},
	}
}

func discoverCandidateToProto(c *lenderdiscover.PreviewCandidate) (*admin.LenderDiscoverCandidateEntity, error) {
	if c == nil {
		return nil, nil
	}
	ent, err := copierx.Trans[lenderdiscover.RawLender, admin.LenderDiscoverCandidateEntity](&c.Raw)
	if err != nil {
		return nil, err
	}
	ent.Name = strings.TrimSpace(ent.Name)
	ent.DedupeStatus = c.DedupeStatus
	ent.NormalizedKey = c.NormalizedKey
	ent.LenderMetadata = discoverLenderMetadataToProto(c.Raw)
	return ent, nil
}

// lenderMetadataToProto handles manual conversion due to nested pointer types
// and int->int32 conversion that copier cannot reliably handle.
func discoverLenderMetadataToProto(raw lenderdiscover.RawLender) *proto.OrganizationLenderMetadataEntity {
	if raw.Metadata.Lender == nil {
		return nil
	}
	ml := raw.Metadata.Lender
	out := &proto.OrganizationLenderMetadataEntity{}
	if ml.LoanAmountRange != nil {
		out.LoanAmountRange = &proto.RangeFloat{Min: ml.LoanAmountRange.Min, Max: ml.LoanAmountRange.Max}
	}
	if ml.LoanTermRange != nil {
		out.LoanTermRange = &proto.RangeInt32{Min: int32(ml.LoanTermRange.Min), Max: int32(ml.LoanTermRange.Max)}
	}
	if ml.LoanRateRange != nil {
		out.LoanRateRange = &proto.RangeFloat{Min: ml.LoanRateRange.Min, Max: ml.LoanRateRange.Max}
	}
	if ml.LoanPointRange != nil {
		out.LoanPointRange = &proto.RangeFloat{Min: ml.LoanPointRange.Min, Max: ml.LoanPointRange.Max}
	}
	if out.LoanAmountRange == nil && out.LoanTermRange == nil && out.LoanRateRange == nil && out.LoanPointRange == nil {
		return nil
	}
	return out
}

func discoverRawLenderFromCommitProto(p *admin.LenderDiscoverCommitLender) (lenderdiscover.RawLender, error) {
	if p == nil {
		return lenderdiscover.RawLender{}, nil
	}
	out, err := copierx.Trans[admin.LenderDiscoverCommitLender, lenderdiscover.RawLender](p)
	if err != nil {
		return lenderdiscover.RawLender{}, err
	}
	if lm := p.GetLenderMetadata(); lm != nil {
		inner := &struct {
			LoanAmountRange *model.RangeFloat64 `json:"loan_amount_range"`
			LoanTermRange   *model.RangeInt     `json:"loan_term_range"`
			LoanRateRange   *model.RangeFloat64 `json:"loan_rate_range"`
			LoanPointRange  *model.RangeFloat64 `json:"loan_point_range"`
		}{}
		if r := lm.GetLoanAmountRange(); r != nil {
			inner.LoanAmountRange = &model.RangeFloat64{Min: r.GetMin(), Max: r.GetMax()}
		}
		if r := lm.GetLoanTermRange(); r != nil {
			inner.LoanTermRange = &model.RangeInt{Min: int(r.GetMin()), Max: int(r.GetMax())}
		}
		if r := lm.GetLoanRateRange(); r != nil {
			inner.LoanRateRange = &model.RangeFloat64{Min: r.GetMin(), Max: r.GetMax()}
		}
		if r := lm.GetLoanPointRange(); r != nil {
			inner.LoanPointRange = &model.RangeFloat64{Min: r.GetMin(), Max: r.GetMax()}
		}
		out.Metadata.Lender = inner
	}
	return *out, nil
}

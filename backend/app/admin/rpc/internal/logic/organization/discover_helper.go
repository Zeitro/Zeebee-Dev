package organizationlogic

import (
	"strings"

	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/lenderdiscover"
	"zeebee-backend/domains/zeebee/proto"
)

func discoverConfigToProto(cfg *lenderdiscover.Config) (*admin.LenderDiscoverSettingsEntity, error) {
	return copierx.Trans[lenderdiscover.Config, admin.LenderDiscoverSettingsEntity](cfg)
}

func discoverRunResultToProto(r *lenderdiscover.RunResult) (*admin.LenderDiscoverRunEntity, error) {
	if r == nil {
		return &admin.LenderDiscoverRunEntity{}, nil
	}
	// RunResult.ErrorMsg and the entity.ErrorMessage differ by name, so the message
	// is mapped by hand; the stats (same names, int->int32) go through copierx.
	stats, err := copierx.Trans[lenderdiscover.RunStats, admin.LenderDiscoverRunStats](&r.Stats)
	if err != nil {
		return nil, err
	}
	return &admin.LenderDiscoverRunEntity{
		TriggerType:  r.TriggerType,
		Status:       r.Status,
		ErrorMessage: r.ErrorMsg,
		Stats:        stats,
	}, nil
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
	if ent.LenderMetadata, err = discoverLenderMetadataToProto(c.Raw); err != nil {
		return nil, err
	}
	return ent, nil
}

// discoverLenderMetadataToProto maps the nested metadata.lender ranges to the
// flat proto entity. copierx handles *model.Range -> *proto.Range and
// int -> int32 in one pass.
func discoverLenderMetadataToProto(raw lenderdiscover.RawLender) (*proto.OrganizationLenderMetadataEntity, error) {
	if raw.Metadata.Lender == nil {
		return nil, nil
	}
	ent, err := copierx.Trans[lenderdiscover.RawLenderMetadata, proto.OrganizationLenderMetadataEntity](raw.Metadata.Lender)
	if err != nil {
		return nil, err
	}
	if ent.LoanAmountRange == nil && ent.LoanTermRange == nil && ent.LoanRateRange == nil && ent.LoanPointRange == nil {
		return nil, nil
	}
	return ent, nil
}

func discoverRawLenderFromCommitProto(p *admin.LenderDiscoverCommitLender) (lenderdiscover.RawLender, error) {
	if p == nil {
		return lenderdiscover.RawLender{}, nil
	}
	out, err := copierx.Trans[admin.LenderDiscoverCommitLender, lenderdiscover.RawLender](p)
	if err != nil {
		return lenderdiscover.RawLender{}, err
	}
	// LenderDiscoverCommitLender.LenderMetadata (flat) <-> RawLender.Metadata.Lender
	// (nested) differ in shape, so the metadata is mapped separately.
	if lm := p.GetLenderMetadata(); lm != nil {
		if out.Metadata.Lender, err = copierx.Trans[proto.OrganizationLenderMetadataEntity, lenderdiscover.RawLenderMetadata](lm); err != nil {
			return lenderdiscover.RawLender{}, err
		}
	}
	return *out, nil
}

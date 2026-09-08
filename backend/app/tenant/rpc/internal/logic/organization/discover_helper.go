package organizationlogic

import (
	"strings"

	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/lenderdiscover"
	"zeebee-backend/domains/zeebee/proto"
)

func discoverConfigToProto(cfg *lenderdiscover.Config) (*tenant.LenderDiscoverSettingsEntity, error) {
	return copierx.Trans[lenderdiscover.Config, tenant.LenderDiscoverSettingsEntity](cfg)
}

func discoverRunResultToProto(r *lenderdiscover.RunResult) (*tenant.LenderDiscoverRunEntity, error) {
	if r == nil {
		return &tenant.LenderDiscoverRunEntity{}, nil
	}
	stats, err := copierx.Trans[lenderdiscover.RunStats, tenant.LenderDiscoverRunStats](&r.Stats)
	if err != nil {
		return nil, err
	}
	return &tenant.LenderDiscoverRunEntity{TriggerType: r.TriggerType, Status: r.Status, ErrorMessage: r.ErrorMsg, Stats: stats}, nil
}

func discoverCandidateToProto(c *lenderdiscover.PreviewCandidate) (*tenant.LenderDiscoverCandidateEntity, error) {
	if c == nil {
		return nil, nil
	}
	ent, err := copierx.Trans[lenderdiscover.RawLender, tenant.LenderDiscoverCandidateEntity](&c.Raw)
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

func discoverRawLenderFromCommitProto(p *tenant.LenderDiscoverCommitLender) (lenderdiscover.RawLender, error) {
	if p == nil {
		return lenderdiscover.RawLender{}, nil
	}
	out, err := copierx.Trans[tenant.LenderDiscoverCommitLender, lenderdiscover.RawLender](p)
	if err != nil {
		return lenderdiscover.RawLender{}, err
	}
	if lm := p.GetLenderMetadata(); lm != nil {
		if out.Metadata.Lender, err = copierx.Trans[proto.OrganizationLenderMetadataEntity, lenderdiscover.RawLenderMetadata](lm); err != nil {
			return lenderdiscover.RawLender{}, err
		}
	}
	return *out, nil
}

package lender_discover

import (
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/app/admin/rpc/pb/admin"
)

func protoSettingsToTypes(p *admin.LenderDiscoverSettingsEntity) types.LenderDiscoverSettings {
	if p == nil {
		return types.LenderDiscoverSettings{}
	}
	return types.LenderDiscoverSettings{
		ScheduleEnabled: p.ScheduleEnabled,
		IntervalDays:    int(p.IntervalDays),
		ScopeState:      p.ScopeState,
		ScopeExtra:      p.ScopeExtra,
		UserPrompt:      p.UserPrompt,
		OpenaiModel:     p.OpenaiModel,
	}
}

func protoRunToTypes(p *admin.LenderDiscoverRunEntity) types.LenderDiscoverRunEntity {
	if p == nil {
		return types.LenderDiscoverRunEntity{Stats: types.LenderDiscoverRunStats{}}
	}
	st := types.LenderDiscoverRunStats{}
	if p.Stats != nil {
		st.RequestedCount = int(p.Stats.RequestedCount)
		st.InsertedCount = int(p.Stats.InsertedCount)
		st.SkippedDuplicateCount = int(p.Stats.SkippedDuplicateCount)
		st.ErrorCount = int(p.Stats.ErrorCount)
	}
	return types.LenderDiscoverRunEntity{
		Id:           p.Id,
		TriggerType:  p.TriggerType,
		Status:       p.Status,
		StartedAt:    p.StartedAt,
		FinishedAt:   p.FinishedAt,
		ErrorMessage: p.ErrorMessage,
		Stats:        st,
	}
}

func protoPreviewDataToTypes(p *admin.LenderDiscoverPreviewResponseData) types.LenderDiscoverPreviewData {
	if p == nil {
		return types.LenderDiscoverPreviewData{}
	}
	out := types.LenderDiscoverPreviewData{
		NewCount:     int(p.NewCount),
		ExistedCount: int(p.ExistedCount),
	}
	for _, c := range p.Candidates {
		out.Candidates = append(out.Candidates, protoCandidateToTypes(c))
	}
	return out
}

func protoCandidateToTypes(c *admin.LenderDiscoverCandidateEntity) types.LenderDiscoverCandidate {
	if c == nil {
		return types.LenderDiscoverCandidate{}
	}
	t := types.LenderDiscoverCandidate{
		DedupeStatus:  c.GetDedupeStatus(),
		NormalizedKey: c.GetNormalizedKey(),
		Name:          c.GetName(),
		ContactEmail:  c.GetContactEmail(),
		ContactPhone:  c.GetContactPhone(),
		ContactName:   c.GetContactName(),
		Website:       c.GetWebsite(),
		Tagline:       c.GetTagline(),
		Logo:          c.GetLogo(),
	}
	if a := c.GetAddress(); a != nil {
		t.Address = types.AddressEntity{
			State:    a.GetState(),
			City:     a.GetCity(),
			Address1: a.GetAddress1(),
			Address2: a.GetAddress2(),
			ZipCode:  a.GetZipCode(),
		}
	}
	if lm := c.GetLenderMetadata(); lm != nil {
		t.LenderMetadata = types.LenderDiscoverLenderMetadata{}
		if r := lm.GetLoanAmountRange(); r != nil {
			t.LenderMetadata.LoanAmountRange = types.RangeFloat64{Min: r.GetMin(), Max: r.GetMax()}
		}
		if r := lm.GetLoanTermRange(); r != nil {
			t.LenderMetadata.LoanTermRange = types.RangeInt{Min: int(r.GetMin()), Max: int(r.GetMax())}
		}
		if r := lm.GetLoanRateRange(); r != nil {
			t.LenderMetadata.LoanRateRange = types.RangeFloat64{Min: r.GetMin(), Max: r.GetMax()}
		}
		if r := lm.GetLoanPointRange(); r != nil {
			t.LenderMetadata.LoanPointRange = types.RangeFloat64{Min: r.GetMin(), Max: r.GetMax()}
		}
	}
	return t
}

func typesCommitRequestToProto(req *types.LenderDiscoverCommitRequest) *admin.LenderDiscoverCommitRequest {
	if req == nil {
		return &admin.LenderDiscoverCommitRequest{Data: &admin.LenderDiscoverCommitRequestData{}}
	}
	list := make([]*admin.LenderDiscoverCommitLender, 0, len(req.Lenders))
	for i := range req.Lenders {
		list = append(list, typesCommitLenderToProto(&req.Lenders[i]))
	}
	return &admin.LenderDiscoverCommitRequest{
		Data: &admin.LenderDiscoverCommitRequestData{Lenders: list},
	}
}

func typesCommitLenderToProto(t *types.LenderDiscoverCommitLender) *admin.LenderDiscoverCommitLender {
	if t == nil {
		return &admin.LenderDiscoverCommitLender{}
	}
	m := t.LenderMetadata
	return &admin.LenderDiscoverCommitLender{
		Name:         t.Name,
		ContactEmail: t.ContactEmail,
		ContactPhone: t.ContactPhone,
		ContactName:  t.ContactName,
		Website:      t.Website,
		Tagline:      t.Tagline,
		Logo:         t.Logo,
		Address: &admin.AddressEntity{
			State:    t.Address.State,
			City:     t.Address.City,
			Address1: t.Address.Address1,
			Address2: t.Address.Address2,
			ZipCode:  t.Address.ZipCode,
		},
		LenderMetadata: &admin.OrganizationLenderMetadataEntity{
			LoanAmountRange: &admin.RangeFloat{Min: m.LoanAmountRange.Min, Max: m.LoanAmountRange.Max},
			LoanTermRange:   &admin.RangeInt32{Min: int32(m.LoanTermRange.Min), Max: int32(m.LoanTermRange.Max)},
			LoanRateRange:   &admin.RangeFloat{Min: m.LoanRateRange.Min, Max: m.LoanRateRange.Max},
			LoanPointRange:  &admin.RangeFloat{Min: m.LoanPointRange.Min, Max: m.LoanPointRange.Max},
		},
	}
}

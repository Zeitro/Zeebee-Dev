package lenderdiscover

import (
	"encoding/json"
	"strings"
	"time"

	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/domains/zeebee/model"
	domainlenderdiscover "zeebee-backend/domains/zeebee/lenderdiscover"
)

func settingsToProto(m *model.LenderDiscoverSettings) *admin.LenderDiscoverSettingsEntity {
	if m == nil {
		return nil
	}
	return &admin.LenderDiscoverSettingsEntity{
		ScheduleEnabled: m.ScheduleEnabled,
		IntervalDays:    m.IntervalDays,
		ScopeState:      m.ScopeState,
		ScopeExtra:      m.ScopeExtra,
		UserPrompt:      m.UserPrompt,
		OpenaiModel:     m.OpenaiModel,
	}
}

func runToProto(r *model.LenderDiscoverRun) *admin.LenderDiscoverRunEntity {
	if r == nil {
		return &admin.LenderDiscoverRunEntity{}
	}
	ent := &admin.LenderDiscoverRunEntity{
		Id:           r.Id,
		TriggerType:  r.TriggerType,
		Status:       r.Status,
		StartedAt:    r.StartedAt.UTC().Format(time.RFC3339),
		ErrorMessage: r.ErrorMessage,
		Stats:        statsFromMap(r.Stats),
	}
	if r.FinishedAt != nil {
		ent.FinishedAt = r.FinishedAt.UTC().Format(time.RFC3339)
	}
	return ent
}

func statsFromMap(m map[string]interface{}) *admin.LenderDiscoverRunStats {
	if len(m) == 0 {
		return &admin.LenderDiscoverRunStats{}
	}
	s := &admin.LenderDiscoverRunStats{
		RequestedCount:        int32(toInt64(m["requested_count"])),
		InsertedCount:         int32(toInt64(m["inserted_count"])),
		SkippedDuplicateCount: int32(toInt64(m["skipped_duplicate_count"])),
		ErrorCount:            int32(toInt64(m["error_count"])),
	}
	return s
}

func toInt64(v interface{}) int64 {
	switch x := v.(type) {
	case int:
		return int64(x)
	case int32:
		return int64(x)
	case int64:
		return x
	case float64:
		return int64(x)
	case json.Number:
		i, _ := x.Int64()
		return i
	default:
		return 0
	}
}

func candidateToProto(c *domainlenderdiscover.PreviewCandidate) *admin.LenderDiscoverCandidateEntity {
	if c == nil {
		return nil
	}
	raw := c.Raw
	ent := &admin.LenderDiscoverCandidateEntity{
		DedupeStatus:  c.DedupeStatus,
		NormalizedKey: c.NormalizedKey,
		Name:          strings.TrimSpace(raw.Name),
		ContactEmail:  raw.ContactEmail,
		ContactPhone:  raw.ContactPhone,
		ContactName:   raw.ContactName,
		Website:       raw.Website,
		Tagline:       raw.Tagline,
		Logo:          raw.Logo,
		Address: &admin.AddressEntity{
			State:    raw.Address.State,
			City:     raw.Address.City,
			Address1: raw.Address.Address1,
			Address2: raw.Address.Address2,
			ZipCode:  raw.Address.ZipCode,
		},
		LenderMetadata: lenderMetadataToProto(raw),
	}
	return ent
}

func lenderMetadataToProto(raw domainlenderdiscover.RawLender) *admin.OrganizationLenderMetadataEntity {
	if raw.Metadata.Lender == nil {
		return nil
	}
	ml := raw.Metadata.Lender
	out := &admin.OrganizationLenderMetadataEntity{}
	if ml.LoanAmountRange != nil {
		out.LoanAmountRange = &admin.RangeFloat{Min: ml.LoanAmountRange.Min, Max: ml.LoanAmountRange.Max}
	}
	if ml.LoanTermRange != nil {
		out.LoanTermRange = &admin.RangeInt32{Min: int32(ml.LoanTermRange.Min), Max: int32(ml.LoanTermRange.Max)}
	}
	if ml.LoanRateRange != nil {
		out.LoanRateRange = &admin.RangeFloat{Min: ml.LoanRateRange.Min, Max: ml.LoanRateRange.Max}
	}
	if ml.LoanPointRange != nil {
		out.LoanPointRange = &admin.RangeFloat{Min: ml.LoanPointRange.Min, Max: ml.LoanPointRange.Max}
	}
	if out.LoanAmountRange == nil && out.LoanTermRange == nil && out.LoanRateRange == nil && out.LoanPointRange == nil {
		return nil
	}
	return out
}

func rawLenderFromCommitProto(p *admin.LenderDiscoverCommitLender) domainlenderdiscover.RawLender {
	if p == nil {
		return domainlenderdiscover.RawLender{}
	}
	wire := struct {
		Name         string        `json:"name"`
		ContactEmail string        `json:"contact_email"`
		ContactPhone string        `json:"contact_phone"`
		ContactName  string        `json:"contact_name"`
		Website      string        `json:"website"`
		Tagline      string        `json:"tagline"`
		Logo         string        `json:"logo"`
		Address      model.Address `json:"address"`
		Metadata     struct {
			Lender *struct {
				LoanAmountRange *model.RangeFloat64 `json:"loan_amount_range"`
				LoanTermRange   *model.RangeInt     `json:"loan_term_range"`
				LoanRateRange   *model.RangeFloat64 `json:"loan_rate_range"`
				LoanPointRange  *model.RangeFloat64 `json:"loan_point_range"`
			} `json:"lender"`
		} `json:"metadata"`
	}{
		Name:         p.GetName(),
		ContactEmail: p.GetContactEmail(),
		ContactPhone: p.GetContactPhone(),
		ContactName:  p.GetContactName(),
		Website:      p.GetWebsite(),
		Tagline:      p.GetTagline(),
		Logo:         p.GetLogo(),
	}
	if a := p.GetAddress(); a != nil {
		wire.Address = model.Address{
			State:    a.GetState(),
			City:     a.GetCity(),
			Address1: a.GetAddress1(),
			Address2: a.GetAddress2(),
			ZipCode:  a.GetZipCode(),
		}
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
		wire.Metadata.Lender = inner
	}
	b, err := json.Marshal(wire)
	if err != nil {
		return domainlenderdiscover.RawLender{}
	}
	var out domainlenderdiscover.RawLender
	_ = json.Unmarshal(b, &out)
	return out
}

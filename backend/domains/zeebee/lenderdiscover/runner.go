package lenderdiscover

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/logic"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/repo"
)

const (
	TriggerManual    = "manual"
	TriggerScheduled = "scheduled"

	DedupeNew     = "new"
	DedupeExisted = "existed"

	RunStatusRunning = "running"
	RunStatusSuccess = "success"
	RunStatusFailed  = "failed"

	ItemInserted         = "inserted"
	ItemSkippedDuplicate = "skipped_duplicate"
	ItemError            = "error"

	DefaultPlaceholderLogoURL = "https://placehold.co/64x64/png?text=L"
)

// Runner executes discover jobs (scheduled auto-insert; manual preview + commit).
type Runner struct {
	Repos              *repo.AllRepos
	TemplateLogic      *logic.OrganizationTemplateLogic
	SettingLogic       *logic.OrganizationSettingLogic
	OpenAIKey          string
	PlaceholderLogoURL string
}

// RawLender matches the JSON shape returned by the model (and sent on commit).
type RawLender struct {
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
}

type lendersEnvelope struct {
	Lenders []RawLender `json:"lenders"`
}

// PreviewCandidate is one row for admin review (no DB writes).
type PreviewCandidate struct {
	DedupeStatus  string
	NormalizedKey string
	Raw           RawLender
}

const systemPrompt = `You help build a comprehensive list of private mortgage lenders, hard-money lenders, bridge lenders, and non-QM lenders operating in the United States.
Output must be a single JSON object only (no markdown fences). The JSON must match the shape described in the user message.
Return as many distinct, real lenders as you can (aim for 30+). Include both well-known national lenders and smaller regional ones.
If you are uncertain about a field, use empty strings or zeros. Do not invent data.`

func buildUserPrompt(settings *model.LenderDiscoverSettings) string {
	scopeClause := "Search nationwide across all US states."
	if s := strings.TrimSpace(settings.ScopeState); s != "" {
		scopeClause = fmt.Sprintf("Focus on lenders operating in or headquartered in: %s. Also include major nationwide lenders.", s)
	}
	extra := strings.TrimSpace(settings.ScopeExtra)
	if extra == "" {
		extra = "private mortgage lenders, hard money lenders, bridge lenders, non-QM lenders, and private lending companies"
	}
	adminInstructions := strings.TrimSpace(settings.UserPrompt)
	if adminInstructions == "" {
		adminInstructions = "Return only factual-sounding entries; leave unknown fields empty. Output must follow the requested JSON shape."
	}

	return fmt.Sprintf(
		`%s
Scope: %s

Admin instructions:
%s

Return JSON with this exact structure:
{"lenders":[{"name":"required company name","contact_email":"","contact_phone":"","contact_name":"","website":"","tagline":"","logo":"","address":{"state":"","city":"","address1":"","address2":"","zip_code":""},"metadata":{"lender":{"loan_amount_range":{"min":0,"max":0},"loan_term_range":{"min":0,"max":0},"loan_rate_range":{"min":0,"max":0},"loan_point_range":{"min":0,"max":0}}}}]}

Rules:
- "lenders" must be an array (use [] if none).
- Each lender must have a distinct "name".
- Return at least 20 lenders, more is better.
- Include a mix of well-known national lenders and regional/local ones.
- Fill in website, contact_email, contact_phone when known.`,
		scopeClause,
		extra,
		adminInstructions,
	)
}

func (r *Runner) fetchAndParseLenders(ctx context.Context, settings *model.LenderDiscoverSettings) (*lendersEnvelope, error) {
	userPrompt := buildUserPrompt(settings)
	jsonStr, err := CallOpenAIJSON(ctx, r.OpenAIKey, settings.OpenaiModel, systemPrompt, userPrompt)
	if err != nil {
		return nil, err
	}
	var env lendersEnvelope
	if err := json.Unmarshal([]byte(jsonStr), &env); err != nil {
		return nil, errorx.WithStack(err)
	}
	return &env, nil
}

func (r *Runner) loadDedupeSet(ctx context.Context) (map[string]struct{}, error) {
	existing, err := r.Repos.OrganizationRepo.AllByTypeAllStatuses(ctx, constants.UserTypeLender)
	if err != nil {
		return nil, err
	}
	seen := make(map[string]struct{}, len(existing))
	for _, o := range existing {
		k := NormalizeCompanyName(o.Name)
		if len(k) > 0 {
			seen[k] = struct{}{}
		}
	}
	return seen, nil
}

func classifyLenders(env *lendersEnvelope, dbSeen map[string]struct{}) []PreviewCandidate {
	working := make(map[string]struct{}, len(dbSeen)+len(env.Lenders))
	for k := range dbSeen {
		working[k] = struct{}{}
	}
	var out []PreviewCandidate
	for _, raw := range env.Lenders {
		name := strings.TrimSpace(raw.Name)
		if name == "" {
			continue
		}
		key := NormalizeCompanyName(name)
		status := DedupeNew
		if _, ok := working[key]; ok {
			status = DedupeExisted
		} else {
			working[key] = struct{}{}
		}
		out = append(out, PreviewCandidate{DedupeStatus: status, NormalizedKey: key, Raw: raw})
	}
	return out
}

// Preview calls OpenAI, parses lenders, and classifies each row against existing orgs (and duplicates within the batch).
func (r *Runner) Preview(ctx context.Context) ([]PreviewCandidate, error) {
	if r.PlaceholderLogoURL == "" {
		r.PlaceholderLogoURL = DefaultPlaceholderLogoURL
	}
	settings, err := r.Repos.LenderDiscoverSettingsRepo.GetSingleton(ctx)
	if err != nil {
		return nil, err
	}
	env, err := r.fetchAndParseLenders(ctx, settings)
	if err != nil {
		return nil, err
	}
	seen, err := r.loadDedupeSet(ctx)
	if err != nil {
		return nil, err
	}
	return classifyLenders(env, seen), nil
}

// Commit inserts selected lenders as Not activated; re-checks dedupe against the DB at commit time.
func (r *Runner) Commit(ctx context.Context, lenders []RawLender) (run *model.LenderDiscoverRun, err error) {
	if r.PlaceholderLogoURL == "" {
		r.PlaceholderLogoURL = DefaultPlaceholderLogoURL
	}
	settings, err := r.Repos.LenderDiscoverSettingsRepo.GetSingleton(ctx)
	if err != nil {
		return nil, err
	}
	scopeSnap := map[string]interface{}{
		"state": settings.ScopeState,
		"extra": settings.ScopeExtra,
	}
	promptSnap := settings.UserPrompt

	run = &model.LenderDiscoverRun{
		TriggerType:    TriggerManual,
		Status:         RunStatusRunning,
		StartedAt:      time.Now().UTC(),
		Stats:          map[string]interface{}{},
		PromptSnapshot: promptSnap,
		ScopeSnapshot:  scopeSnap,
	}
	run, err = r.Repos.LenderDiscoverRunRepo.Create(ctx, run)
	if err != nil {
		return nil, err
	}

	finish := func(status, errMsg string, stats map[string]interface{}) {
		now := time.Now().UTC()
		_ = r.Repos.LenderDiscoverRunRepo.UpdateById(context.WithoutCancel(ctx), run.Id, map[string]interface{}{
			"status":          status,
			"finished_at":     &now,
			"error_message":   errMsg,
			"stats":           stats,
			"prompt_snapshot": promptSnap,
			"scope_snapshot":  scopeSnap,
		})
		run.Status = status
		run.FinishedAt = &now
		run.ErrorMessage = errMsg
		run.Stats = stats
	}

	seen, err := r.loadDedupeSet(ctx)
	if err != nil {
		finish(RunStatusFailed, err.Error(), map[string]interface{}{})
		return run, err
	}

	stats := map[string]interface{}{
		"requested_count":         len(lenders),
		"inserted_count":          0,
		"skipped_duplicate_count": 0,
		"error_count":             0,
	}
	insertedN, skippedN, errN := r.processLenderList(ctx, run, settings, seen, lenders)
	stats["inserted_count"] = insertedN
	stats["skipped_duplicate_count"] = skippedN
	stats["error_count"] = errN
	finish(RunStatusSuccess, "", stats)
	return run, nil
}

// Run is for the scheduler only: OpenAI + insert all non-duplicate lenders.
func (r *Runner) Run(ctx context.Context, trigger string) (run *model.LenderDiscoverRun, err error) {
	if trigger != TriggerScheduled {
		return nil, errorx.New("manual discover: use Preview then Commit")
	}
	if r.PlaceholderLogoURL == "" {
		r.PlaceholderLogoURL = DefaultPlaceholderLogoURL
	}
	settings, err := r.Repos.LenderDiscoverSettingsRepo.GetSingleton(ctx)
	if err != nil {
		return nil, err
	}

	scopeSnap := map[string]interface{}{
		"state": settings.ScopeState,
		"extra": settings.ScopeExtra,
	}
	promptSnap := settings.UserPrompt

	run = &model.LenderDiscoverRun{
		TriggerType:    trigger,
		Status:         RunStatusRunning,
		StartedAt:      time.Now().UTC(),
		Stats:          map[string]interface{}{},
		PromptSnapshot: promptSnap,
		ScopeSnapshot:  scopeSnap,
	}
	run, err = r.Repos.LenderDiscoverRunRepo.Create(ctx, run)
	if err != nil {
		return nil, err
	}

	finish := func(status, errMsg string, stats map[string]interface{}) {
		now := time.Now().UTC()
		_ = r.Repos.LenderDiscoverRunRepo.UpdateById(context.WithoutCancel(ctx), run.Id, map[string]interface{}{
			"status":          status,
			"finished_at":     &now,
			"error_message":   errMsg,
			"stats":           stats,
			"prompt_snapshot": promptSnap,
			"scope_snapshot":  scopeSnap,
		})
		run.Status = status
		run.FinishedAt = &now
		run.ErrorMessage = errMsg
		run.Stats = stats
	}

	env, err := r.fetchAndParseLenders(ctx, settings)
	if err != nil {
		finish(RunStatusFailed, err.Error(), map[string]interface{}{
			"requested_count":         0,
			"inserted_count":          0,
			"skipped_duplicate_count": 0,
			"error_count":             0,
		})
		return run, err
	}

	seen, err := r.loadDedupeSet(ctx)
	if err != nil {
		finish(RunStatusFailed, err.Error(), map[string]interface{}{})
		return run, err
	}

	stats := map[string]interface{}{
		"requested_count":         len(env.Lenders),
		"inserted_count":          0,
		"skipped_duplicate_count": 0,
		"error_count":             0,
	}
	insertedN, skippedN, errN := r.processLenderList(ctx, run, settings, seen, env.Lenders)
	stats["inserted_count"] = insertedN
	stats["skipped_duplicate_count"] = skippedN
	stats["error_count"] = errN
	finish(RunStatusSuccess, "", stats)
	return run, nil
}

func (r *Runner) processLenderList(ctx context.Context, run *model.LenderDiscoverRun, settings *model.LenderDiscoverSettings, seen map[string]struct{}, lenders []RawLender) (insertedN, skippedN, errN int) {
	for _, raw := range lenders {
		name := strings.TrimSpace(raw.Name)
		if len(name) == 0 {
			continue
		}
		key := NormalizeCompanyName(name)
		if _, dup := seen[key]; dup {
			skippedN++
			_ = r.Repos.LenderDiscoverRunItemRepo.Create(ctx, &model.LenderDiscoverRunItem{
				RunId:         run.Id,
				RawName:       name,
				NormalizedKey: key,
				Action:        ItemSkippedDuplicate,
			})
			continue
		}

		org := r.buildOrganization(raw)
		var created *model.Organization
		txErr := r.Repos.OrganizationRepo.Transaction(ctx, func(tx *gorm.DB) error {
			o, inner := r.Repos.OrganizationRepo.CreateWithTx(ctx, tx, org, false)
			if inner != nil {
				return inner
			}
			created = o
			if inner = r.TemplateLogic.SyncByOrganizationIdWithTx(ctx, tx, o); inner != nil {
				return inner
			}
			if inner = r.SettingLogic.InitLenderOrganizationSettingWithTx(ctx, tx, o); inner != nil {
				return inner
			}
			return nil
		})
		if txErr != nil {
			errN++
			_ = r.Repos.LenderDiscoverRunItemRepo.Create(ctx, &model.LenderDiscoverRunItem{
				RunId:         run.Id,
				RawName:       name,
				NormalizedKey: key,
				Action:        ItemError,
				ErrorDetail:   txErr.Error(),
			})
			continue
		}

		seen[key] = struct{}{}
		insertedN++
		_ = r.Repos.LenderDiscoverRunItemRepo.Create(ctx, &model.LenderDiscoverRunItem{
			RunId:          run.Id,
			RawName:        name,
			NormalizedKey:  key,
			Action:         ItemInserted,
			OrganizationId: created.Id,
		})
	}
	return insertedN, skippedN, errN
}

func (r *Runner) buildOrganization(raw RawLender) *model.Organization {
	email := strings.TrimSpace(raw.ContactEmail)
	if email == "" {
		email = "pending-discover@placeholder.local"
	}
	phone := strings.TrimSpace(raw.ContactPhone)
	if phone == "" {
		phone = "0000000000"
	}
	contactName := strings.TrimSpace(raw.ContactName)
	if contactName == "" {
		contactName = "Pending"
	}
	addr := raw.Address
	logo := strings.TrimSpace(raw.Logo)
	if logo == "" {
		logo = r.PlaceholderLogoURL
	}
	meta := &model.OrganizationMetadata{}
	if raw.Metadata.Lender != nil {
		meta.Lender = &model.OrganizationLenderMetadata{}
		if raw.Metadata.Lender.LoanAmountRange != nil {
			meta.Lender.LoanAmountRange = *raw.Metadata.Lender.LoanAmountRange
		}
		if raw.Metadata.Lender.LoanTermRange != nil {
			meta.Lender.LoanTermRange = *raw.Metadata.Lender.LoanTermRange
		}
		if raw.Metadata.Lender.LoanRateRange != nil {
			meta.Lender.LoanRateRange = *raw.Metadata.Lender.LoanRateRange
		}
		if raw.Metadata.Lender.LoanPointRange != nil {
			meta.Lender.LoanPointRange = *raw.Metadata.Lender.LoanPointRange
		}
	}

	return &model.Organization{
		Name:         strings.TrimSpace(raw.Name),
		Status:       constants.OrganizationStatusNotActivated,
		Source:       constants.OrganizationSourceSearch,
		Order:        1,
		OwnerId:      0,
		Address:      addr,
		ContactEmail: email,
		ContactPhone: phone,
		ContactName:  contactName,
		OwnerType:    constants.OrganizationOwnerTypeThirdParty,
		Type:         constants.UserTypeLender,
		Logo:         logo,
		Website:      strings.TrimSpace(raw.Website),
		Tagline:      strings.TrimSpace(raw.Tagline),
		Metadata:     meta,
	}
}

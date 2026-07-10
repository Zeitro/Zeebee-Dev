package lenderdiscover

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"gorm.io/gorm"

	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/logic"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/repo"
)

// Config holds the lender discover settings stored in organization_settings JSONB.
type Config struct {
	ScheduleEnabled bool   `json:"schedule_enabled"`
	IntervalDays    int32  `json:"interval_days"`
	ScopeState      string `json:"scope_state"`
	ScopeExtra      string `json:"scope_extra"`
	UserPrompt      string `json:"user_prompt"`
	OpenaiModel     string `json:"openai_model"`
}

const (
	// DefaultUserPrompt matches the seed in deploy/zeebee/sql/20260404120000_lender_discover.sql and buildUserPrompt empty fallback.
	DefaultUserPrompt         = "Return only factual-sounding entries; leave unknown fields empty. Output must follow the requested JSON shape."
	defaultOpenaiModel        = "gpt-4o-mini"
	defaultIntervalDays int32 = 7
)

func applyDiscoverConfigDefaults(cfg *Config) {
	if cfg.IntervalDays < 1 {
		cfg.IntervalDays = defaultIntervalDays
	}
	if strings.TrimSpace(cfg.OpenaiModel) == "" {
		cfg.OpenaiModel = defaultOpenaiModel
	}
	if strings.TrimSpace(cfg.UserPrompt) == "" {
		cfg.UserPrompt = DefaultUserPrompt
	}
}

// LoadConfig reads discover settings from organization_settings (default public org id, purpose=LenderDiscover).
// Returns the parsed config and the row ID (for subsequent updates).
func LoadConfig(ctx context.Context, settingRepo *repo.OrganizationSettingRepo) (*Config, int64, error) {
	row, err := settingRepo.DetailByPurpose(ctx, constants.OrganizationIdPublicDiscoverSettings, constants.OrganizationPurposeLenderDiscover, false)
	if err != nil {
		return nil, 0, errorx.WithStack(err)
	}
	cfg := &Config{IntervalDays: defaultIntervalDays, OpenaiModel: defaultOpenaiModel}
	raw, jsonErr := json.Marshal(row.Setting)
	if jsonErr != nil {
		return nil, 0, errorx.WithStack(jsonErr)
	}
	if jsonErr = json.Unmarshal(raw, cfg); jsonErr != nil {
		return nil, 0, errorx.WithStack(jsonErr)
	}
	applyDiscoverConfigDefaults(cfg)
	return cfg, row.Id, nil
}

// SaveConfig writes discover settings back to the organization_settings row.
func SaveConfig(ctx context.Context, settingRepo *repo.OrganizationSettingRepo, rowId int64, cfg *Config) error {
	raw, err := json.Marshal(cfg)
	if err != nil {
		return errorx.WithStack(err)
	}
	// Do not pass map[string]interface{} for jsonb: pgx cannot encode it (OID 0). Bind JSON text as jsonb.
	err = settingRepo.Client(ctx).Where("id = ?", rowId).Updates(map[string]interface{}{
		"setting": gorm.Expr("CAST(? AS jsonb)", string(raw)),
	}).Error
	return errorx.WithStack(err)
}

const (
	TriggerManual    = "manual"
	TriggerScheduled = "scheduled"

	DedupeNew     = "new"
	DedupeExisted = "existed"

	RunStatusSuccess = "success"
	RunStatusFailed  = "failed"

	DefaultPlaceholderLogoURL = "https://placehold.co/64x64/png?text=L"
)

// RunResult is a lightweight in-memory result returned by Run/Commit (no DB persistence).
type RunResult struct {
	TriggerType string
	Status      string
	ErrorMsg    string
	Stats       RunStats
}

type RunStats struct {
	RequestedCount        int
	InsertedCount         int
	SkippedDuplicateCount int
	ErrorCount            int
}

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

func buildUserPrompt(cfg *Config) string {
	scopeClause := "Search nationwide across all US states."
	if s := strings.TrimSpace(cfg.ScopeState); s != "" {
		scopeClause = fmt.Sprintf("Focus on lenders operating in or headquartered in: %s. Also include major nationwide lenders.", s)
	}
	extra := strings.TrimSpace(cfg.ScopeExtra)
	if extra == "" {
		extra = "private mortgage lenders, hard money lenders, bridge lenders, non-QM lenders, and private lending companies"
	}
	adminInstructions := strings.TrimSpace(cfg.UserPrompt)
	if adminInstructions == "" {
		adminInstructions = DefaultUserPrompt
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

func (r *Runner) fetchAndParseLenders(ctx context.Context, cfg *Config) (*lendersEnvelope, error) {
	userPrompt := buildUserPrompt(cfg)
	jsonStr, err := CallOpenAIJSON(ctx, r.OpenAIKey, cfg.OpenaiModel, systemPrompt, userPrompt)
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
	cfg, _, err := LoadConfig(ctx, r.Repos.OrganizationSettingRepo)
	if err != nil {
		return nil, err
	}
	env, err := r.fetchAndParseLenders(ctx, cfg)
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
func (r *Runner) Commit(ctx context.Context, lenders []RawLender) (*RunResult, error) {
	if r.PlaceholderLogoURL == "" {
		r.PlaceholderLogoURL = DefaultPlaceholderLogoURL
	}

	seen, err := r.loadDedupeSet(ctx)
	if err != nil {
		return &RunResult{TriggerType: TriggerManual, Status: RunStatusFailed, ErrorMsg: err.Error()}, err
	}

	insertedN, skippedN, errN := r.processLenderList(ctx, seen, lenders)
	return &RunResult{
		TriggerType: TriggerManual,
		Status:      RunStatusSuccess,
		Stats: RunStats{
			RequestedCount:        len(lenders),
			InsertedCount:         insertedN,
			SkippedDuplicateCount: skippedN,
			ErrorCount:            errN,
		},
	}, nil
}

// Run is for the scheduler only: OpenAI + insert all non-duplicate lenders.
func (r *Runner) Run(ctx context.Context, trigger string) (*RunResult, error) {
	if trigger != TriggerScheduled {
		return nil, errorx.New("manual discover: use Preview then Commit")
	}
	if r.PlaceholderLogoURL == "" {
		r.PlaceholderLogoURL = DefaultPlaceholderLogoURL
	}
	cfg, _, err := LoadConfig(ctx, r.Repos.OrganizationSettingRepo)
	if err != nil {
		return nil, err
	}

	env, err := r.fetchAndParseLenders(ctx, cfg)
	if err != nil {
		return &RunResult{TriggerType: trigger, Status: RunStatusFailed, ErrorMsg: err.Error()}, err
	}

	seen, err := r.loadDedupeSet(ctx)
	if err != nil {
		return &RunResult{TriggerType: trigger, Status: RunStatusFailed, ErrorMsg: err.Error()}, err
	}

	insertedN, skippedN, errN := r.processLenderList(ctx, seen, env.Lenders)
	return &RunResult{
		TriggerType: trigger,
		Status:      RunStatusSuccess,
		Stats: RunStats{
			RequestedCount:        len(env.Lenders),
			InsertedCount:         insertedN,
			SkippedDuplicateCount: skippedN,
			ErrorCount:            errN,
		},
	}, nil
}

func (r *Runner) processLenderList(ctx context.Context, seen map[string]struct{}, lenders []RawLender) (insertedN, skippedN, errN int) {
	for _, raw := range lenders {
		name := strings.TrimSpace(raw.Name)
		if len(name) == 0 {
			continue
		}
		key := NormalizeCompanyName(name)
		if _, dup := seen[key]; dup {
			skippedN++
			continue
		}

		org := r.buildOrganization(raw)
		txErr := r.Repos.OrganizationRepo.Transaction(ctx, func(tx *gorm.DB) error {
			o, inner := r.Repos.OrganizationRepo.CreateWithTx(ctx, tx, org, false)
			if inner != nil {
				return inner
			}
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
			continue
		}

		seen[key] = struct{}{}
		insertedN++
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

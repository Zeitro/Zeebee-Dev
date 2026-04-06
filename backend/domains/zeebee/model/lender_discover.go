package model

import "time"

// LenderDiscoverSettings is a singleton row (id=1) for OpenAI discover jobs.
type LenderDiscoverSettings struct {
	Base
	ScheduleEnabled  bool   `gorm:"column:schedule_enabled;not null;default:false" json:"schedule_enabled"`
	IntervalDays     int32  `gorm:"column:interval_days;not null;default:7" json:"interval_days"`
	ScopeState       string `gorm:"column:scope_state;type:varchar(8);not null;default:''" json:"scope_state"`
	ScopeExtra       string `gorm:"column:scope_extra;type:text;not null;default:''" json:"scope_extra"`
	UserPrompt       string `gorm:"column:user_prompt;type:text;not null;default:''" json:"user_prompt"`
	OpenaiModel      string `gorm:"column:openai_model;type:varchar(64);not null;default:'gpt-4o-mini'" json:"openai_model"`
	UpdatedByUserId  int64  `gorm:"column:updated_by_user_id;not null;default:0" json:"updated_by_user_id"`
}

func (LenderDiscoverSettings) TableName() string {
	return "lender_discover_settings"
}

// LenderDiscoverRun records one execution of the discover pipeline.
type LenderDiscoverRun struct {
	Base
	TriggerType     string                 `gorm:"column:trigger_type;type:varchar(32);not null;default:''" json:"trigger_type"`
	Status          string                 `gorm:"column:status;type:varchar(32);not null;default:''" json:"status"`
	StartedAt       time.Time              `gorm:"column:started_at;not null;default:now()" json:"started_at"`
	FinishedAt      *time.Time             `gorm:"column:finished_at" json:"finished_at"`
	ErrorMessage    string                 `gorm:"column:error_message;type:text" json:"error_message"`
	Stats           map[string]interface{} `gorm:"column:stats;type:jsonb;serializer:json;not null;default:'{}'" json:"stats"`
	PromptSnapshot  string                 `gorm:"column:prompt_snapshot;type:text;not null;default:''" json:"prompt_snapshot"`
	ScopeSnapshot   map[string]interface{} `gorm:"column:scope_snapshot;type:jsonb;serializer:json;not null;default:'{}'" json:"scope_snapshot"`
}

func (LenderDiscoverRun) TableName() string {
	return "lender_discover_runs"
}

// LenderDiscoverRunItem is one lender row processed in a run.
type LenderDiscoverRunItem struct {
	Base
	RunId          int64  `gorm:"column:run_id;not null;index" json:"run_id"`
	RawName        string `gorm:"column:raw_name;type:text;not null;default:''" json:"raw_name"`
	NormalizedKey  string `gorm:"column:normalized_key;type:text;not null;default:''" json:"normalized_key"`
	Action         string `gorm:"column:action;type:varchar(32);not null;default:''" json:"action"`
	OrganizationId int64  `gorm:"column:organization_id;default:0" json:"organization_id"`
	ErrorDetail    string `gorm:"column:error_detail;type:text;not null;default:''" json:"error_detail"`
}

func (LenderDiscoverRunItem) TableName() string {
	return "lender_discover_run_items"
}

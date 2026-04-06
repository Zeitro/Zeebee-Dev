package repo

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/model"
)

type LenderDiscoverSettingsRepo struct {
	db *gorm.DB
}

func NewLenderDiscoverSettingsRepo(db *gorm.DB) *LenderDiscoverSettingsRepo {
	return &LenderDiscoverSettingsRepo{db: db}
}

func (r *LenderDiscoverSettingsRepo) client(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

// GetSingleton returns the first settings row (bootstrap migration inserts one).
func (r *LenderDiscoverSettingsRepo) GetSingleton(ctx context.Context) (out *model.LenderDiscoverSettings, err error) {
	err = r.client(ctx).Model(&model.LenderDiscoverSettings{}).Order("id asc").Limit(1).First(&out).Error
	return out, errorx.WithStack(err)
}

func (r *LenderDiscoverSettingsRepo) UpdateById(ctx context.Context, id int64, cols map[string]interface{}) error {
	err := r.client(ctx).Model(&model.LenderDiscoverSettings{}).Where("id = ?", id).Updates(cols).Error
	return errorx.WithStack(err)
}

type LenderDiscoverRunRepo struct {
	db *gorm.DB
}

func NewLenderDiscoverRunRepo(db *gorm.DB) *LenderDiscoverRunRepo {
	return &LenderDiscoverRunRepo{db: db}
}

func (r *LenderDiscoverRunRepo) client(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

func (r *LenderDiscoverRunRepo) Create(ctx context.Context, row *model.LenderDiscoverRun) (*model.LenderDiscoverRun, error) {
	err := r.client(ctx).Create(row).Error
	return row, errorx.WithStack(err)
}

func (r *LenderDiscoverRunRepo) UpdateById(ctx context.Context, id int64, cols map[string]interface{}) error {
	err := r.client(ctx).Model(&model.LenderDiscoverRun{}).Where("id = ?", id).Updates(cols).Error
	return errorx.WithStack(err)
}

func (r *LenderDiscoverRunRepo) GetById(ctx context.Context, id int64) (out *model.LenderDiscoverRun, err error) {
	err = r.client(ctx).Model(&model.LenderDiscoverRun{}).Where("id = ?", id).First(&out).Error
	return out, errorx.WithStack(err)
}

func (r *LenderDiscoverRunRepo) Latest(ctx context.Context) (out *model.LenderDiscoverRun, err error) {
	err = r.client(ctx).Model(&model.LenderDiscoverRun{}).Order("started_at desc").Limit(1).First(&out).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errorx.WithStack(err)
	}
	return out, nil
}

// LastSuccessfulScheduledFinished returns finished_at of the latest successful scheduled run, if any.
func (r *LenderDiscoverRunRepo) LastSuccessfulScheduledFinished(ctx context.Context) (t *time.Time, err error) {
	var run model.LenderDiscoverRun
	err = r.client(ctx).
		Model(&model.LenderDiscoverRun{}).
		Where("trigger_type = ? AND status = ?", "scheduled", "success").
		Where("finished_at IS NOT NULL").
		Order("finished_at desc").
		Limit(1).
		First(&run).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errorx.WithStack(err)
	}
	return run.FinishedAt, nil
}

type LenderDiscoverRunItemRepo struct {
	db *gorm.DB
}

func NewLenderDiscoverRunItemRepo(db *gorm.DB) *LenderDiscoverRunItemRepo {
	return &LenderDiscoverRunItemRepo{db: db}
}

func (r *LenderDiscoverRunItemRepo) client(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

func (r *LenderDiscoverRunItemRepo) Create(ctx context.Context, row *model.LenderDiscoverRunItem) error {
	err := r.client(ctx).Create(row).Error
	return errorx.WithStack(err)
}

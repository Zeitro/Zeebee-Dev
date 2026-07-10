package repo

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
)

type StripeEventRepo struct {
	*gormx.BaseRepo[*model.StripeEvent]
}

func NewStripeEventRepo(db *gorm.DB) *StripeEventRepo {
	return &StripeEventRepo{
		BaseRepo: gormx.NewBaseRepo[*model.StripeEvent](db),
	}
}

// CreateIgnoreConflict inserts the event and reports whether a new row was
// created. A false return means the event id already exists (duplicate delivery).
func (repo *StripeEventRepo) CreateIgnoreConflict(ctx context.Context, event *model.StripeEvent) (created bool, err error) {
	result := repo.Client(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "event_id"}},
			DoNothing: true,
		}).
		Create(event)
	if result.Error != nil {
		return false, errorx.WithStack(result.Error)
	}
	return result.RowsAffected > 0, nil
}

func (repo *StripeEventRepo) FindOneByEventId(ctx context.Context, eventId string) (result *model.StripeEvent, err error) {
	err = repo.Client(ctx).Where("event_id = ?", eventId).First(&result).Error
	return result, errorx.WithStack(err)
}

func (repo *StripeEventRepo) UpdateStatusByEventId(ctx context.Context, eventId string, status constants.StripeEventStatus, errMsg string) error {
	err := repo.Client(ctx).Model(&model.StripeEvent{}).
		Where("event_id = ?", eventId).
		Updates(map[string]interface{}{
			"status": status,
			"error":  errMsg,
		}).Error
	return errorx.WithStack(err)
}

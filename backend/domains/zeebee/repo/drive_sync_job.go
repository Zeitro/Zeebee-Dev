package repo

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
)

type DriveSyncJobRepo struct {
	*gormx.BaseRepo[*model.DriveSyncJob]
}

func NewDriveSyncJobRepo(db *gorm.DB) *DriveSyncJobRepo {
	return &DriveSyncJobRepo{
		BaseRepo: gormx.NewBaseRepo[*model.DriveSyncJob](db),
	}
}

// GetLatestByTransaction returns the most recent sync job for a transaction, or
// (nil, nil) when none exists yet.
func (repo *DriveSyncJobRepo) GetLatestByTransaction(ctx context.Context, transactionId, organizationId int64, organizationType constants.UserType) (result *model.DriveSyncJob, err error) {
	err = repo.Client(ctx).
		Where("transaction_id = ? and organization_id = ? and organization_type = ?", transactionId, organizationId, organizationType).
		Order("id desc").
		First(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errorx.WithStack(err)
	}
	return result, nil
}

// UpdateResult persists the terminal outcome of a finished job.
func (repo *DriveSyncJobRepo) UpdateResult(ctx context.Context, id int64, columns map[string]interface{}) (err error) {
	return errorx.WithStack(repo.UpdateColumnsById(ctx, id, columns))
}

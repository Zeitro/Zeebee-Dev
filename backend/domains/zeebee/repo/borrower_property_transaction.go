package repo

import (
	"context"

	"gorm.io/gorm"

	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
)

type BorrowerPropertyTransactionRepo struct {
	*gormx.BaseRepo[*model.BorrowerPropertyTransaction]
}

func NewBorrowerPropertyTransactionRepo(db *gorm.DB) *BorrowerPropertyTransactionRepo {
	return &BorrowerPropertyTransactionRepo{
		BaseRepo: gormx.NewBaseRepo[*model.BorrowerPropertyTransaction](db),
	}
}

func (repo *BorrowerPropertyTransactionRepo) GetOne(ctx context.Context, id int64) (result *model.BorrowerPropertyTransaction, err error) {
	err = repo.Client(ctx).
		Where("id = ?", id).
		Preload("BorrowerProperty").
		First(&result).Error
	return result, errorx.WithStack(err)
}

func (repo *BorrowerPropertyTransactionRepo) GetOrganizationTransaction(ctx context.Context, id int64, organizationId int64) (result *model.BorrowerPropertyTransaction, err error) {
	err = repo.Client(ctx).
		Where("id = ? and organization_id = ?", id, organizationId).
		Preload("BorrowerProperty").
		First(&result).Error
	return result, errorx.WithStack(err)
}

func (repo *BorrowerPropertyTransactionRepo) GetByCode(ctx context.Context, code string) (result *model.BorrowerPropertyTransaction, err error) {
	err = repo.Client(ctx).
		Where("code = ?", code).
		First(&result).Error
	return result, errorx.WithStack(err)
}

func (repo *BorrowerPropertyTransactionRepo) AllByOrganization(ctx context.Context, organizationId int64) (result []*model.BorrowerPropertyTransaction, err error) {
	err = repo.Client(ctx).
		Where("organization_id = ? and status in ?", organizationId, []constants.BorrowerPropertyTransactionStatus{constants.BorrowerPropertyTransactionStatusPending, constants.BorrowerPropertyTransactionStatusClosed}).
		Preload("BorrowerProperty").
		Order("id desc").
		Find(&result).Error
	return result, errorx.WithStack(err)
}

func (repo *BorrowerPropertyTransactionRepo) UpdateOne(ctx context.Context, id int64, data *model.BorrowerPropertyTransaction) (result *model.BorrowerPropertyTransaction, err error) {
	err = repo.Client(ctx).
		Model(new(model.BorrowerPropertyTransaction)).
		Where("id = ?", id).
		Updates(data).Error
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	return data, nil
}

func (repo *BorrowerPropertyTransactionRepo) IsBoundToLoan(ctx context.Context, id int64) (bool, error) {
	return false, nil
}

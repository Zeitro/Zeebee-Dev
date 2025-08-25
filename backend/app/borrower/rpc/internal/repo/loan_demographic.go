package repo

import (
	"context"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"zeebee-backend/app/borrower/rpc/internal/repo/model"
	"zeebee-backend/common/gormx"
)

type LoanDemographicRepo struct {
	*gormx.BaseRepo[*model.LoanDemographic]
	m *model.LoanDemographic
}

func NewLoanDemographicRepo(db *gorm.DB) *LoanDemographicRepo {
	return &LoanDemographicRepo{
		BaseRepo: gormx.NewBaseRepo[*model.LoanDemographic](db),
		m:        new(model.LoanDemographic),
	}
}

// UpdateOneWithTx
func (repo *LoanDemographicRepo) UpdateOneWithTx(ctx context.Context, tx *gorm.DB, id int64, data *model.LoanDemographic) (result *model.LoanDemographic, err error) {
	err = tx.WithContext(ctx).Select("*").Omit("id", "created_at").Where("id = ?", id).Updates(data).Error
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return data, nil
}

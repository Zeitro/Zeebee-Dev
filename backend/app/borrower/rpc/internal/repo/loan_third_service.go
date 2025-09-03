package repo

import (
	"context"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"zeebee-backend/app/borrower/rpc/internal/repo/model"
	"zeebee-backend/common/gormx"
)

type LoanThirdServiceRepo struct {
	*gormx.BaseRepo[*model.LoanThirdService]
	m *model.LoanThirdService
}

func NewLoanThirdServiceRepo(db *gorm.DB) *LoanThirdServiceRepo {
	return &LoanThirdServiceRepo{
		BaseRepo: gormx.NewBaseRepo[*model.LoanThirdService](db),
		m:        new(model.LoanThirdService),
	}
}

// UpdateOneWithTx
func (repo *LoanThirdServiceRepo) UpdateOneWithTx(ctx context.Context, tx *gorm.DB, id int64, data *model.LoanThirdService) (result *model.LoanThirdService, err error) {
	err = tx.WithContext(ctx).Select("*").Omit("id", "created_at").Where("id = ?", id).Updates(data).Error
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return data, nil
}

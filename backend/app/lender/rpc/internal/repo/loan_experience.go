package repo

import (
	"context"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"zeebee-backend/app/lender/rpc/internal/repo/model"
	"zeebee-backend/common/gormx"
)

type LoanExperienceRepo struct {
	*gormx.BaseRepo[*model.LoanExperience]
	m *model.LoanExperience
}

func NewLoanExperienceRepo(db *gorm.DB) *LoanExperienceRepo {
	return &LoanExperienceRepo{
		BaseRepo: gormx.NewBaseRepo[*model.LoanExperience](db),
		m:        new(model.LoanExperience),
	}
}

// UpdateOneWithTx
func (repo *LoanExperienceRepo) UpdateOneWithTx(ctx context.Context, tx *gorm.DB, id int64, data *model.LoanExperience) (result *model.LoanExperience, err error) {
	err = tx.WithContext(ctx).Select("*").Omit("id", "created_at", "loan_id").Where("id = ?", id).Updates(data).Error
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return data, nil
}

func (repo *LoanExperienceRepo) DeleteByLoanIdWithTx(ctx context.Context, tx *gorm.DB, loanId int64) (err error) {
	err = tx.WithContext(ctx).
		Model(&model.LoanExperience{}).
		Where("loan_id = ?", loanId).
		Delete(&model.LoanExperience{}).
		Error
	return errors.WithStack(err)
}

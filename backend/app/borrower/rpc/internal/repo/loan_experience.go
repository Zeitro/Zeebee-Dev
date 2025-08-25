package repo

import (
	"context"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"zeebee-backend/app/borrower/rpc/internal/repo/model"
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
	err = tx.WithContext(ctx).Select("*").Omit("id", "created_at").Where("id = ?", id).Updates(data).Error
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return data, nil
}

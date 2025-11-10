package repo

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"zeebee-backend/app/lender/rpc/internal/repo/model"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/common/timex"
)

type BorrowerEntityRepo struct {
	*gormx.BaseRepo[*model.BorrowerEntity]
	m *model.BorrowerEntity
}

func NewBorrowerEntityRepo(db *gorm.DB) *BorrowerEntityRepo {
	return &BorrowerEntityRepo{
		BaseRepo: gormx.NewBaseRepo[*model.BorrowerEntity](db),
		m:        new(model.BorrowerEntity),
	}
}

func (repo *BorrowerEntityRepo) GetBorrowerEntityById(ctx context.Context, id int64) (result *model.BorrowerEntity, err error) {
	err = repo.Client(ctx).
		Where("id = ? and borrower_id = ?", id).
		First(&result).
		Error
	return result, errors.WithStack(err)
}

func (repo *BorrowerEntityRepo) TranToLoanBorrowerWithId(ctx context.Context, id int64) (result *model.LoanBorrower, err error) {
	var entity *model.BorrowerEntity
	entity, err = repo.GetBorrowerEntityById(ctx, id)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	result = &model.LoanBorrower{}
	err = copier.CopyWithOption(result, entity, copier.Option{
		CaseSensitive: true,
		DeepCopy:      true,
		Converters: []copier.TypeConverter{
			copierx.StrToSqlTimePtr(timex.LayoutUsStyle),
		},
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	result.Id = 0
	result.BorrowerEntityId = entity.Id

	return result, nil
}

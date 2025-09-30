package repo

import (
	"context"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"zeebee-backend/app/lender/rpc/internal/repo/model"
	"zeebee-backend/common/gormx"
)

type LoanDocumentRepo struct {
	*gormx.BaseRepo[*model.LoanDocument]
	m *model.LoanDocument
}

func NewLoanDocumentRepo(db *gorm.DB) *LoanDocumentRepo {
	return &LoanDocumentRepo{
		BaseRepo: gormx.NewBaseRepo[*model.LoanDocument](db),
		m:        new(model.LoanDocument),
	}
}

// GetFileTypesByLoanId
func (r *LoanDocumentRepo) GetFileTypesByLoanId(ctx context.Context, loanId int64) (fileTypes []string, err error) {
	err = r.BaseRepo.Client(ctx).
		Select("file_type").
		Where("loan_id = ?", loanId).
		Pluck("file_type", &fileTypes).Error
	return fileTypes, errors.WithStack(err)
}

func (repo *LoanDocumentRepo) BatchCreateWithTx(ctx context.Context, tx *gorm.DB, data []*model.LoanDocument) ([]*model.LoanDocument, error) {
	err := tx.WithContext(ctx).Create(&data).Error
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return data, nil
}

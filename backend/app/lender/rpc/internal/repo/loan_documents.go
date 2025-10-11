package repo

import (
	"context"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"zeebee-backend/app/lender/rpc/internal/repo/model"
	"zeebee-backend/common/constants"
	"zeebee-backend/common/gormx"
)

type LoanDocumentOption struct {
	Id       int64                          `json:"id"`
	Category constants.LoanDocumentCategory `json:"category"`
	FileType constants.LoanDocumentFileType `json:"file_type"`
}

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

// GetOptsByLoanId return loan document options
func (repo *LoanDocumentRepo) GetOptsByLoanId(ctx context.Context, loanId int64, categories []constants.LoanDocumentCategory) (opts []LoanDocumentOption, err error) {
	query := repo.Client(ctx).
		Select("id", "category", "file_type").
		Where("loan_id = ?", loanId)
	if len(categories) > 0 {
		query = query.Where("category in ?", categories)
	}
	err = query.Order("category asc").Scan(&opts).Error
	return opts, errors.WithStack(err)
}

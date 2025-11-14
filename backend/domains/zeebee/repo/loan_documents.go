package repo

import (
	"context"
	"database/sql"
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"zeebee-backend/common/constants"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/model"
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
	return fileTypes, errorx.WithStack(err)
}

func (repo *LoanDocumentRepo) BatchCreateWithTx(ctx context.Context, tx *gorm.DB, data []*model.LoanDocument) ([]*model.LoanDocument, error) {
	err := tx.WithContext(ctx).Create(&data).Error
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	return data, nil
}

func (repo *LoanDocumentRepo) DeleteLoanDocumentsByLoanIdWithTx(ctx context.Context, tx *gorm.DB, loanId int64) (err error) {
	var documentIds []int64
	// get document ids
	err = tx.WithContext(ctx).
		Model(&model.LoanDocument{}).
		Where("loan_id = ?", loanId).
		Pluck("id", &documentIds).Error
	if err != nil {
		return errorx.WithStack(err)
	}

	if len(documentIds) == 0 {
		return nil
	}

	// delete files
	err = tx.WithContext(ctx).
		Model(&model.File{}).
		Where("entity_type = ?", constants.FileEntityTypeLoanDocument).
		Where("entity_id in ?", documentIds).
		Delete(&model.File{}).
		Error
	if err != nil {
		return errorx.WithStack(err)
	}

	// delete loan documents
	err = tx.WithContext(ctx).
		Model(&model.LoanDocument{}).
		Where("loan_id = ?", loanId).
		Delete(&model.LoanDocument{}).
		Error
	return errorx.WithStack(err)
}

// SaveApplyFormDocument
func (repo *LoanDocumentRepo) SaveApplyFormDocumentWithTx(ctx context.Context, tx *gorm.DB, loanId int64, fileUrl string) (err error) {
	// get or create apply form loan document
	applyFormLoanDocument, _ := repo.GetLoanDocumentWithTx(ctx, tx, loanId, constants.LoanDocumentFileTypeApplyForm)
	if applyFormLoanDocument == nil || applyFormLoanDocument.Id == 0 {
		now := &sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
		applyFormLoanDocument = &model.LoanDocument{
			LoanId:         loanId,
			FileType:       constants.LoanDocumentFileTypeApplyForm,
			Category:       constants.LoanDocumentCategoryInit,
			Status:         constants.LoanDocumentStatusReadyForReview,
			LastUploadedAt: now,
		}
		err = tx.WithContext(ctx).
			Create(applyFormLoanDocument).
			Error
		if err != nil {
			return errorx.WithStack(err)
		}
	}
	// save file
	var file *model.File
	if len(applyFormLoanDocument.Files) == 0 {
		file = &model.File{
			EntityType: constants.FileEntityTypeLoanDocument,
			EntityId:   applyFormLoanDocument.Id,
		}
	} else {
		file = applyFormLoanDocument.Files[0]
	}
	file.FileUrl = fileUrl
	err = tx.WithContext(ctx).
		Save(file).
		Error

	return errorx.WithStack(err)
}

func (repo *LoanDocumentRepo) GetLoanDocumentWithTx(ctx context.Context, tx *gorm.DB, loanId int64, fileType constants.LoanDocumentFileType) (loanDocument *model.LoanDocument, err error) {
	loanDocument = &model.LoanDocument{}
	err = tx.WithContext(ctx).
		Model(&model.LoanDocument{}).
		Where("loan_id = ?", loanId).
		Where("file_type = ?", fileType).
		Preload("Files").
		First(&loanDocument).
		Error
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	if loanDocument == nil || loanDocument.Id == 0 {
		err = errorx.New("Loan Document Not Found")
	}
	return loanDocument, errorx.WithStack(err)
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

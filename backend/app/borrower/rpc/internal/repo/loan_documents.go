package repo

import (
	"context"
	"database/sql"
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"zeebee-backend/app/borrower/rpc/internal/repo/model"
	"zeebee-backend/common/constants"
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

func (repo *LoanDocumentRepo) DeleteLoanDocumentsByLoanIdWithTx(ctx context.Context, tx *gorm.DB, loanId int64) (err error) {
	var documentIds []int64
	// get document ids
	err = tx.WithContext(ctx).
		Model(&model.LoanDocument{}).
		Where("loan_id = ?", loanId).
		Pluck("id", &documentIds).Error
	if err != nil {
		return errors.WithStack(err)
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
		return errors.WithStack(err)
	}

	// delete loan documents
	err = tx.WithContext(ctx).
		Model(&model.LoanDocument{}).
		Where("loan_id = ?", loanId).
		Delete(&model.LoanDocument{}).
		Error
	return errors.WithStack(err)
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
			return errors.WithStack(err)
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

	return errors.WithStack(err)
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
		return nil, errors.WithStack(err)
	}
	if loanDocument == nil || loanDocument.Id == 0 {
		err = errors.New("loan document not found")
	}
	return loanDocument, errors.WithStack(err)
}

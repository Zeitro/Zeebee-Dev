package loanlogic

import (
	"context"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"zeebee-backend/app/borrower/rpc/internal/repo/model"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/constants"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type UpdateLoanDocumentsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateLoanDocumentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLoanDocumentsLogic {
	return &UpdateLoanDocumentsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateLoanDocumentsLogic) UpdateLoanDocuments(in *borrower.LoanUpdateDocumentsRequest) (resp *borrower.EmptyResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	// get loan detail
	loan, err := l.svcCtx.LoanRepo.GetBorrowerLoan(l.ctx, in.Data.Id, in.Data.BorrowerId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get loan")
	}
	if loan == nil {
		return nil, errors.WithStack(errors.New("loan not found"))
	}

	// check status
	if loan.Stage != constants.LoanStagePending {
		return nil, errors.WithStack(errors.New("loan is not in pending stage"))
	}
	if loan.PendingStatus != constants.LoanPendingStatusNone && loan.PendingStatus != constants.LoanPendingStatusQuoting {
		return nil, errors.WithStack(errors.New("loan is not in pending status"))
	}

	// check source
	if loan.Source != constants.LoanSourceSystem {
		return nil, errors.WithStack(errors.New("loan is not from system"))
	}

	err = l.svcCtx.LoanRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		var innerErr error
		oldDocuments := loan.Documents
		newDocuments := make([]*model.File, 0, len(in.Data.Documents))
		existingDocuments := make([]*model.File, 0, len(in.Data.Documents))

		for _, doc := range in.Data.Documents {
			doc := doc
			if constants.FileType(doc.FileType) == constants.FileTypeApplyForm {
				continue
			}
			if !constants.FileType(doc.FileType).IsValid() {
				return errors.WithStack(errors.New("invalid document type: " + doc.FileType))
			}

			item, ok := slice.FindBy(oldDocuments, func(index int, oldDoc *model.File) bool {
				return oldDoc.FileType == constants.FileType(doc.FileType)
			},
			)
			if ok {
				item.File = doc.File
				existingDocuments = append(existingDocuments, item)
			} else {
				newDocuments = append(newDocuments, &model.File{
					FileType: constants.FileType(doc.FileType),
					BizType:  constants.FileBizTypeLoan,
					RefId:    loan.Id,
					File:     doc.File,
				})
			}
		}

		// update documents
		for _, doc := range existingDocuments {
			_, innerErr = l.svcCtx.FileRepo.UpdateOneWithTx(l.ctx, tx, doc.Id, doc)
			if innerErr != nil {
				return innerErr
			}
		}

		// create  documents
		if len(newDocuments) > 0 {
			_, innerErr = l.svcCtx.FileRepo.BatchCreateWithTx(l.ctx, tx, newDocuments, false)
			if innerErr != nil {
				return innerErr
			}
		}
		return innerErr
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to update loan documents")
	}

	return &borrower.EmptyResponse{}, nil

}

package loanlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/repo/model"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type BatchCreateBorrowerLoanLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

type ImportLoanItem struct {
	Loan *model.Loan
	// Property   *model.Property
	Applicants []*model.Applicant
}

func NewBatchCreateBorrowerLoanLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchCreateBorrowerLoanLogic {
	return &BatchCreateBorrowerLoanLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BatchCreateBorrowerLoanLogic) BatchCreateBorrowerLoan(in *borrower.BatchCreateLoanRequest) (resp *borrower.EmptyResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	// createLogic := NewCreateBorrowerLoanLogic(l.ctx, l.svcCtx)
	// insertLoans := make([]ImportLoanItem, 0, len(in.Data))
	// for _, createLoan := range in.Data {
	// 	loan, property, applicants, err := createLogic.buildCreateModels(createLoan)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	insertLoans = append(insertLoans, ImportLoanItem{
	// 		Loan:       loan,
	// 		Property:   property,
	// 		Applicants: applicants,
	// 	})
	// }
	// // create loan
	// propertyIds := map[int64]struct{}{}
	// err = l.svcCtx.LoanRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
	// 	var innerErr error
	// 	for _, insertLoan := range insertLoans {
	// 		loan := insertLoan.Loan
	// 		applicants := insertLoan.Applicants
	// 		property := insertLoan.Property
	// 		property, innerErr = l.svcCtx.PropertyRepo.GetOrCreateWithTx(l.ctx, tx, insertLoan.Property)
	// 		if innerErr != nil {
	// 			return innerErr
	// 		}
	//
	// 		loan.PropertyId = property.Id
	// 		loan, innerErr = l.svcCtx.LoanRepo.CreateWithTx(l.ctx, tx, loan)
	// 		if innerErr != nil {
	// 			return innerErr
	// 		}
	//
	// 		applicants = slice.Map(applicants, func(index int, item *model.Applicant) *model.Applicant {
	// 			item.LoanId = loan.Id
	// 			item.PropertyId = property.Id
	// 			return item
	// 		})
	// 		_, innerErr = l.svcCtx.ApplicantRepo.BatchCreateWithTx(l.ctx, tx, applicants)
	// 		if innerErr != nil {
	// 			return innerErr
	// 		}
	// 		// applicant.LoanId = loan.Id
	// 		// applicant.PropertyId = property.Id
	// 		// applicant, innerErr = l.svcCtx.ApplicantRepo.CreateWithTx(l.ctx, tx, applicant)
	// 		// if innerErr != nil {
	// 		// 	return innerErr
	// 		// }
	// 		propertyIds[property.Id] = struct{}{}
	//
	// 	}
	// 	// sync total amount
	// 	for propertyId := range propertyIds {
	// 		innerErr = l.svcCtx.PropertyRepo.SyncTotalAmountWithTx(l.ctx, tx, propertyId)
	// 		if innerErr != nil {
	// 			return innerErr
	// 		}
	// 	}
	// 	return innerErr
	// })

	return &borrower.EmptyResponse{}, nil
}

package loanlogic

import (
	"context"

	"gorm.io/gorm"

	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/constants"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteBorrowerLoanLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteBorrowerLoanLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteBorrowerLoanLogic {
	return &DeleteBorrowerLoanLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// DeleteBorrowerLoan delete borrower loan
func (l *DeleteBorrowerLoanLogic) DeleteBorrowerLoan(in *borrower.LoanDeleteRequest) (resp *borrower.EmptyResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	loanModel, err := l.svcCtx.LoanRepo.GetBorrowerLoan(l.ctx, in.Data.Id, in.Data.OrganizationId)
	if err != nil {
		return nil, errorx.Wrap(err, "Failed To Get Loan")
	}

	if loanModel.Source != constants.LoanSourceBorrower && loanModel.Stage != constants.LoanStageDraft {
		return nil, errorx.WithStack(errorx.New("Loan Source Is Not 3rd Party Or Loan Is Not In Draft Stage"))
	}
	// todo check organization
	if loanModel.BorrowerOrganizationId != in.Data.OrganizationId {
		return nil, errorx.WithStack(errorx.New("Loan Does Not Belong To Borrower"))
	}

	err = l.svcCtx.LoanRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		var innerErr error
		var loanProperty *model.LoanProperty
		var propertyBase *model.BorrowerProperty
		loanProperty = loanModel.LoanProperty
		if loanProperty != nil {
			propertyBase = loanProperty.BorrowerProperty
		}

		// delete property_base
		if propertyBase != nil {
			var loanCount int64
			loanCount, innerErr = l.svcCtx.LoanPropertyRepo.GetLoanPropertyCountByBorrowerPropertyId(l.ctx, propertyBase.Id)
			if innerErr != nil {
				return innerErr
			}
			if loanCount == 1 {
				_, innerErr = l.svcCtx.BorrowerPropertyRepo.DeleteByIDWithTx(l.ctx, tx, propertyBase.Id)
				if innerErr != nil {
					return innerErr
				}
			}
		}

		// delete loan_property
		_, innerErr = l.svcCtx.LoanPropertyRepo.DeleteByIDWithTx(l.ctx, tx, loanProperty.Id)
		if innerErr != nil {
			return innerErr
		}

		// delete applicants
		for _, applicant := range loanModel.Applicants {
			_, innerErr = l.svcCtx.LoanApplicantRepo.DeleteByIDWithTx(l.ctx, tx, applicant.Id)
			if innerErr != nil {
				return innerErr
			}
		}

		// deleter loan borrower
		innerErr = l.svcCtx.LoanBorrowerRepo.DeleteByLoanIdWithTx(l.ctx, tx, loanModel.Id)
		if innerErr != nil {
			return innerErr
		}

		// delete loan documents
		innerErr = l.svcCtx.LoanDocumentRepo.DeleteLoanDocumentsByLoanIdWithTx(l.ctx, tx, loanModel.Id)
		if innerErr != nil {
			return innerErr
		}

		// delete loan borrower
		innerErr = l.svcCtx.LoanBorrowerRepo.DeleteByLoanIdWithTx(l.ctx, tx, loanModel.Id)
		if innerErr != nil {
			return innerErr
		}
		// delete loan demographics
		innerErr = l.svcCtx.LoanDemographicRepo.DeleteByLoanIdWithTx(l.ctx, tx, loanModel.Id)
		if innerErr != nil {
			return innerErr
		}

		// delete loan experience
		innerErr = l.svcCtx.LoanExperienceRepo.DeleteByLoanIdWithTx(l.ctx, tx, loanModel.Id)
		if innerErr != nil {
			return innerErr
		}
		// delete loan third service
		innerErr = l.svcCtx.LoanThirdServiceRepo.DeleteByLoanIdWithTx(l.ctx, tx, loanModel.Id)
		if innerErr != nil {
			return innerErr
		}

		// delete loan
		_, innerErr = l.svcCtx.LoanRepo.DeleteByIDWithTx(l.ctx, tx, loanModel.Id)
		if innerErr != nil {
			return innerErr
		}

		return innerErr
	})

	if err != nil {
		return nil, errorx.Wrap(err, "Failed To Delete Loan")
	}
	return &borrower.EmptyResponse{}, nil
}

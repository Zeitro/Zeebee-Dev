package loanlogic

import (
	"context"

	"gorm.io/gorm"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/proto"
)

type DeleteBorrowerLoanLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewDeleteBorrowerLoanLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteBorrowerLoanLogic {
	return &DeleteBorrowerLoanLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

// DeleteBorrowerLoan delete borrower loan
func (l *DeleteBorrowerLoanLogic) DeleteBorrowerLoan(in *borrower.LoanDeleteRequest) (resp *proto.EmptyResponse, err error) {

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx, constants.UserTypeBorrower)
	if err != nil {
		return nil, err
	}

	loanModel, err := l.svcCtx.LoanRepo.GetBorrowerLoan(l.ctx, in.Data.Id, claim.OrganizationId)
	if err != nil {
		return nil, errorx.Wrap(err, "fail to get loan")
	}

	if loanModel.Source != constants.LoanSourceBorrower && loanModel.Stage != constants.LoanStageDraft {
		return nil, errorx.WithStack(errorx.New("loan source is not 3rd party or loan is not in draft stage"))
	}
	// todo check organization
	if loanModel.BorrowerOrganizationId != claim.OrganizationId {
		return nil, errorx.WithStack(errorx.New("loan does not belong to borrower"))
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
		return nil, errorx.Wrap(err, "fail to delete loan")
	}
	return &proto.EmptyResponse{}, nil
}

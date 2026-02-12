package logic

import (
	"context"

	"gorm.io/gorm"

	"zeebee-backend/common/errorx"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/repo"
)

type LoanLogic struct {
	repos *repo.AllRepos
}

// NewLoanLogic creates a new instance of LoanLogic.
// It initializes the logic with the provided repository dependencies.
func NewLoanLogic(repos *repo.AllRepos) *LoanLogic {
	return &LoanLogic{
		repos: repos,
	}
}

// ReassignLender
func (l *LoanLogic) ReassignLender(ctx context.Context, id int64, lenderOrganizationId int64, claim utils.JwtUser) (err error) {
	loan, err := l.repos.LoanRepo.Detail(ctx, id, false)
	if err != nil {
		return err
	}
	if loan.Stage != constants.LoanStageDraft && loan.Source != constants.LoanSourcePreapproval {
		err = errorx.New("Only draft can be reassign")
		return err
	}
	if loan.LenderOrganizationId == lenderOrganizationId {
		err = errorx.New("Lender organization is same")
		return err
	}

	lenderOrganization, err := l.repos.OrganizationRepo.GetOrganization(ctx, lenderOrganizationId, constants.UserTypeLender, false)
	if err != nil {
		return err
	}
	if lenderOrganization.Status != constants.OrganizationStatusActive {
		err = errorx.New("Lender organization is not active")
		return err
	}

	err = l.repos.LoanRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var innerErr error
		originalLenderOrganizationId := loan.LenderOrganizationId

		loan.LenderOrganizationId = lenderOrganization.Id
		loan.LenderCompany = lenderOrganization.Name
		_, err = l.repos.LoanRepo.UpdateOneByIdWithTx(ctx, tx, loan.Id, loan)

		operationLogic := NewLoanOperationLogic(l.repos)
		input := dto.LoanOperationCreateDto{
			LoanId:     loan.Id,
			OperatorId: claim.OrganizationUserId,
			From: data.LoanOperationStatus{
				Stage:         loan.Stage,
				PendingStatus: loan.PendingStatus,
			},
			To: data.LoanOperationStatus{
				Stage:         loan.Stage,
				PendingStatus: loan.PendingStatus,
			},
			Action:        constants.LoanOperationActionReassignLender,
			OperationType: constants.LoanOperationOperatorTypeAdmin,
			Metadata: &data.LoanOperationMetadata{
				Comment:                      "Admin reassign lender",
				OriginalLenderOrganizationId: originalLenderOrganizationId,
				NewLenderOrganizationId:      lenderOrganizationId,
			},
		}
		innerErr = operationLogic.CreateLoanOperationLog(ctx, tx, input)

		return innerErr
	})

	return errorx.WithStack(err)
}

package loanlogic

import (
	"context"

	"github.com/duke-git/lancet/v2/mathutil"
	"github.com/duke-git/lancet/v2/slice"
	"gorm.io/gorm"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/timex"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
)

type UpdateBorrowerLoanLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewUpdateBorrowerLoanLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateBorrowerLoanLogic {
	return &UpdateBorrowerLoanLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *UpdateBorrowerLoanLogic) UpdateBorrowerLoan(in *borrower.LoanUpdateRequest) (resp *borrower.EmptyResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	loanModel, err := l.svcCtx.LoanRepo.GetBorrowerLoan(l.ctx, in.Data.Id, claim.OrganizationId)
	if err != nil {
		return nil, err
	}

	loan, loanProperty, applicants, loanBorrower, err := l.buildUpdateModels(in.Data, loanModel, claim)
	if err != nil {
		return nil, err
	}

	oldApplicants := loanModel.Applicants
	if len(oldApplicants) == 0 {
		return nil, errorx.WithStack(errorx.New("Applicant is not found"))
	}
	if loanModel.Source != constants.LoanSourceBorrower && loanModel.Source != constants.LoanSourcePreapproval {
		return nil, errorx.WithStack(errorx.New("Loan source is not third party"))
	}
	loanPropertyModel := loanModel.LoanProperty
	if loanPropertyModel == nil {
		return nil, errorx.WithStack(errorx.New("Loan property is not found"))
	}

	// create loan
	err = l.svcCtx.LoanRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		var innerErr error
		previousBorrowerPropertyId := loanPropertyModel.BorrowerPropertyId

		// update loan
		loan.Id = in.Data.Id
		loan.RemoveRelationships()
		_, innerErr = l.svcCtx.LoanRepo.UpdateOneWithTx(l.ctx, tx, in.Data.Id, loan)
		if innerErr != nil {
			return innerErr
		}

		// update loan property
		loanProperty.Id = loanPropertyModel.Id
		_, innerErr = l.svcCtx.LoanPropertyRepo.UpdateOneWithTx(l.ctx, tx, loanPropertyModel.Id, loanProperty)
		if innerErr != nil {
			return innerErr
		}

		// loan borrower
		if loanModel.LoanBorrower == nil {
			loanBorrower.LoanId = loanModel.Id
			loanBorrower, innerErr = l.svcCtx.LoanBorrowerRepo.CreateWithTx(l.ctx, tx, loanBorrower, false)
			if innerErr != nil {
				return innerErr
			}
		} else {
			loanBorrower.LoanId = loanModel.Id
			loanBorrower.Id = loanModel.LoanBorrower.Id
			_, innerErr = l.svcCtx.LoanBorrowerRepo.UpdateOneWithTx(l.ctx, tx, loanBorrower.Id, loanBorrower)
			if innerErr != nil {
				return innerErr
			}

		}

		// handle applicant
		removedApplicants := slice.DifferenceWith(oldApplicants, applicants, func(item1, item2 *model.LoanApplicant) bool {
			return item1.Id == item2.Id
		})
		removedApplicantIds := slice.Map(removedApplicants, func(index int, item *model.LoanApplicant) int64 {
			return item.Id
		})
		newApplicants := slice.DifferenceWith(applicants, oldApplicants, func(item1, item2 *model.LoanApplicant) bool {
			return item1.Id == item2.Id
		})
		newApplicants = slice.Map(newApplicants, func(index int, newApplicant *model.LoanApplicant) *model.LoanApplicant {
			newApplicant.LoanId = loan.Id
			newApplicant.LoanPropertyId = loanPropertyModel.Id
			newApplicant.LoanBorrowerId = loanBorrower.Id
			return newApplicant
		})

		updatedApplicants := slice.Filter(applicants, func(index int, newApplicant *model.LoanApplicant) bool {
			_, ok := slice.FindBy(oldApplicants, func(index int, oldApplicant *model.LoanApplicant) bool {
				return newApplicant.Id == oldApplicant.Id
			})
			return ok
		})

		updatedApplicants = slice.Map(updatedApplicants, func(index int, updateApplicant *model.LoanApplicant) *model.LoanApplicant {
			updateApplicant.LoanId = loan.Id
			updateApplicant.LoanPropertyId = loanPropertyModel.Id
			updateApplicant.LoanBorrowerId = loanBorrower.Id
			return updateApplicant
		})

		if len(removedApplicantIds) > 0 {
			_, innerErr = l.svcCtx.LoanApplicantRepo.DeleteByIDsWithTx(l.ctx, tx, removedApplicantIds)
			if innerErr != nil {
				return innerErr
			}
		}

		if len(newApplicants) > 0 {
			_, innerErr = l.svcCtx.LoanApplicantRepo.BatchCreateWithTx(l.ctx, tx, newApplicants, false)
			if innerErr != nil {
				return innerErr
			}
		}

		for _, updatedApplicant := range updatedApplicants {
			_, innerErr = l.svcCtx.LoanApplicantRepo.UpdateOneWithTx(l.ctx, tx, updatedApplicant.Id, updatedApplicant)
			if innerErr != nil {
				return innerErr
			}
		}

		// sync property total amount
		_ = l.svcCtx.BorrowerPropertyRepo.SyncTotalAmountAndLoanCountWithTx(l.ctx, tx, previousBorrowerPropertyId)
		if in.Data.LoanProperty.BorrowerPropertyId != previousBorrowerPropertyId {
			_ = l.svcCtx.BorrowerPropertyRepo.SyncTotalAmountAndLoanCountWithTx(l.ctx, tx, in.Data.LoanProperty.BorrowerPropertyId)
		}
		return innerErr
	})
	if err != nil {
		return nil, errorx.Wrap(err, "Failed To Update Loan")
	}

	return &borrower.EmptyResponse{}, nil
}

func (l *UpdateBorrowerLoanLogic) buildUpdateModels(in *borrower.LoanUpdateRequestData, oldLoan *model.Loan, claim utils.JwtUser) (loan *model.Loan, loanProperty *model.LoanProperty, applicants []*model.LoanApplicant, loanBorrower *model.LoanBorrower, err error) {

	// build applicants
	var totalOwnPercent int32
	for _, applicant := range in.Applicants {
		totalOwnPercent += applicant.OwnPercent
		borrowerApplicant, err := l.svcCtx.BorrowerApplicantRepo.GetBorrowerApplicantById(l.ctx, applicant.BorrowerApplicantId, claim.OrganizationId)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		applicantItem := &model.LoanApplicant{}
		applicantItem, err = copierx.Trans[model.BorrowerApplicant, model.LoanApplicant](borrowerApplicant, copierx.StrToSqlTimePtr(timex.LayoutUsStyle))
		if err != nil {
			return nil, nil, nil, nil, err
		}
		applicantItem.Id = applicant.Id
		applicantItem.IsPrimary = applicant.IsPrimary
		applicantItem.BorrowerApplicantId = applicant.BorrowerApplicantId
		applicantItem.InvitationStatus = constants.ApplicantInvitationStatusAccepted
		applicants = append(applicants, applicantItem)
	}
	if totalOwnPercent > 100 {
		return nil, nil, nil, nil, errorx.New("Total own percent must be less than 100")
	}

	slice.SortBy(applicants, func(a, b *model.LoanApplicant) bool {
		return b.IsPrimary
	})
	applicants = slice.UniqueByComparator(applicants, func(item *model.LoanApplicant, other *model.LoanApplicant) bool {
		return item.FullName == other.FullName
	})
	if len(applicants) == 0 {
		return nil, nil, nil, nil, errorx.New("No applicant")
	}
	primaryCount := slice.CountBy(applicants, func(index int, item *model.LoanApplicant) bool {
		return item.IsPrimary
	})
	if primaryCount == 0 {
		return nil, nil, nil, nil, errorx.New("No primary applicant")
	}
	if primaryCount > 1 {
		return nil, nil, nil, nil, errorx.New("Multiple primary applicants")
	}

	// LoanBorrower
	loanBorrower = &model.LoanBorrower{}
	loanBorrower.LoanId = in.Id

	// loan property
	borrowerProperty, err := l.svcCtx.BorrowerPropertyRepo.GetOneByIdForOrganization(l.ctx, in.LoanProperty.BorrowerPropertyId, claim.OrganizationId)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	loanProperty = &model.LoanProperty{}
	loanProperty, err = copierx.Trans[model.BorrowerProperty, model.LoanProperty](borrowerProperty, copierx.StrToSqlTimePtr(timex.LayoutUsStyle))
	if err != nil {
		return nil, nil, nil, nil, err
	}

	err = copierx.Assign[borrower.ThirdPartyLoanPropertyEditEntity, model.LoanProperty](in.LoanProperty, loanProperty, copierx.StrToSqlTimePtr(timex.LayoutUsStyle))
	if err != nil {
		return nil, nil, nil, nil, err
	}

	loanProperty.Id = 0
	if loanProperty.YearlyLiabilityExpenseAmount > 0 {
		loanProperty.Dscr = mathutil.RoundToFloat((loanProperty.EstYearlyRentalIncomeAmount/loanProperty.YearlyLiabilityExpenseAmount)*100, 2)
	}

	// loan
	loan = &model.Loan{}
	loan, err = copierx.Trans[model.Loan, model.Loan](oldLoan)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	err = copierx.Assign[borrower.LoanUpdateRequestData, model.Loan](in, loan, copierx.StrToSqlTimePtr(timex.LayoutUsStyle))
	if err != nil {
		return nil, nil, nil, nil, err
	}

	loan.RemoveRelationships()
	if loan.TotalAmount < loan.InitialAmount {
		return nil, nil, nil, nil, errorx.New("Total amount cannot be less than initial amount")
	}
	switch constants.LoanStage(in.Stage) {
	case constants.LoanStagePending:
		if loan.ClosingDate == nil {
			return nil, nil, nil, nil, errorx.New("Closing date is required")
		}
		loan.HoldbackAmount = loan.TotalAmount - loan.InitialAmount
		loan.PrincipalAmount = loan.TotalAmount
		loan.PendingStatus = constants.LoanPendingStatus(in.PendingStatus)

	case constants.LoanStageActive:
		if loan.MaturityDate == nil {
			return nil, nil, nil, nil, errorx.New("Maturity date is required")
		}
	case constants.LoanStageArchived:

		if loan.PaidOffDate == nil {
			return nil, nil, nil, nil, errorx.New("Paid off date is required")
		}
	default:
		return nil, nil, nil, nil, errorx.New("Invalid loan stage")
	}
	return
}

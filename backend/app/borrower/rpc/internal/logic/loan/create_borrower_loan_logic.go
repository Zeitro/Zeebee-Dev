package loanlogic

import (
	"context"

	"github.com/duke-git/lancet/v2/slice"
	"gorm.io/gorm"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/timex"
	"zeebee-backend/common/utils"
	"zeebee-backend/common/uuid"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
)

type CreateBorrowerLoanLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewCreateBorrowerLoanLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateBorrowerLoanLogic {
	return &CreateBorrowerLoanLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *CreateBorrowerLoanLogic) CreateBorrowerLoan(in *borrower.LoanCreateRequest) (resp *borrower.LoanCreateResponse, err error) {

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx, constants.UserTypeBorrower)
	if err != nil {
		return nil, err
	}

	loan, loanProperty, applicants, loanBorrower, err := l.buildCreateModels(in.Data, claim)
	if err != nil {
		return nil, err
	}

	// set transaction binding if provided
	if in.Data.BorrowerPropertyTransactionId > 0 {
		_, err = l.svcCtx.BorrowerPropertyTransactionRepo.GetOne(l.ctx, in.Data.BorrowerPropertyTransactionId)
		if err != nil {
			return nil, errorx.Wrap(err, "Transaction not found")
		}
		loan.BorrowerPropertyTransactionId = in.Data.BorrowerPropertyTransactionId
	}

	// create loan
	err = l.svcCtx.LoanRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		var innerErr error

		// loan.PropertyId = property.Id
		uid, innerErr := uuid.GetSonyflakeIDString()
		if innerErr != nil {
			return innerErr
		}
		loan.Code = uid
		loan, innerErr = l.svcCtx.LoanRepo.CreateWithTx(l.ctx, tx, loan, false)
		if innerErr != nil {
			return innerErr
		}

		// create loan property
		loanProperty.LoanId = loan.Id
		// loanProperty.PropertyBaseId = propertyBase.Id
		loanProperty, innerErr = l.svcCtx.LoanPropertyRepo.CreateWithTx(l.ctx, tx, loanProperty, false)
		if innerErr != nil {
			return innerErr
		}

		// loan borrower
		loanBorrower.LoanId = loan.Id
		loanBorrower, innerErr = l.svcCtx.LoanBorrowerRepo.CreateWithTx(l.ctx, tx, loanBorrower, false)
		if innerErr != nil {
			return innerErr
		}

		// create applicants
		applicants = slice.Map(applicants, func(index int, item *model.LoanApplicant) *model.LoanApplicant {
			item.LoanId = loan.Id
			item.LoanPropertyId = loanProperty.Id
			item.LoanBorrowerId = loanBorrower.Id
			return item
		})
		_, innerErr = l.svcCtx.LoanApplicantRepo.BatchCreateWithTx(l.ctx, tx, applicants, false)
		if innerErr != nil {
			return innerErr
		}

		// sync total amount
		innerErr = l.svcCtx.BorrowerPropertyRepo.SyncTotalAmountAndLoanCountWithTx(l.ctx, tx, loanProperty.BorrowerPropertyId)
		return innerErr
	})
	if err != nil {
		return nil, errorx.Wrap(err, "fail to create loan")
	}

	return &borrower.LoanCreateResponse{
		Data: &borrower.LoanCreateResponseData{
			Id: loan.Id,
		},
	}, nil
}

// buildCreateModels
func (l *CreateBorrowerLoanLogic) buildCreateModels(in *borrower.LoanCreateRequestData, claim utils.JwtUser) (loan *model.Loan, loanProperty *model.LoanProperty, applicants []*model.LoanApplicant, loanBorrower *model.LoanBorrower, err error) {
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
		applicantItem.Id = 0
		applicantItem.IsPrimary = applicant.IsPrimary
		applicantItem.BorrowerApplicantId = applicant.BorrowerApplicantId
		applicantItem.InvitationStatus = constants.ApplicantInvitationStatusAccepted
		applicantItem.BorrowerOrganizationId = claim.OrganizationId
		applicants = append(applicants, applicantItem)
	}
	if totalOwnPercent > 100 {
		return nil, nil, nil, nil, errorx.New("Total own percent cannot be greater than 100")
	}

	slice.SortBy(applicants, func(a, b *model.LoanApplicant) bool {
		return b.IsPrimary
	})

	primaryCount := slice.CountBy(applicants, func(index int, item *model.LoanApplicant) bool {
		return item.IsPrimary
	})
	if primaryCount == 0 {
		return nil, nil, nil, nil, errorx.New("No primary applicant")
	}
	if primaryCount > 1 {
		return nil, nil, nil, nil, errorx.New("Multiple primary applicants")
	}

	// loan borrower
	loanBorrower = &model.LoanBorrower{}

	// LoanProperty
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

	// loan
	loan = &model.Loan{}
	loan, err = copierx.Trans[borrower.LoanCreateRequestData, model.Loan](in, copierx.StrToSqlTimePtr(timex.LayoutUsStyle))
	if err != nil {
		return nil, nil, nil, nil, err
	}

	loan.CreatorId = claim.OrganizationUserId
	loan.CreatorName = claim.Username
	loan.BorrowerOrganizationId = claim.OrganizationId
	loan.PendingStatus = constants.LoanPendingStatusNone
	loan.Source = constants.LoanSourceBorrower
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
		// loan.PrincipalAmount = in.PrincipalAmount
		if loan.MaturityDate == nil {
			return nil, nil, nil, nil, errorx.New("Maturity date is required")
		}

	case constants.LoanStageArchived:

		// loan.PrincipalAmount = in.PrincipalAmount
		if loan.PaidOffDate == nil {
			return nil, nil, nil, nil, errorx.New("Paid off date is required")
		}

	default:
		return nil, nil, nil, nil, errorx.WithStack(errorx.New("Invalid loan stage"))
	}
	return
}

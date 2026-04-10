package logic

import (
	"context"
	"time"

	"github.com/duke-git/lancet/v2/datetime"
	"github.com/duke-git/lancet/v2/mathutil"
	"gorm.io/gorm"

	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/common/timex"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
)

// LoanLogic provides loan business logic
type LoanLogic struct {
	*BaseLogic
}

// NewLoanLogic creates a new loan logic instance
func NewLoanLogic(baseLogic *BaseLogic) *LoanLogic {
	return &LoanLogic{
		BaseLogic: baseLogic,
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

		operationLogic := NewLoanOperationLogic(l.BaseLogic)
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

// BorrowerLoanSummary get borrower loan summary
func (l *LoanLogic) BorrowerLoanSummary(ctx context.Context, req dto.LoanSummaryRequestDto) (result *dto.LoanSummaryResponseDto, err error) {
	// get loan summary
	result, err = l.repos.LoanRepo.SummaryForBorrower(ctx, req)
	if err != nil {
		return
	}
	// handle holdback
	result.TotalHoldbackRemainingAmount = result.TotalLoanAmount - result.TotalPrincipalAmount - result.TotalReservesAmount
	result.TotalHoldbackReleasedAmount = result.TotalHoldbackAmount - result.TotalHoldbackRemainingAmount

	// get page result
	pageReq := dto.LoanBorrowerPageRequestDto{
		BasePageRequestDto: dto.BasePageRequestDto{
			Page:  1,
			Limit: 5,
		},
		Filter: dto.LoanBorrowerPageRequestFilter{
			OrganizationId: req.OrganizationId,
			Stage:          req.Stage,
			Source:         req.Source,
		},
	}
	pageRes, err := l.repos.LoanRepo.PageForBorrower(ctx, pageReq)
	if err != nil {
		return
	}
	page, err := copierx.Trans[[]*model.Loan, []*dto.LoanEntityDto](&pageRes.List,
		copierx.SqlTimePtrToRFC3339Str(),
		copierx.TimeToRFC3339Str(),
	)
	if err != nil {
		return
	}
	result.List = *page
	result.TotalLoanCount = pageRes.Total

	result.TotalLoanAmount = mathutil.RoundToFloat(result.TotalLoanAmount, 2)
	result.TotalPrincipalAmount = mathutil.RoundToFloat(result.TotalPrincipalAmount, 2)
	result.TotalReservesAmount = mathutil.RoundToFloat(result.TotalReservesAmount, 2)
	result.TotalHoldbackAmount = mathutil.RoundToFloat(result.TotalHoldbackAmount, 2)
	result.TotalHoldbackRemainingAmount = mathutil.RoundToFloat(result.TotalHoldbackRemainingAmount, 2)
	result.TotalHoldbackReleasedAmount = mathutil.RoundToFloat(result.TotalHoldbackReleasedAmount, 2)
	result.AverageRate = mathutil.RoundToFloat(result.AverageRate, 2)
	result.TotalMonthlyPaymentAmount = mathutil.RoundToFloat(result.TotalMonthlyPaymentAmount, 2)
	result.TotalClosingFundAmount = mathutil.RoundToFloat(result.TotalClosingFundAmount, 2)

	l.handleDueInfo(result.List, int(req.DueInDays))

	return
}

func (l *LoanLogic) BorrowerLoanPage(ctx context.Context, req dto.LoanBorrowerPageRequestDto) (result *dto.PageResultDto[*dto.LoanEntityDto], err error) {
	pageResult, err := l.repos.LoanRepo.PageForBorrower(ctx, req)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	result, err = copierx.Trans[gormx.PageResponse[*model.Loan], dto.PageResultDto[*dto.LoanEntityDto]](&pageResult,
		copierx.TimeToRFC3339Str(),
		copierx.SqlTimePtrToStr(timex.LayoutUsStyle),
	)
	if err != nil {
		return nil, err
	}
	l.handleDueInfo(result.List, int(req.Filter.DueInDays))

	return
}

func (l *LoanLogic) handleDueInfo(items []*dto.LoanEntityDto, days int) {
	today := datetime.EndOfDay(time.Now())
	dueDay := datetime.BeginOfDay(time.Now().AddDate(0, 0, days))

	for _, item := range items {
		item.HoldbackRemainingAmount = item.TotalAmount - item.PrincipalAmount - item.ReservesAmount
		if item.Stage != constants.LoanStageActive {
			continue
		}
		if item.MaturityDate == nil || !item.MaturityDate.Valid {
			continue
		}
		maturityDate := item.MaturityDate.Time
		if maturityDate.Before(today) {
			item.DueStatus = constants.LoanDueStatusOverdue
			item.DueDays = int32(datetime.DaysBetween(maturityDate, today))
		} else if maturityDate.Before(dueDay) {
			item.DueStatus = constants.LoanDueStatusDueSoon
			item.DueDays = int32(datetime.DaysBetween(today, dueDay))
		} else {
			item.DueStatus = constants.LoanDueStatusNone
			item.DueDays = int32(datetime.DaysBetween(today, maturityDate))
		}
	}

}

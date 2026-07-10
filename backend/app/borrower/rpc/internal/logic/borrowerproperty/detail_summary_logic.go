package borrowerpropertylogic

import (
	"context"

	"github.com/duke-git/lancet/v2/mathutil"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/timex"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/proto"
)

type DetailSummaryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewDetailSummaryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailSummaryLogic {
	return &DetailSummaryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *DetailSummaryLogic) DetailSummary(in *borrower.PropertyDetailRequest) (resp *borrower.PropertyDetailSummaryResponse, err error) {
	defer func() {
		if err != nil {
			l.svcCtx.Logger.WithContext(l.ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
			).Error(err)
		}
	}()

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	// get property
	borrowerProperty, err := l.svcCtx.BorrowerPropertyRepo.GetOne(l.ctx, in.Data.Id)

	if err != nil {
		return nil, err
	}
	if borrowerProperty == nil || borrowerProperty.Id == 0 {
		return nil, errorx.WithStack(errorx.New("property is invalid"))
	}

	if borrowerProperty.OrganizationId != claim.OrganizationId {
		return nil, errorx.WithStack(errorx.New("access denied"))
	}

	loanProperties := borrowerProperty.LoanProperties

	respData := &borrower.PropertyDetailSummaryResponseData{}
	respData, err = copierx.Trans[model.BorrowerProperty, borrower.PropertyDetailSummaryResponseData](borrowerProperty)
	if err != nil {
		return nil, err
	}

	var totalClosingFundAmount float64
	var totalPrincipalAmount float64
	var totalMonthlyPaymentAmount float64
	var totalRate float32
	var loanCount int

	respData.PendingLoans = make([]*proto.LoanEntity, 0)
	respData.ActiveLoans = make([]*proto.LoanEntity, 0)
	respData.ArchivedLoans = make([]*proto.LoanEntity, 0)

	loanIdMap := make(map[int64]struct{})
	if len(loanProperties) > 0 {
		for i, loanProperty := range loanProperties {
			loan := loanProperty.Loan
			if loan == nil || loan.Id == 0 {
				continue
			}

			_, ok := loanIdMap[loan.Id]
			if ok {
				continue
			}

			loanEntry := &proto.LoanEntity{
				Applicants: []*proto.LoanApplicantEntity{},
			}
			loan.LoanProperty = loanProperties[i]
			err = copierx.Assign[model.Loan, proto.LoanEntity](loan, loanEntry,
				copierx.TimeToRFC3339Str(),
				copierx.SqlTimePtrToStr(timex.LayoutUsStyle),
			)
			if err != nil {
				return nil, err
			}

			err = copierx.Assign[[]*model.LoanApplicant, []*proto.LoanApplicantEntity](&loanProperty.Applicants, &loanEntry.Applicants,
				copierx.TimeToRFC3339Str(),
				copierx.SqlTimePtrToStr(timex.LayoutUsStyle),
			)
			if err != nil {
				return nil, err
			}

			switch loan.Stage {
			case constants.LoanStagePending:
				respData.PendingLoans = append(respData.PendingLoans, loanEntry)
			case constants.LoanStageActive:
				respData.ActiveLoans = append(respData.ActiveLoans, loanEntry)
				totalMonthlyPaymentAmount += loan.MonthlyPayment
				totalPrincipalAmount += loan.PrincipalAmount
				respData.TotalAmount += loan.TotalAmount
				totalRate += loan.Rate

			case constants.LoanStageArchived:
				respData.ArchivedLoans = append(respData.ArchivedLoans, loanEntry)
			}
			totalClosingFundAmount += loan.ClosingFund

			loanCount++
		}
	}
	slice.SortBy(respData.PendingLoans, func(item, other *proto.LoanEntity) bool {
		return item.Id < other.Id
	})
	slice.SortBy(respData.ActiveLoans, func(item, other *proto.LoanEntity) bool {
		return item.Id < other.Id
	})
	slice.SortBy(respData.ArchivedLoans, func(item, other *proto.LoanEntity) bool {
		return item.Id < other.Id
	})

	// remove duplicate loans
	respData.PendingLoans = slice.UniqueBy(respData.PendingLoans, func(item *proto.LoanEntity) int64 {
		return item.Id
	})
	respData.ActiveLoans = slice.UniqueBy(respData.ActiveLoans, func(item *proto.LoanEntity) int64 {
		return item.Id
	})
	respData.ArchivedLoans = slice.UniqueBy(respData.ArchivedLoans, func(item *proto.LoanEntity) int64 {
		return item.Id
	})

	respData.TotalClosingFundAmount = totalClosingFundAmount
	// respData.TotalLoanAmount = property.TotalAmount
	respData.TotalPrincipalAmount = totalPrincipalAmount
	respData.TotalMonthlyPaymentAmount = totalMonthlyPaymentAmount
	if loanCount > 0 {
		respData.AvgRate = totalRate / float32(loanCount)
		respData.AvgRate = float32(mathutil.RoundToFloat(respData.AvgRate, 2))
	}

	return &borrower.PropertyDetailSummaryResponse{
		Data: respData,
	}, nil
}

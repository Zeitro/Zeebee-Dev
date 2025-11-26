package borrowerpropertylogic

import (
	"context"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/jinzhu/copier"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/constants"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/timex"
)

type DetailSummaryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDetailSummaryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailSummaryLogic {
	return &DetailSummaryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DetailSummaryLogic) DetailSummary(in *borrower.PropertyDetailRequest) (resp *borrower.PropertyDetailSummaryResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	// get property
	borrowerProperty, err := l.svcCtx.BorrowerPropertyRepo.GetOne(l.ctx, in.Data.Id)

	if err != nil {
		return nil, err
	}
	if borrowerProperty == nil || borrowerProperty.Id == 0 {
		return nil, errorx.WithStack(errorx.New("Property Is Invalid"))
	}
	// todo check organization
	if borrowerProperty.OrganizationId != in.Data.OrganizationId {
		return nil, errorx.WithStack(errorx.New("Access Denied"))
	}

	loanProperties := borrowerProperty.LoanProperties

	respData := &borrower.PropertyDetailSummaryResponseData{}
	err = copier.CopyWithOption(respData, borrowerProperty, copier.Option{
		DeepCopy:      true,
		CaseSensitive: true,
	})
	if err != nil {
		return nil, errorx.GetFailedTransReqErr(err)
	}

	var totalClosingFundAmount float64
	var totalPrincipalAmount float64
	var totalMonthlyPaymentAmount float64
	var totalRate float32
	var loanCount int

	respData.PendingLoans = make([]*borrower.LoanEntry, 0)
	respData.ActiveLoans = make([]*borrower.LoanEntry, 0)
	respData.ArchivedLoans = make([]*borrower.LoanEntry, 0)

	if len(loanProperties) > 0 {
		for _, loanProperty := range loanProperties {
			loan := loanProperty.Loan
			if loan == nil || loan.Id == 0 {
				continue
			}

			loanEntry := &borrower.LoanEntry{
				Applicants: []*borrower.LoanApplicantEntry{},
			}
			err = copier.CopyWithOption(loanEntry, loan, copier.Option{
				CaseSensitive: true,
				Converters: []copier.TypeConverter{
					copierx.TimeToRFC3339Str(),
					copierx.SqlTimePtrToStr(timex.LayoutUsStyle),
				},
			})

			if err != nil {
				return nil, errorx.GetFailedTransRespErr(err)
			}

			err = copier.CopyWithOption(&loanEntry.Applicants, loanProperty.Applicants, copier.Option{
				CaseSensitive: true,
				Converters: []copier.TypeConverter{
					copierx.TimeToRFC3339Str(),
					copierx.SqlTimePtrToStr(timex.LayoutUsStyle),
				},
			})
			if err != nil {
				return nil, errorx.GetFailedTransRespErr(err)
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
	slice.SortBy(respData.PendingLoans, func(item, other *borrower.LoanEntry) bool {
		return item.Id < other.Id
	})
	slice.SortBy(respData.ActiveLoans, func(item, other *borrower.LoanEntry) bool {
		return item.Id < other.Id
	})
	slice.SortBy(respData.ArchivedLoans, func(item, other *borrower.LoanEntry) bool {
		return item.Id < other.Id
	})

	respData.TotalClosingFundAmount = totalClosingFundAmount
	// respData.TotalLoanAmount = property.TotalAmount
	respData.TotalPrincipalAmount = totalPrincipalAmount
	respData.TotalMonthlyPaymentAmount = totalMonthlyPaymentAmount
	if loanCount > 0 {
		respData.AvgRate = totalRate / float32(loanCount)
	}

	return &borrower.PropertyDetailSummaryResponse{
		Data: respData,
	}, nil
}

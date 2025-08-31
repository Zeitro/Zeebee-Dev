package propertylogic

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/constants"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/timex"
)

type PropertyDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPropertyDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PropertyDetailLogic {
	return &PropertyDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PropertyDetailLogic) PropertyDetail(in *borrower.PropertyDetailRequest) (resp *borrower.PropertyDetailResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	// get property
	propertyBase, err := l.svcCtx.PropertyBaseRepo.GetOne(l.ctx, in.Data.Id)

	if err != nil {
		return nil, err
	}
	if propertyBase == nil || propertyBase.Id == 0 {
		return nil, errors.WithStack(errors.New("property is invalid"))
	}

	if propertyBase.BorrowerId != in.Data.BorrowerId {
		return nil, errors.WithStack(errors.New("access denied"))
	}

	loanProperties := propertyBase.LoanProperties

	respData := &borrower.PropertyDetailResponseData{
		Id: propertyBase.Id,
		// PropertyType: string(loanProperty.PropertyType),
		// TotalAmount:  property.TotalAmount,
		Unit: propertyBase.Unit,
		Address: &borrower.AddressEntry{
			State:    propertyBase.Address.State,
			City:     propertyBase.Address.City,
			Address1: propertyBase.Address.Address1,
			Address2: propertyBase.Address.Address2,
			ZipCode:  propertyBase.Address.ZipCode,
		},
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
				Applicants: []*borrower.ApplicantEntry{},
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

			case constants.LoanStageArchived:
				respData.ArchivedLoans = append(respData.ArchivedLoans, loanEntry)
			}
			totalClosingFundAmount += loan.ClosingFund
			totalRate += loan.Rate
			loanCount++
		}
	}

	respData.TotalClosingFundAmount = totalClosingFundAmount
	// respData.TotalLoanAmount = property.TotalAmount
	respData.TotalPrincipalAmount = totalPrincipalAmount
	respData.TotalMonthlyPaymentAmount = totalMonthlyPaymentAmount
	if loanCount > 0 {
		respData.AvgRate = totalRate / float32(loanCount)
	}

	return &borrower.PropertyDetailResponse{
		Data: respData,
	}, nil
}

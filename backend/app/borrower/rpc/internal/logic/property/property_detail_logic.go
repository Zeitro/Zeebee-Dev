package propertylogic

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"

	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/constants"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/timex"

	"github.com/zeromicro/go-zero/core/logx"
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
	property, err := l.svcCtx.PropertyRepo.GetOne(l.ctx, in.Data.Id)
	if err != nil {
		return nil, err
	}
	if property == nil || property.Id == 0 {
		return nil, errors.WithStack(errors.New("property is invalid"))
	}

	if property.BorrowerId != in.Data.BorrowerId {
		return nil, errors.WithStack(errors.New("access denied"))
	}

	respData := &borrower.PropertyDetailResponseData{
		Id:           property.Id,
		PropertyType: string(property.PropertyType),
		TotalAmount:  property.TotalAmount,
		Unit:         property.Unit,
		Address: &borrower.AddressEntry{
			State:    property.Address.State,
			City:     property.Address.City,
			Address1: property.Address.Address1,
			Address2: property.Address.Address2,
			ZipCode:  property.Address.ZipCode,
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

	for _, loan := range property.Loans {

		loanEntry := &borrower.LoanEntry{}
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

		switch loan.Stage {
		case constants.LoanStagePending:
			respData.PendingLoans = append(respData.PendingLoans, loanEntry)
		case constants.LoanStageActive:
			respData.ActiveLoans = append(respData.ActiveLoans, loanEntry)
			totalMonthlyPaymentAmount += loan.MonthlyPayment
			totalPrincipalAmount += loan.PrincipalAmount
		case constants.LoanStageArchived:
			respData.ArchivedLoans = append(respData.ArchivedLoans, loanEntry)
		}
		totalClosingFundAmount += loan.ClosingFund
		totalRate += loan.Rate
		loanCount++
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

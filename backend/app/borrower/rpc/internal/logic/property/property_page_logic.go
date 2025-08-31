package propertylogic

import (
	"context"
	"database/sql"
	"time"

	"github.com/duke-git/lancet/v2/datetime"
	"github.com/duke-git/lancet/v2/slice"

	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/constants"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type PropertyPageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPropertyPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PropertyPageLogic {
	return &PropertyPageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PropertyPageLogic) PropertyPage(in *borrower.PropertyPageRequest) (resp *borrower.PropertyPageResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	propertyBasePageResult, err := l.svcCtx.PropertyBaseRepo.Page(l.ctx, in.Data)
	if err != nil {
		return nil, err
	}

	properties := make([]*borrower.PropertyPageEntry, len(propertyBasePageResult.List))
	now := datetime.BeginOfDay(time.Now())
	for i, propertyBase := range propertyBasePageResult.List {
		propertyPageEntry := &borrower.PropertyPageEntry{
			Id: propertyBase.Id,
			// PropertyType: string(property.PropertyType),
			TotalAmount: propertyBase.TotalAmount,
			Address: &borrower.AddressEntry{
				State:    propertyBase.Address.State,
				City:     propertyBase.Address.City,
				Address1: propertyBase.Address.Address1,
				Address2: propertyBase.Address.Address2,
				ZipCode:  propertyBase.Address.ZipCode,
			},
		}
		var maturityDate *sql.NullTime
		var totalRate float32 = 0.0
		for _, loan := range propertyBase.Loans {
			totalRate += loan.Rate
			propertyPageEntry.Lenders = append(propertyPageEntry.Lenders, loan.LenderCompany)
			propertyPageEntry.TotalPrincipalAmount += loan.PrincipalAmount
			switch loan.Stage {
			case constants.LoanStagePending:
				propertyPageEntry.PendingLoanCount++
			case constants.LoanStageActive:
				propertyPageEntry.ActiveLoanCount++
				propertyPageEntry.TotalMonthlyPaymentAmount += loan.MonthlyPayment
				// MaturityDate is valid , less than maturityDate,greater than now
				if loan.MaturityDate.Valid && loan.MaturityDate.Time.After(now) && (maturityDate == nil || loan.MaturityDate.Time.Before(maturityDate.Time)) {
					maturityDate = loan.MaturityDate
				}
			case constants.LoanStageArchived:
				propertyPageEntry.ArchivedLoanCount++
			}
		}

		for _, applicant := range propertyBase.Applicants {
			propertyPageEntry.Applicants = append(propertyPageEntry.Applicants, applicant.FullName)
		}
		if totalRate > 0 {
			propertyPageEntry.AvgRate = totalRate / float32(len(propertyBase.Loans))
		}
		if maturityDate != nil && maturityDate.Valid {
			propertyPageEntry.NearestMaturityDate = maturityDate.Time.Format("01/02/2006")
		}
		propertyPageEntry.Applicants = slice.Unique(propertyPageEntry.Applicants)
		propertyPageEntry.Lenders = slice.Unique(propertyPageEntry.Lenders)
		properties[i] = propertyPageEntry

	}

	return &borrower.PropertyPageResponse{
		Data: &borrower.PropertyPageResponseData{
			Page:      propertyBasePageResult.Page,
			Limit:     propertyBasePageResult.Limit,
			Total:     propertyBasePageResult.Total,
			TotalPage: propertyBasePageResult.TotalPage,
			List:      properties,
		},
	}, nil
}

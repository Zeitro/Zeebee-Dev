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
	propertyPageResult, err := l.svcCtx.PropertyRepo.Page(l.ctx, in.Data)
	if err != nil {
		return nil, err
	}

	properties := make([]*borrower.PropertyPageEntry, len(propertyPageResult.List))
	now := datetime.BeginOfDay(time.Now())
	for i, property := range propertyPageResult.List {
		propertyPageEntry := &borrower.PropertyPageEntry{
			Id:           property.Id,
			PropertyType: string(property.PropertyType),
			TotalAmount:  property.TotalAmount,
			Address: &borrower.AddressEntry{
				State:    property.Address.State,
				City:     property.Address.City,
				Address1: property.Address.Address1,
				Address2: property.Address.Address2,
				ZipCode:  property.Address.ZipCode,
			},
		}
		var maturityDate *sql.NullTime
		var totalRate float32 = 0.0
		for _, loan := range property.Loans {
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

		for _, applicant := range property.Applicants {
			propertyPageEntry.Applicants = append(propertyPageEntry.Applicants, applicant.Name)
		}
		if totalRate > 0 {
			propertyPageEntry.AvgRate = totalRate / float32(len(property.Loans))
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
			Page:      propertyPageResult.Page,
			Limit:     propertyPageResult.Limit,
			Total:     propertyPageResult.Total,
			TotalPage: propertyPageResult.TotalPage,
			List:      properties,
		},
	}, nil
}

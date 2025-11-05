package borrowerpropertylogic

import (
	"context"
	"database/sql"
	"time"

	"github.com/duke-git/lancet/v2/datetime"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/jinzhu/copier"

	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/constants"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type PageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PageLogic {
	return &PageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PageLogic) Page(in *borrower.PropertyPageRequest) (resp *borrower.PropertyPageResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	borrowerPropertyPageResult, err := l.svcCtx.BorrowerPropertyRepo.Page(l.ctx, in.Data)
	if err != nil {
		return nil, err
	}

	applicantRepo := l.svcCtx.LoanApplicantRepo
	properties := make([]*borrower.PropertyPageEntry, len(borrowerPropertyPageResult.List))
	now := datetime.BeginOfDay(time.Now())
	for i, borrowerProperty := range borrowerPropertyPageResult.List {
		// propertyPageEntry := &borrower.PropertyPageEntry{
		// 	Id: propertyBase.Id,
		// 	// PropertyType: string(property.PropertyType),
		// 	// TotalAmount: propertyBase.TotalAmount,
		// 	Address: &borrower.AddressEntry{
		// 		State:    propertyBase.Address.State,
		// 		City:     propertyBase.Address.City,
		// 		Address1: propertyBase.Address.Address1,
		// 		Address2: propertyBase.Address.Address2,
		// 		ZipCode:  propertyBase.Address.ZipCode,
		// 	},
		// }

		propertyPageEntry := &borrower.PropertyPageEntry{}
		err = copier.CopyWithOption(propertyPageEntry, borrowerProperty, copier.Option{
			DeepCopy:      true,
			CaseSensitive: true,
		})
		if err != nil {
			return nil, errorx.GetFailedTransReqErr(err)
		}
		var maturityDate *sql.NullTime
		var totalRate float32 = 0.0
		for _, loan := range borrowerProperty.Loans {
			switch loan.Stage {
			case constants.LoanStagePending:
				propertyPageEntry.PendingLoanCount++
			case constants.LoanStageActive:
				propertyPageEntry.TotalPrincipalAmount += loan.PrincipalAmount
				propertyPageEntry.Lenders = append(propertyPageEntry.Lenders, loan.LenderCompany)
				totalRate += loan.Rate
				propertyPageEntry.ActiveLoanCount++
				propertyPageEntry.TotalMonthlyPaymentAmount += loan.MonthlyPayment
				propertyPageEntry.TotalAmount += loan.TotalAmount
				// MaturityDate is valid , less than maturityDate,greater than now
				if (maturityDate == nil || loan.MaturityDate.Time.Before(maturityDate.Time)) && loan.MaturityDate != nil && loan.MaturityDate.Valid && loan.MaturityDate.Time.After(now) {
					maturityDate = loan.MaturityDate
				}
				// get loan applicants
				applicants, err := applicantRepo.GetApplicantsFullNamesByLoanId(l.ctx, loan.Id)
				if err != nil {
					return nil, err
				}
				propertyPageEntry.Applicants = append(propertyPageEntry.Applicants, applicants...)

			case constants.LoanStageArchived:
				propertyPageEntry.ArchivedLoanCount++
			}

		}

		// for _, applicant := range propertyBase.Applicants {
		// 	propertyPageEntry.Applicants = append(propertyPageEntry.Applicants, applicant.FullName)
		// }
		if totalRate > 0 {
			propertyPageEntry.AvgRate = totalRate / float32(len(borrowerProperty.Loans))
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
			Page:      borrowerPropertyPageResult.Page,
			Limit:     borrowerPropertyPageResult.Limit,
			Total:     borrowerPropertyPageResult.Total,
			TotalPage: borrowerPropertyPageResult.TotalPage,
			List:      properties,
		},
	}, nil
}

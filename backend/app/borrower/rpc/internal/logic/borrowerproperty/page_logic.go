package borrowerpropertylogic

import (
	"context"
	"database/sql"
	"time"

	"github.com/duke-git/lancet/v2/datetime"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/jinzhu/copier"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
)

type PageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PageLogic {
	return &PageLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.NewBaseLogic(ctx, svcCtx),
	}
}

func (l *PageLogic) Page(in *borrower.PropertyPageRequest) (resp *borrower.PropertyPageResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	pageReq := dto.BorrowerPropertyPageRequestDto{}
	_ = copier.CopyWithOption(&pageReq, in.Data, copier.Option{
		DeepCopy:      true,
		CaseSensitive: true,
	})
	pageReq.Filter.OrganizationId = claim.OrganizationId
	pageReq.Filter.UserId = claim.OrganizationUserId

	borrowerPropertyPageResult, err := l.svcCtx.BorrowerPropertyRepo.PageForBorrower(l.ctx, pageReq)
	if err != nil {
		return nil, err
	}

	applicantRepo := l.svcCtx.LoanApplicantRepo
	properties := make([]*borrower.PropertyPageEntry, len(borrowerPropertyPageResult.List))
	now := datetime.BeginOfDay(time.Now())
	for i, borrowerProperty := range borrowerPropertyPageResult.List {

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

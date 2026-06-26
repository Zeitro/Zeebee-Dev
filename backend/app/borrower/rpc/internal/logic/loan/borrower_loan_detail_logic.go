package loanlogic

import (
	"context"

	"github.com/duke-git/lancet/v2/slice"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/timex"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type BorrowerLoanDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewBorrowerLoanDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BorrowerLoanDetailLogic {
	return &BorrowerLoanDetailLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *BorrowerLoanDetailLogic) BorrowerLoanDetail(in *borrower.LoanDetailRequest) (resp *borrower.LoanDetailResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
				logx.Field("service_name", l.svcCtx.Config.Name),
			).Error(err)
		}
	}()

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	loan, err := l.svcCtx.LoanRepo.GetBorrowerLoan(l.ctx, in.Data.Id, claim.OrganizationId)
	if err != nil {
		return nil, err
	}
	resp = &borrower.LoanDetailResponse{
		Data: &borrower.LoanEntity{},
	}
	resp.Data, err = copierx.Trans[model.Loan, borrower.LoanEntity](loan,
		copierx.TimeToRFC3339Str(),
		copierx.SqlTimePtrToStr(timex.LayoutUsStyle),
	)
	if err != nil {
		return nil, err
	}

	// loan experiences
	var transExperiences []*borrower.LoanExperienceEntity
	var realEstateExperiences []*borrower.LoanExperienceEntity
	slice.ForEach(loan.LoanExperiences, func(index int, item *model.LoanExperience) {
		experience := &borrower.LoanExperienceEntity{}
		experience, _ = copierx.Trans[model.LoanExperience, borrower.LoanExperienceEntity](item,
			copierx.TimeToRFC3339Str(),
			copierx.TimePtrToRFC3339Str(),
			copierx.SqlTimePtrToStr(timex.LayoutUsStyle),
		)
		if item.Type == constants.LoanExperienceTypeTransaction {
			transExperiences = append(transExperiences, experience)
		} else {
			realEstateExperiences = append(realEstateExperiences, experience)
		}
	})
	// requested borrower loan document count
	slice.ForEach(loan.LoanDocuments, func(index int, item *model.LoanDocument) {
		if item.Category != constants.LoanDocumentCategoryBorrower && item.Category != constants.LoanDocumentCategoryExtend {
			return
		}
		if item.Status == constants.LoanDocumentStatusRequested {
			resp.Data.RequestedLoanDocumentCount++
		}
	})

	resp.Data.TransactionExperiences = transExperiences
	resp.Data.RealEstateExperiences = realEstateExperiences

	return resp, err
}

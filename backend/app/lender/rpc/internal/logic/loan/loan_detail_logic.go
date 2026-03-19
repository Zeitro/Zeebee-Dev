package loanlogic

import (
	"context"

	"github.com/duke-git/lancet/v2/slice"

	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/timex"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type LoanDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewLoanDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoanDetailLogic {
	return &LoanDetailLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *LoanDetailLogic) LoanDetail(in *lender.LoanDetailRequest) (resp *lender.LoanDetailResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	loan, err := l.svcCtx.LoanRepo.GetLenderLoan(l.ctx, in.Data.Id, claim.OrganizationId, true)
	if err != nil {
		return nil, err
	}
	resp = &lender.LoanDetailResponse{
		Data: &lender.LoanEntity{},
	}
	resp.Data, err = copierx.Trans[model.Loan, lender.LoanEntity](loan,
		copierx.TimeToRFC3339Str(),
		copierx.SqlTimePtrToStr(timex.LayoutUsStyle),
	)
	if err != nil {
		return nil, err
	}

	// loan experiences
	var transExperiences []*lender.LoanExperienceEntity
	var realEstateExperiences []*lender.LoanExperienceEntity
	slice.ForEach(loan.LoanExperiences, func(index int, item *model.LoanExperience) {
		experience := &lender.LoanExperienceEntity{}
		experience, _ = copierx.Trans[model.LoanExperience, lender.LoanExperienceEntity](item,
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
	resp.Data.TransactionExperiences = transExperiences
	resp.Data.RealEstateExperiences = realEstateExperiences

	return resp, err
}

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
	"zeebee-backend/domains/zeebee/proto"
)

type LoanDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewLoanDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoanDetailLogic {
	return &LoanDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *LoanDetailLogic) LoanDetail(in *proto.PathIdRequest) (resp *lender.LoanDetailResponse, err error) {

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx, constants.UserTypeLender)
	if err != nil {
		return nil, err
	}

	loan, err := l.svcCtx.LoanRepo.GetLenderLoan(l.ctx, in.Data.Id, claim.OrganizationId, true)
	if err != nil {
		return nil, err
	}
	resp = &lender.LoanDetailResponse{
		Data: &proto.LoanEntity{},
	}
	resp.Data, err = copierx.Trans[model.Loan, proto.LoanEntity](loan,
		copierx.TimeToRFC3339Str(),
		copierx.SqlTimePtrToStr(timex.LayoutUsStyle),
	)
	if err != nil {
		return nil, err
	}

	// loan experiences
	var transExperiences []*proto.LoanExperienceEntity
	var realEstateExperiences []*proto.LoanExperienceEntity
	slice.ForEach(loan.LoanExperiences, func(index int, item *model.LoanExperience) {
		experience := &proto.LoanExperienceEntity{}
		experience, _ = copierx.Trans[model.LoanExperience, proto.LoanExperienceEntity](item,
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

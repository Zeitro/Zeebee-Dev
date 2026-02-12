package loanlogic

import (
	"context"

	"github.com/duke-git/lancet/v2/slice"

	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
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
}

func NewLoanDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoanDetailLogic {
	return &LoanDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoanDetailLogic) LoanDetail(in *admin.IdRequest) (resp *admin.LoanDetailResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	loan, err := l.svcCtx.LoanRepo.Detail(l.ctx, in.Data, true)
	if err != nil {
		return nil, err
	}
	resp = &admin.LoanDetailResponse{
		Data: &admin.LoanEntity{},
	}
	data, err := copierx.Trans[model.Loan, admin.LoanEntity](loan,
		copierx.TimeToRFC3339Str(),
		copierx.SqlTimePtrToStr(timex.LayoutUsStyle),
	)
	if err != nil {
		return nil, err
	}
	// loan experiences
	var transExperiences []*admin.LoanExperienceEntity
	var realEstateExperiences []*admin.LoanExperienceEntity
	slice.ForEach(loan.LoanExperiences, func(index int, item *model.LoanExperience) {
		experience, _ := copierx.Trans[model.LoanExperience, admin.LoanExperienceEntity](item,
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
	data.TransactionExperiences = transExperiences
	data.RealEstateExperiences = realEstateExperiences

	return &admin.LoanDetailResponse{
		Data: data,
	}, nil
}

package loanlogic

import (
	"context"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/jinzhu/copier"

	"zeebee-backend/app/lender/rpc/internal/repo/model"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/constants"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/timex"

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

func (l *LoanDetailLogic) LoanDetail(in *lender.LoanDetailRequest) (resp *lender.LoanDetailResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	loan, err := l.svcCtx.LoanRepo.GetLenderLoan(l.ctx, in.Data.Id, in.Data.LenderCompanyId, true)
	if err != nil {
		return nil, err
	}
	resp = &lender.LoanDetailResponse{
		Data: &lender.LoanEntry{},
	}
	err = copier.CopyWithOption(&resp.Data, loan, copier.Option{
		CaseSensitive: true,
		Converters: []copier.TypeConverter{
			copierx.TimeToRFC3339Str(),
			copierx.SqlTimePtrToStr(timex.LayoutUsStyle),
		},
	})
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	// loan experiences
	var transExperiences []*lender.LoanExperienceEntry
	var realEstateExperiences []*lender.LoanExperienceEntry
	slice.ForEach(loan.LoanExperiences, func(index int, item *model.LoanExperience) {
		experience := &lender.LoanExperienceEntry{}
		_ = copier.CopyWithOption(experience, item, copier.Option{
			CaseSensitive: true,
			DeepCopy:      true,
			IgnoreEmpty:   false,
			Converters: []copier.TypeConverter{
				copierx.TimeToRFC3339Str(),
				copierx.SqlTimePtrToStr(timex.LayoutUsStyle),
			},
		})
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

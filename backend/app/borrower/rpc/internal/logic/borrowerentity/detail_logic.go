package borrowerentitylogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/timex"
	"zeebee-backend/domains/zeebee/model"

	"zeebee-backend/common/errorx"
)

type DetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *DetailLogic) Detail(in *borrower.BorrowerEntityDetailRequest) (resp *borrower.BorrowerEntityDetailResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	respData := &borrower.BorrowerEntityEntity{}
	result, err := l.svcCtx.BorrowerEntityRepo.GetOneByIdForOrganization(l.ctx, in.Data.Id, claim.OrganizationId)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	respData, err = copierx.Trans[model.BorrowerEntity, borrower.BorrowerEntityEntity](result, copierx.TimeToRFC3339Str(),
		copierx.SqlTimePtrToStr(timex.LayoutUsStyle),
	)
	if err != nil {
		return nil, err
	}

	return &borrower.BorrowerEntityDetailResponse{
		Data: respData,
	}, nil
}

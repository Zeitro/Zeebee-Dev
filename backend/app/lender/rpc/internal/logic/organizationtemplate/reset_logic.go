package organizationtemplatelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/data/bo"

	"zeebee-backend/common/errorx"
)

type ResetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewResetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetLogic {
	return &ResetLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.NewBaseLogic(ctx, svcCtx),
	}
}

func (l *ResetLogic) Reset(in *lender.IdRequest) (resp *lender.OrganizationTemplateDetailResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	organizationTemplate, err := l.baseLogic.OrganizationTemplateLogic.Reset(l.ctx, in.Data, claim.OrganizationId)
	if err != nil {
		return nil, err
	}
	result, err := copierx.Trans[bo.OrganizationTemplateBo, lender.OrganizationTemplateEntity](
		organizationTemplate,
		copierx.TimeToRFC3339Str(),
	)

	return &lender.OrganizationTemplateDetailResponse{
		Data: result,
	}, err
}

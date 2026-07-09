package organizationtemplatelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/proto"

	"zeebee-backend/common/errorx"
)

type UpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *UpdateLogic) Update(in *lender.UpdateOrganizationTemplateRequest) (resp *lender.OrganizationTemplateDetailResponse, err error) {
	defer func() {
		if err != nil {
			l.svcCtx.Logger.WithContext(l.ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
			).Error(err)
		}
	}()
	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	input, err := copierx.Trans[lender.UpdateOrganizationTemplateRequestData, dto.OrganizationTemplateUpdateVariableGroupDto](
		in.Data,
	)
	if err != nil {
		return nil, err
	}
	input.OrganizationId = claim.OrganizationId
	result, err := l.baseLogic.OrganizationTemplateLogic.UpdateVariableGroup(l.ctx, *input)
	if err != nil {
		return nil, err
	}
	data, err := copierx.Trans[dto.OrganizationTemplateEntityDto, proto.OrganizationTemplateEntity](
		result,
		copierx.TimeToRFC3339Str(),
	)

	return &lender.OrganizationTemplateDetailResponse{
		Data: data,
	}, err
}

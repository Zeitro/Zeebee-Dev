package organizationtemplatelogic

import (
	"context"

	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"

	"zeebee-backend/common/errorx"
)

type PageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PageLogic {
	return &PageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *PageLogic) Page(in *lender.OrganizationTemplatePageRequest) (resp *lender.OrganizationTemplatePageResponse, err error) {

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx, constants.UserTypeLender)
	if err != nil {
		return nil, errorx.Wrap(err, "fail to Get User Claim From Context")
	}

	input, err := copierx.Trans[lender.OrganizationTemplatePageRequestData, dto.OrganizationTemplatePageRequestDto](
		in.Data,
	)
	if err != nil {
		return nil, err
	}
	input.Filter.OrganizationId = claim.OrganizationId
	input.Filter.Status = constants.CommonStatusActive
	page, err := l.baseLogic.OrganizationTemplateLogic.Page(l.ctx, *input)
	if err != nil {
		return nil, err
	}
	data, err := copierx.Trans[dto.PageResultDto[dto.OrganizationTemplatePageEntityDto], lender.OrganizationTemplatePageResponseData](&page, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}

	return &lender.OrganizationTemplatePageResponse{
		Data: data,
	}, nil
}

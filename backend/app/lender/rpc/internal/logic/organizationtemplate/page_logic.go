package organizationtemplatelogic

import (
	"context"

	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
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
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *PageLogic) Page(in *lender.OrganizationTemplatePageRequest) (resp *lender.OrganizationTemplatePageResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, errorx.Wrap(err, "Failed To Get User Claim From Context")
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

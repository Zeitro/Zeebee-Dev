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

type DetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *DetailLogic) Detail(in *proto.IdRequest) (resp *lender.OrganizationTemplateDetailResponse, err error) {
	defer func() {
		if err != nil {
			l.svcCtx.Logger.WithContext(l.ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
			).Error(err)
		}
	}()
	detail, err := l.baseLogic.OrganizationTemplateLogic.Detail(l.ctx, in.Data)
	if err != nil {
		return nil, err
	}
	data, err := copierx.Trans[dto.OrganizationTemplateEntityDto, proto.OrganizationTemplateEntity](
		detail,
		copierx.TimeToRFC3339Str(),
	)
	if err != nil {
		return nil, err
	}

	return &lender.OrganizationTemplateDetailResponse{
		Data: data,
	}, nil
}

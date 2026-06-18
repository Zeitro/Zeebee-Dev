package basetemplatelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/data/dto"
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

// Page retrieves paginated base template list
func (l *PageLogic) Page(in *admin.BaseTemplatePageRequest) (resp *admin.BaseTemplatePageResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
				logx.Field("service_name", l.svcCtx.Config.Name),
			).Error(err)
		}
	}()

	// Convert proto request to DTO
	pageReq, err := copierx.Trans[admin.BaseTemplatePageRequestData, dto.BaseTemplatePageRequestDto](in.Data)
	if err != nil {
		return nil, err
	}

	// Get paginated results from logic layer
	pageResult, err := l.baseLogic.BaseTemplateLogic.Page(l.ctx, *pageReq)
	if err != nil {
		return nil, err
	}

	// Convert BO page result to proto response
	protoData, err := copierx.Trans[dto.PageResultDto[dto.BaseTemplatePageEntityDto], admin.BaseTemplatePageResponseData](&pageResult, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}

	return &admin.BaseTemplatePageResponse{
		Data: protoData,
	}, nil
}

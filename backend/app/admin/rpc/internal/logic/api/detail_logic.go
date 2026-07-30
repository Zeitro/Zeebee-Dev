package apilogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/proto"
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

func (l *DetailLogic) Detail(in *admin.ApiDetailRequest) (resp *admin.ApiDetailResponse, err error) {

	result, err := l.baseLogic.ApiLogic.Detail(l.ctx, in.Data.Id, constants.UserType(in.Data.Type))
	if err != nil {
		return nil, err
	}
	respData, err := copierx.Trans[dto.ApiEntityDto, proto.ApiEntity](result, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}

	return &admin.ApiDetailResponse{
		Data: respData,
	}, nil
}

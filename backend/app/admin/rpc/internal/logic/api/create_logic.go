package apilogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/proto"
)

type CreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *CreateLogic) Create(in *admin.ApiCreateRequest) (resp *admin.ApiDetailResponse, err error) {

	data := in.Data
	if data == nil {
		return nil, logic.ErrorInvalidRequest
	}
	input, err := copierx.Trans[admin.ApiCreateRequestData, dto.ApiCreateRequestDto](data)
	if err != nil {
		return nil, err
	}

	result, err := l.baseLogic.ApiLogic.CreateApi(l.ctx, input)
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

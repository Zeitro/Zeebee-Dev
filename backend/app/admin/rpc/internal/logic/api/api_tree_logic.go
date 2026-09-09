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

type ApiTreeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewApiTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiTreeLogic {
	return &ApiTreeLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *ApiTreeLogic) ApiTree(in *admin.ApiTreeRequest) (resp *admin.ApiTreeResponse, err error) {

	trees, err := l.baseLogic.ApiLogic.ApiTree(l.ctx, constants.UserType(in.Data.Type))
	if err != nil {
		return nil, err
	}

	data, err := copierx.Trans[[]*dto.ApiEntityDto, []*proto.ApiTreeEntity](&trees, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}

	return &admin.ApiTreeResponse{
		Data: *data,
	}, nil
}

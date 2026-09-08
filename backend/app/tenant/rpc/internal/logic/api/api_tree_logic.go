package apilogic

import (
	"context"

	"zeebee-backend/app/tenant/rpc/internal/logic"
	"zeebee-backend/app/tenant/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
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

// ApiTree returns the api registry tree for the target platform.
func (l *ApiTreeLogic) ApiTree(in *tenant.ApiTreeRequest) (resp *tenant.ApiTreeResponse, err error) {
	if in.Data == nil {
		return nil, errorx.ErrorEmptyRequestData
	}
	callerType, err := logic.ParseUserType(in.Data.UserType)
	if err != nil {
		return nil, err
	}
	targetType, err := targetType(callerType, in.Data.Type)
	if err != nil {
		return nil, err
	}

	trees, err := l.baseLogic.ApiLogic.ApiTree(l.ctx, targetType)
	if err != nil {
		return nil, err
	}

	data, err := copierx.Trans[[]*dto.ApiEntityDto, []*proto.ApiTreeEntity](&trees, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}
	return &tenant.ApiTreeResponse{Data: *data}, nil
}

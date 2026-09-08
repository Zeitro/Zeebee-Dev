package apilogic

import (
	"context"

	"zeebee-backend/app/tenant/rpc/internal/logic"
	"zeebee-backend/app/tenant/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/proto"
)

type DeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

// Delete removes an api endpoint and its children.
func (l *DeleteLogic) Delete(in *tenant.ApiDeleteRequest) (resp *proto.EmptyResponse, err error) {
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

	if err = l.baseLogic.ApiLogic.DeleteApi(l.ctx, in.Data.Id, targetType); err != nil {
		return nil, err
	}
	return &proto.EmptyResponse{}, nil
}

package authzlogic

import (
	"context"

	"zeebee-backend/app/tenant/rpc/internal/svc"
	"zeebee-backend/domains/zeebee/proto"

	"zeebee-backend/common/errorx"
)

type AuthzCheckLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAuthzCheckLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuthzCheckLogic {
	return &AuthzCheckLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AuthzCheckLogic) AuthzCheck(in *proto.AuthCheckRequest) (resp *proto.AuthCheckResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	// todo: add your logic here and delete this line

	return &proto.AuthCheckResponse{}, nil
}

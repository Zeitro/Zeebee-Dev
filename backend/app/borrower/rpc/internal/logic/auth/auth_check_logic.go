package authlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/domains/zeebee/proto"

	"zeebee-backend/common/errorx"
)

type AuthCheckLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAuthCheckLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuthCheckLogic {
	return &AuthCheckLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AuthCheckLogic) AuthCheck(in *proto.AuthCheckRequest) (resp *proto.AuthCheckResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	// todo: add your logic here and delete this line

	return &proto.AuthCheckResponse{}, nil
}

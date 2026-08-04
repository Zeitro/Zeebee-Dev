package userlogic

import (
	"context"

	"zeebee-backend/app/tenant/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"

	"zeebee-backend/common/errorx"
)

type CreateAdminUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateAdminUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAdminUserLogic {
	return &CreateAdminUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateAdminUserLogic) CreateAdminUser(in *tenant.CreateAdminUserRequest) (resp *tenant.CreateAdminUserResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	// todo: add your logic here and delete this line

	return &tenant.CreateAdminUserResponse{}, nil
}

package userlogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"

	"zeebee-backend/common/errorx"
)

type UserSwitchStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewUserSwitchStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserSwitchStatusLogic {
	return &UserSwitchStatusLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.NewBaseLogic(ctx, svcCtx),
	}
}

func (l *UserSwitchStatusLogic) UserSwitchStatus(in *admin.UserSwitchStatusRequest) (resp *admin.EmptyResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	// todo handle switch  status
	// err = l.baseLogic.UserLogic.SwitchOrganizationUserStatus(l.ctx, in.Data.Id, in.Data.OrganizationId, constants.OrganizationAndUserStatus(in.Data.Status), constants.UserTypeAdmin)

	return &admin.EmptyResponse{}, nil
}

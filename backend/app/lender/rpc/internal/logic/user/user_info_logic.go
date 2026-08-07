package userlogic

import (
	"context"

	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	tenantpb "zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/domains/zeebee/proto"
)

type UserInfoLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{ctx: ctx, svcCtx: svcCtx, baseLogic: logic.GetBaseLogic(svcCtx)}
}

// UserInfo proxies to tenant. user_id is resolved from the caller's JWT and the
// lender rpc server interceptor forwards the token so tenant resolves userType
// from the claim.
func (l *UserInfoLogic) UserInfo(in *proto.EmptyRequest) (resp *proto.UserinfoResponse, err error) {
	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	return l.svcCtx.TenantClient.User.UserInfo(l.ctx, &tenantpb.UserInfoRequest{
		Data: &tenantpb.UserInfoRequestData{
			UserId: claim.UserId,
		},
	})
}

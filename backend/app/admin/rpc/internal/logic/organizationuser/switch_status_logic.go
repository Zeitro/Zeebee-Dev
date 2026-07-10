package organizationuserlogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/proto"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type SwitchStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewSwitchStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SwitchStatusLogic {
	return &SwitchStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *SwitchStatusLogic) SwitchStatus(in *admin.OrganizationUserSwitchStatusRequest) (resp *proto.EmptyResponse, err error) {
	defer func() {
		if err != nil {
			l.svcCtx.Logger.WithContext(l.ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
			).Error(err)
		}
	}()
	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	err = l.baseLogic.OrganizationUserLogic.SwitchStatus(l.ctx, in.Data.Id, claim.OrganizationId, constants.OrganizationUserStatus(in.Data.Status), claim.UserType)
	if err != nil {
		return nil, err
	}

	return &proto.EmptyResponse{}, nil
}

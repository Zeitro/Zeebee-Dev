package organizationuserlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/email_tpl"
	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/proto"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type SendInvitationEmailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewSendInvitationEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendInvitationEmailLogic {
	return &SendInvitationEmailLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *SendInvitationEmailLogic) SendInvitationEmail(in *proto.PathIdRequest) (resp *proto.EmptyResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
				logx.Field("service_name", l.svcCtx.Config.Name),
			).Error(err)
		}
	}()
	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	input := dto.OrganizationUserSendInviteEmailRequestDto{
		Id:             in.Data.Id,
		UserType:       claim.UserType,
		OrganizationId: claim.OrganizationId,
		EmailConfig: dto.SendInviteUserEmailConfig{
			HostEmail:         l.svcCtx.Config.EmailConf.Username,
			TplName:           email_tpl.UserInviteEmail,
			ExpireSeconds:     l.svcCtx.Config.InviteUserExpireSeconds,
			Url:               l.svcCtx.Config.InviteUserUrl,
			SendPeriodSeconds: l.svcCtx.Config.SendInviteEmailPeriodSeconds,
		},
		Operator: claim.Username,
	}
	err = l.baseLogic.OrganizationUserLogic.SendInviteEmail(l.ctx, input)
	if err != nil {
		return nil, err
	}

	return &proto.EmptyResponse{}, nil
}

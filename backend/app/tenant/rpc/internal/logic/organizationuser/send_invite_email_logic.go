package organizationuserlogic

import (
	"context"

	"zeebee-backend/app/tenant/rpc/internal/logic"
	"zeebee-backend/app/tenant/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/proto"
)

type SendInviteEmailLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewSendInviteEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendInviteEmailLogic {
	return &SendInviteEmailLogic{ctx: ctx, svcCtx: svcCtx, baseLogic: logic.GetBaseLogic(svcCtx)}
}

// SendInviteEmail (re)sends an invitation email to a pending member.
func (l *SendInviteEmailLogic) SendInviteEmail(in *tenant.OrganizationUserSendInviteEmailRequest) (resp *proto.EmptyResponse, err error) {
	if in.Data == nil {
		return nil, errorx.ErrorEmptyRequestData
	}
	userType, err := logic.ParseUserType(in.Data.UserType)
	if err != nil {
		return
	}

	input := dto.OrganizationUserSendInviteEmailRequestDto{
		Id:             in.Data.Id,
		UserType:       userType,
		OrganizationId: in.Data.OrganizationId,
		EmailConfig:    toInviteEmailConfig(in.Data.EmailConfig),
		Operator:       in.Data.Operator,
	}
	if err = l.baseLogic.OrganizationUserLogic.SendInviteEmail(l.ctx, input); err != nil {
		return
	}
	l.baseLogic.RecordAuditFromClaim(l.ctx, constants.OperationEntityTypeOrganizationUser, in.Data.Id, constants.OperationActionMemberInviteEmailSent)
	return &proto.EmptyResponse{}, nil
}

// toInviteEmailConfig maps the proto invite-email config to the shared dto.
func toInviteEmailConfig(c *tenant.OrganizationUserInviteEmailConfig) dto.SendInviteUserEmailConfig {
	cfg := dto.SendInviteUserEmailConfig{}
	if c == nil {
		return cfg
	}
	cfg.HostEmail = c.HostEmail
	cfg.TplName = c.TplName
	cfg.ExpireSeconds = int(c.ExpireSeconds)
	cfg.SendPeriodSeconds = int(c.SendPeriodSeconds)
	cfg.Url = c.Url
	return cfg
}

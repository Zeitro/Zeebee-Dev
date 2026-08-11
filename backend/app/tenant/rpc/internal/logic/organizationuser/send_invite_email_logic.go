package organizationuserlogic

import (
	"context"

	"zeebee-backend/app/tenant/config"
	"zeebee-backend/app/tenant/rpc/internal/email_tpl"
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
	profile, err := l.baseLogic.PlatformFor(userType)
	if err != nil {
		return
	}

	input := dto.OrganizationUserSendInviteEmailRequestDto{
		Id:             in.Data.Id,
		UserType:       userType,
		OrganizationId: in.Data.OrganizationId,
		EmailConfig:    toInviteEmailConfig(in.Data.EmailConfig, profile),
		Operator:       in.Data.Operator,
	}
	if err = l.baseLogic.OrganizationUserLogic.SendInviteEmail(l.ctx, input); err != nil {
		return
	}
	l.baseLogic.RecordAuditFromClaim(l.ctx, constants.OperationEntityTypeOrganizationUser, in.Data.Id, constants.OperationActionMemberInviteEmailSent)
	return &proto.EmptyResponse{}, nil
}

// toInviteEmailConfig maps the proto invite-email config to the shared dto.
func toInviteEmailConfig(c *tenant.OrganizationUserInviteEmailConfig, profile config.PlatformProfile) dto.SendInviteUserEmailConfig {
	cfg := dto.SendInviteUserEmailConfig{
		HostEmail:         profile.EmailConf.Username,
		ExpireSeconds:     profile.InviteUserExpireSeconds,
		SendPeriodSeconds: profile.SendInviteEmailPeriodSeconds,
		Url:               profile.InviteUserUrl,
		Logo:              profile.LogoUrl,
		Address:           profile.FooterAddress,
	}
	cfg.TplName = email_tpl.UserInviteEmail
	if c == nil {
		return cfg
	}
	if c.HostEmail != "" {
		cfg.HostEmail = c.HostEmail
	}
	if c.Logo != "" {
		cfg.Logo = c.Logo
	}
	if c.Address != "" {
		cfg.Address = c.Address
	}
	if c.TplName != "" {
		cfg.TplName = c.TplName
	}
	if c.ExpireSeconds > 0 {
		cfg.ExpireSeconds = int(c.ExpireSeconds)
	}
	if c.SendPeriodSeconds > 0 {
		cfg.SendPeriodSeconds = int(c.SendPeriodSeconds)
	}
	if c.Url != "" {
		cfg.Url = c.Url
	}
	return cfg
}

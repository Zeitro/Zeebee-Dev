package organizationuserlogic

import (
	"context"

	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/domains/zeebee/constants"
)

type CreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *CreateLogic) Create(in *lender.OrganizationUserCreateRequest) (resp *lender.OrganizationUserDetailResponse, err error) {

	claimUser, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	input := &tenant.OrganizationUserInviteRequestData{
		UserType:       string(constants.UserTypeLender),
		OrganizationId: claimUser.OrganizationId,
		Email:          in.Data.Email,
		Phone:          in.Data.Phone,
		Remark:         in.Data.Remark,
		RoleIds:        in.Data.RoleIds,
		Operator:       claimUser.Username,
		EmailConfig: &tenant.OrganizationUserInviteEmailConfig{
			HostEmail:         l.svcCtx.Config.EmailConf.Username,
			ExpireSeconds:     int32(l.svcCtx.Config.InviteUserExpireSeconds),
			Url:               l.svcCtx.Config.InviteUserUrl,
			SendPeriodSeconds: int32(l.svcCtx.Config.SendInviteEmailPeriodSeconds),
		},
	}
	result, err := l.svcCtx.TenantClient.OrganizationUser.InviteOrganizationUser(l.ctx, &tenant.OrganizationUserInviteRequest{Data: input})
	if err != nil {
		return nil, err
	}

	return &lender.OrganizationUserDetailResponse{
		Data: result.Data,
	}, nil
}

package organizationuserlogic

import (
	"context"

	"zeebee-backend/app/lender/rpc/internal/email_tpl"
	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/proto"

	"zeebee-backend/common/errorx"
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

	claimUser, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx, constants.UserTypeLender)
	if err != nil {
		return nil, err
	}
	input := dto.OrganizationUserCreateRequestDto{
		Email:          in.Data.Email,
		Phone:          in.Data.Phone,
		Remark:         in.Data.Remark,
		RoleIds:        in.Data.RoleIds,
		UserType:       constants.UserTypeLender,
		OrganizationId: claimUser.OrganizationId,
		EmailConfig: dto.SendInviteUserEmailConfig{
			HostEmail:         l.svcCtx.Config.EmailConf.Username,
			TplName:           email_tpl.UserInviteEmail,
			ExpireSeconds:     l.svcCtx.Config.InviteUserExpireSeconds,
			Url:               l.svcCtx.Config.InviteUserUrl,
			SendPeriodSeconds: l.svcCtx.Config.SendInviteEmailPeriodSeconds,
		},
		Operator: claimUser.Username,
	}

	result, err := l.baseLogic.OrganizationUserLogic.InviteOrganizationUser(l.ctx, input)
	if err != nil {
		return nil, err
	}
	data, err := copierx.Trans[dto.OrganizationUserEntityDto, proto.OrganizationUserEntity](result, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	return &lender.OrganizationUserDetailResponse{
		Data: data,
	}, nil
}

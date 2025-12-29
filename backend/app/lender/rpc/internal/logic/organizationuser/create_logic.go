package organizationuserlogic

import (
	"context"

	"github.com/jinzhu/copier"

	"zeebee-backend/app/lender/rpc/internal/email_tpl"
	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type CreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.NewBaseLogic(ctx, svcCtx),
	}
}

func (l *CreateLogic) Create(in *lender.OrganizationUserCreateRequest) (resp *lender.OrganizationUserCreateResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	claimUser, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
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
	data := &lender.OrganizationUserEntity{}
	err = copier.CopyWithOption(data, result, copier.Option{
		CaseSensitive: true,
		Converters: []copier.TypeConverter{
			copierx.TimeToRFC3339Str(),
		},
	})
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	return &lender.OrganizationUserCreateResponse{
		Data: data,
	}, nil
}

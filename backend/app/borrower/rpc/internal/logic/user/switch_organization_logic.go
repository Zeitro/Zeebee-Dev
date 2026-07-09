package userlogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/proto"

	"zeebee-backend/common/errorx"
)

type SwitchOrganizationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewSwitchOrganizationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SwitchOrganizationLogic {
	return &SwitchOrganizationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *SwitchOrganizationLogic) SwitchOrganization(in *proto.SwitchOrganizationRequest) (resp *proto.UserLoginResponse, err error) {
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

	user, err := l.baseLogic.UserLogic.SwitchOrganization(l.ctx, claim.UserId, in.Data.OrganizationId, constants.UserTypeBorrower)
	if err != nil {
		return nil, err
	}

	// token, err := l.baseLogic.UserLogic.GenerateJwtToken(user, organizationUser, constants.UserTypeBorrower)
	token, err := l.baseLogic.UserLogic.GenerateJwtTokenById(l.ctx, user.Id, user.DefaultOrganizationId, constants.UserTypeBorrower)
	if err != nil {
		return nil, err
	}
	data := &proto.UserLoginResponseData{}
	data, err = copierx.Trans[dto.UserJwtTokenDto, proto.UserLoginResponseData](token)
	if err != nil {
		return nil, err
	}
	return &proto.UserLoginResponse{Data: data}, nil
}

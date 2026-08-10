package userlogic

import (
	"context"

	"zeebee-backend/app/tenant/rpc/internal/logic"
	"zeebee-backend/app/tenant/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/proto"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

// Login authenticates a user across the three platforms and issues a TokenSet
// signed with the matching platform's JWT key. user_type selects platform.
func (l *LoginLogic) Login(in *tenant.LoginRequest) (resp *proto.UserLoginResponse, err error) {
	userType, err := logic.ParseUserType(in.Data.UserType)
	if err != nil {
		return nil, err
	}

	loginReq, err := copierx.Trans[tenant.LoginRequestData, dto.UserLoginRequestDto](in.Data)
	if err != nil {
		return nil, err
	}
	loginReq.UserType = userType
	limiter := dto.LoginLimiterDto{
		LimiterItem: l.svcCtx.Config.LimiterConf.Login,
		RedisCli:    l.svcCtx.RedisCli,
	}

	token, err := l.baseLogic.UserLogic.Login(l.ctx, *loginReq, limiter)
	if err != nil {
		l.baseLogic.RecordAudit(l.ctx, dto.OperationCreateDto{
			EntityType:   constants.OperationEntityTypeUser,
			OperatorName: in.Data.Email,
			OperatorType: userType,
			Operation:    constants.OperationActionUserLoginFailed,
		})
		return nil, err
	}
	l.baseLogic.RecordAudit(l.ctx, dto.OperationCreateDto{
		EntityType:   constants.OperationEntityTypeUser,
		OperatorName: in.Data.Email,
		OperatorType: userType,
		Operation:    constants.OperationActionUserLogin,
	})

	data, err := copierx.Trans[dto.UserJwtTokenDto, proto.UserLoginResponseData](token)
	if err != nil {
		return nil, err
	}
	return &proto.UserLoginResponse{Data: data}, nil
}

package authlogic

import (
	"context"

	"zeebee-backend/app/tenant/rpc/internal/logic"
	"zeebee-backend/app/tenant/rpc/internal/svc"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/proto"
)

type AuthCheckLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewAuthCheckLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuthCheckLogic {
	return &AuthCheckLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

// AuthCheck is the central per-request authorization gate (direction B).
// The caller forwards the user's platform JWT in gRPC metadata; the platform is
// given by request.type. Auth resolves the matching secret + policy internally.
func (l *AuthCheckLogic) AuthCheck(in *proto.AuthCheckRequest) (resp *proto.AuthCheckResponse, err error) {
	if in.Data == nil {
		return nil, errorx.New("request data is nil")
	}
	userType, err := logic.ParseUserType(in.Data.Type)
	if err != nil {
		return nil, err
	}

	input, err := copierx.Trans[proto.AuthCheckRequestData, dto.AuthCheckRequestDto](in.Data)
	if err != nil {
		return nil, err
	}
	input.Type = userType

	result, err := l.baseLogic.AuthLogic.AuthCheck(l.ctx, *input)
	if err != nil {
		return nil, err
	}

	data, err := copierx.Trans[dto.AuthCheckResponseDto, proto.AuthCheckResponseData](result)
	if err != nil {
		return nil, err
	}
	return &proto.AuthCheckResponse{Data: data}, nil
}

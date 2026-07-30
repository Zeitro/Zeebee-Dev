package authlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/common/copierx"
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

func (l *AuthCheckLogic) AuthCheck(in *proto.AuthCheckRequest) (resp *proto.AuthCheckResponse, err error) {

	input, err := copierx.Trans[proto.AuthCheckRequestData, dto.AuthCheckRequestDto](in.Data)
	if err != nil {
		return nil, err
	}
	result, err := l.baseLogic.AuthLogic.AuthCheck(l.ctx, *input)
	if err != nil {
		return nil, err
	}
	data, err := copierx.Trans[dto.AuthCheckResponseDto, proto.AuthCheckResponseData](result)
	if err != nil {
		return nil, err
	}

	return &proto.AuthCheckResponse{
		Data: data,
	}, nil
}

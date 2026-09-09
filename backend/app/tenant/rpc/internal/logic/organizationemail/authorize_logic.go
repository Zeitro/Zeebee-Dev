package organizationemaillogic

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

type AuthorizeLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewAuthorizeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuthorizeLogic {
	return &AuthorizeLogic{ctx: ctx, svcCtx: svcCtx, baseLogic: logic.GetBaseLogic(svcCtx)}
}

func (l *AuthorizeLogic) Authorize(in *tenant.OrganizationEmailAuthorizeRequest) (*proto.EmptyResponse, error) {
	if in == nil || in.Data == nil {
		return nil, errorx.ErrorEmptyRequestData
	}
	if in.Data.Code == "" {
		return nil, errorx.New("code is empty")
	}
	userType, err := logic.ParseUserType(in.Data.UserType)
	if err != nil {
		return nil, err
	}
	generator, err := gmailAuthCodeGenerator(l.baseLogic, userType)
	if err != nil {
		return nil, err
	}
	_, _, err = l.baseLogic.OrganizationEmailLogic.EmailAuthorize(l.ctx, dto.OrganizationEmailAuthorizeRequestDto{
		Code:             in.Data.Code,
		State:            in.Data.State,
		EmailProvider:    constants.EmailProvider(in.Data.EmailProvider),
		Purpose:          constants.EmailPurpose(in.Data.State),
		OrganizationId:   in.Data.OrganizationId,
		UserId:           in.Data.UserId,
		OrganizationType: userType,
		Oauth2Config:     generator.GetConfig(),
	})
	if err != nil {
		return nil, err
	}
	return &proto.EmptyResponse{}, nil
}

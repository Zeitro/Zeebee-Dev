package organizationemaillogic

import (
	"context"
	"time"

	"zeebee-backend/app/tenant/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
)

type GetAuthorizationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAuthorizationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAuthorizationLogic {
	return &GetAuthorizationLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *GetAuthorizationLogic) GetAuthorization(in *tenant.OrganizationEmailGetAuthorizationRequest) (*tenant.OrganizationEmailGetAuthorizationResponse, error) {
	if in.Data == nil {
		return nil, errorx.ErrorEmptyRequestData
	}
	organizationEmail, err := l.svcCtx.OrganizationEmailRepo.GetOrganizationEmailByOrganizationId(l.ctx, in.Data.OrganizationId, constants.EmailPurpose(in.Data.Purpose))
	if err != nil {
		return nil, err
	}
	userId := int64(0)
	if organizationEmail.AuthStrategy == constants.EmailAuthStrategyUser {
		userId = in.Data.UserId
	}
	authorization, err := l.svcCtx.EmailAuthorizationRepo.GetEmailAuthorizationByOrganizationEmailId(l.ctx, organizationEmail.Id, userId, organizationEmail.Purpose)
	if err != nil {
		return nil, err
	}
	data := &tenant.OrganizationEmailGetAuthorizationResponseData{
		Email: string(authorization.Email), EmailProvider: string(authorization.EmailProvider), IsAuthorized: authorization.IsAuthorized,
	}
	if authorization.Token != nil && authorization.Token.Token != nil {
		data.AccessToken = authorization.Token.AccessToken
		data.RefreshToken = authorization.Token.RefreshToken
		data.TokenType = authorization.Token.TokenType
		data.Expiry = authorization.Token.Expiry.Format(time.RFC3339)
	}
	return &tenant.OrganizationEmailGetAuthorizationResponse{Data: data}, nil
}

package organizationemaillogic

import (
	"context"

	"zeebee-backend/app/tenant/rpc/internal/logic"
	"zeebee-backend/app/tenant/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/errorx"
)

type GenerateAuthCodeUrlLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewGenerateAuthCodeUrlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenerateAuthCodeUrlLogic {
	return &GenerateAuthCodeUrlLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

// GenerateAuthCodeUrl builds the Gmail OAuth consent URL. Gmail OAuth sending is
// deferred (SMTP-first scope); the tenant wrapper returns an explicit error until
// OAuth client config is wired (see contracts §2.6).
func (l *GenerateAuthCodeUrlLogic) GenerateAuthCodeUrl(in *tenant.OrganizationEmailGenerateAuthCodeUrlRequest) (resp *tenant.OrganizationEmailGenerateAuthCodeUrlResponse, err error) {
	if in.Data == nil {
		return nil, errorx.ErrorEmptyRequestData
	}
	if _, err = logic.ParseUserType(in.Data.UserType); err != nil {
		return
	}

	return nil, errorx.New("gmail oauth not configured in tenant (deferred)")
}

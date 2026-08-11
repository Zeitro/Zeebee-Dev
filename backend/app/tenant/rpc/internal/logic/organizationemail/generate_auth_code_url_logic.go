package organizationemaillogic

import (
	"context"

	"zeebee-backend/app/tenant/rpc/internal/logic"
	"zeebee-backend/app/tenant/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
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

func (l *GenerateAuthCodeUrlLogic) GenerateAuthCodeUrl(in *tenant.OrganizationEmailGenerateAuthCodeUrlRequest) (resp *tenant.OrganizationEmailGenerateAuthCodeUrlResponse, err error) {
	if in.Data == nil {
		return nil, errorx.ErrorEmptyRequestData
	}
	userType, err := logic.ParseUserType(in.Data.UserType)
	if err != nil {
		return
	}
	generator, err := gmailAuthCodeGenerator(l.baseLogic, userType)
	if err != nil {
		return nil, err
	}
	url, err := l.baseLogic.OrganizationEmailLogic.GenerateAuthCodeUrl(l.ctx, dto.OrganizationEmailGenerateAuthCodeUrlRequestDto{
		OrganizationId:    in.Data.OrganizationId,
		UserId:            in.Data.UserId,
		EmailProvider:     constants.EmailProvider(in.Data.EmailProvider),
		Purpose:           constants.EmailPurpose(in.Data.Purpose),
		OrganizationType:  userType,
		AuthCodeGenerator: generator,
	})
	if err != nil {
		return nil, err
	}
	return &tenant.OrganizationEmailGenerateAuthCodeUrlResponse{Data: &tenant.OrganizationEmailGenerateAuthCodeUrlResponseData{Url: url}}, nil
}

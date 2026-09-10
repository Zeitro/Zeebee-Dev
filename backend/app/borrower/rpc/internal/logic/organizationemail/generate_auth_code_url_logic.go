package organizationemaillogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/constants"
)

type GenerateAuthCodeUrlLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewGenerateAuthCodeUrlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenerateAuthCodeUrlLogic {
	return &GenerateAuthCodeUrlLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *GenerateAuthCodeUrlLogic) GenerateAuthCodeUrl(in *borrower.OrganizationEmailGenerateAuthCodeUrlRequest) (resp *borrower.OrganizationEmailGenerateAuthCodeUrlResponse, err error) {

	claim, err := l.baseLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	input, err := copierx.Trans[borrower.OrganizationEmailGenerateAuthCodeUrlRequestData, tenant.OrganizationEmailGenerateAuthCodeUrlRequestData](in.Data)
	if err != nil {
		return nil, err
	}
	input.OrganizationId = claim.OrganizationId
	input.UserId = claim.UserId
	input.UserType = string(constants.UserTypeBorrower)
	result, err := l.svcCtx.TenantClient.OrganizationEmail.GenerateAuthCodeUrl(l.ctx, &tenant.OrganizationEmailGenerateAuthCodeUrlRequest{Data: input})
	if err != nil {
		return nil, err
	}

	return &borrower.OrganizationEmailGenerateAuthCodeUrlResponse{
		Data: &borrower.OrganizationEmailGenerateAuthCodeUrlResponseData{
			Url: result.Data.Url,
		},
	}, nil
}

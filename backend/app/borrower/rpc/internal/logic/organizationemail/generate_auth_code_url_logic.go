package organizationemaillogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/global"
	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type GenerateAuthCodeUrlLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewGenerateAuthCodeUrlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenerateAuthCodeUrlLogic {
	return &GenerateAuthCodeUrlLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *GenerateAuthCodeUrlLogic) GenerateAuthCodeUrl(in *borrower.OrganizationEmailGenerateAuthCodeUrlRequest) (resp *borrower.OrganizationEmailGenerateAuthCodeUrlResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
				logx.Field("service_name", l.svcCtx.Config.Name),
			).Error(err)
		}
	}()
	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	input, err := copierx.Trans[borrower.OrganizationEmailGenerateAuthCodeUrlRequestData, dto.OrganizationEmailGenerateAuthCodeUrlRequestDto](in.Data)
	if err != nil {
		return nil, err
	}
	input.OrganizationId = claim.OrganizationId
	input.UserId = claim.UserId
	input.OrganizationType = constants.UserTypeBorrower

	authCodeUrlGenerator, ok := global.EmailAuthUrlGenerator[input.EmailProvider]
	if !ok {
		return nil, errorx.Errorf("invalid email provider %v", input.EmailProvider)
	}
	input.AuthCodeGenerator = authCodeUrlGenerator

	url, err := l.baseLogic.OrganizationEmailLogic.GenerateAuthCodeUrl(l.ctx, *input)
	if err != nil {
		return nil, err
	}

	return &borrower.OrganizationEmailGenerateAuthCodeUrlResponse{
		Data: &borrower.OrganizationEmailGenerateAuthCodeUrlResponseData{
			Url: url,
		},
	}, nil
}

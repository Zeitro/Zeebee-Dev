package organizationemaillogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/global"
	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/proto"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type AuthorizeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewAuthorizeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuthorizeLogic {
	return &AuthorizeLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *AuthorizeLogic) Authorize(in *borrower.OrganizationEmailAuthorizeRequest) (resp *proto.EmptyResponse, err error) {
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
	input, err := copierx.Trans[borrower.OrganizationEmailAuthorizeRequestData, dto.OrganizationEmailAuthorizeRequestDto](in.Data)
	if err != nil {
		return nil, err
	}
	err = l.validate(input, claim)
	if err != nil {
		return nil, err
	}

	authCodeUrlGenerator, ok := global.EmailAuthUrlGenerator[input.EmailProvider]
	if !ok {
		return nil, errorx.Errorf("invalid email provider %v", input.EmailProvider)
	}
	input.Oauth2Config = authCodeUrlGenerator.GetConfig()
	input.UserId = claim.UserId
	input.OrganizationType = constants.UserTypeBorrower
	input.OrganizationId = claim.OrganizationId
	input.Purpose = constants.EmailPurpose(input.State)

	email, emailClient, err := l.baseLogic.OrganizationEmailLogic.EmailAuthorize(l.ctx, *input)
	if err != nil {
		return nil, err
	}
	global.EmailAuthManager.Register(email, emailClient)
	return &proto.EmptyResponse{}, nil
}

// validate
func (l *AuthorizeLogic) validate(input *dto.OrganizationEmailAuthorizeRequestDto, claim utils.JwtUser) (err error) {
	if input.Code == "" {
		return errorx.New("code is empty")
	}
	return err
}

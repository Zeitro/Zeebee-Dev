package drivelogic

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

type AuthorizeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewAuthorizeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuthorizeLogic {
	return &AuthorizeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *AuthorizeLogic) Authorize(in *borrower.DriveAuthorizeRequest) (resp *borrower.DriveAuthStatusResponse, err error) {
	defer func() {
		if err != nil {
			l.svcCtx.Logger.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	claim, err := l.baseLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	if global.DriveOAuthConfig == nil {
		return nil, errorx.New("Google Drive OAuth is not configured")
	}
	if in.Data == nil {
		return nil, errorx.New("Request data is empty")
	}

	input := dto.DriveAuthorizeRequestDto{
		Code:             in.Data.Code,
		State:            in.Data.State,
		OrganizationId:   claim.OrganizationId,
		UserId:           claim.UserId,
		OrganizationType: constants.UserTypeBorrower,
		Oauth2Config:     global.DriveOAuthConfig,
	}

	result, err := l.baseLogic.DriveLogic.Authorize(l.ctx, input)
	if err != nil {
		return nil, err
	}

	data, err := copierx.Trans[dto.DriveIntegrationDto, borrower.DriveIntegrationEntity](result)
	if err != nil {
		return nil, err
	}

	return &borrower.DriveAuthStatusResponse{Data: data}, nil
}

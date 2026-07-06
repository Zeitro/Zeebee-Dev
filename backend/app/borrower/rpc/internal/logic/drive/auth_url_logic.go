package drivelogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/global"
	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/proto"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type AuthUrlLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewAuthUrlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuthUrlLogic {
	return &AuthUrlLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *AuthUrlLogic) AuthUrl(in *proto.EmptyRequest) (resp *borrower.DriveAuthUrlResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	if global.DriveAuthUrlGenerator == nil {
		return nil, errorx.New("Google Drive OAuth is not configured")
	}

	input := dto.DriveAuthUrlRequestDto{
		OrganizationId:    claim.OrganizationId,
		UserId:            claim.UserId,
		OrganizationType:  constants.UserTypeBorrower,
		AuthCodeGenerator: global.DriveAuthUrlGenerator,
	}

	url, err := l.baseLogic.DriveLogic.GenerateAuthUrl(l.ctx, input)
	if err != nil {
		return nil, err
	}

	return &borrower.DriveAuthUrlResponse{
		Data: &borrower.DriveAuthUrlResponseData{
			Url: url,
		},
	}, nil
}

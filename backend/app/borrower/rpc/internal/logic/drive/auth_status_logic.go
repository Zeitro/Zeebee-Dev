package drivelogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/proto"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type AuthStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewAuthStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuthStatusLogic {
	return &AuthStatusLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *AuthStatusLogic) AuthStatus(in *proto.EmptyRequest) (resp *borrower.DriveAuthStatusResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	result, err := l.baseLogic.DriveLogic.GetIntegration(l.ctx, claim.OrganizationId, constants.UserTypeBorrower)
	if err != nil {
		return nil, err
	}

	data, err := copierx.Trans[dto.DriveIntegrationDto, borrower.DriveIntegrationEntity](result)
	if err != nil {
		return nil, err
	}

	return &borrower.DriveAuthStatusResponse{Data: data}, nil
}

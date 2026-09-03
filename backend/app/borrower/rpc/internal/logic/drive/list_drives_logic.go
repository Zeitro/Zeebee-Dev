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
	"zeebee-backend/domains/zeebee/drive/client"
	"zeebee-backend/domains/zeebee/proto"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type ListDrivesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewListDrivesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListDrivesLogic {
	return &ListDrivesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *ListDrivesLogic) ListDrives(in *proto.EmptyRequest) (resp *borrower.DriveListDrivesResponse, err error) {
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

	input := dto.DriveListDrivesRequestDto{
		OrganizationId:   claim.OrganizationId,
		OrganizationType: constants.UserTypeBorrower,
		Oauth2Config:     global.DriveOAuthConfig,
	}

	drives, err := l.baseLogic.DriveLogic.ListDrives(l.ctx, input)
	if err != nil {
		return nil, err
	}

	data, err := copierx.Trans[[]client.DriveInfo, []*borrower.DriveItemEntity](&drives)
	if err != nil {
		return nil, err
	}

	return &borrower.DriveListDrivesResponse{Data: *data}, nil
}

package drivelogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/proto"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type UpdateFolderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewUpdateFolderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateFolderLogic {
	return &UpdateFolderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *UpdateFolderLogic) UpdateFolder(in *borrower.DriveUpdateFolderRequest) (resp *proto.EmptyResponse, err error) {
	defer func() {
		if err != nil {
			l.svcCtx.Logger.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx, constants.UserTypeBorrower)
	if err != nil {
		return nil, err
	}
	if in.Data == nil {
		return nil, errorx.New("Request data is empty")
	}

	input := dto.DriveUpdateFolderRequestDto{
		OrganizationId:   claim.OrganizationId,
		OrganizationType: constants.UserTypeBorrower,
		FolderName:       in.Data.FolderName,
	}

	err = l.baseLogic.DriveLogic.UpdateFolder(l.ctx, input)
	if err != nil {
		return nil, err
	}

	return &proto.EmptyResponse{}, nil
}

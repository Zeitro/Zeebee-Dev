package drive

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
)

type DriveUpdateFolderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// update folder name
func NewDriveUpdateFolderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DriveUpdateFolderLogic {
	return &DriveUpdateFolderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DriveUpdateFolderLogic) DriveUpdateFolder(req *types.DriveUpdateFolderRequest) (resp *types.EmptyResponse, err error) {
	_, err = l.svcCtx.DriveRpcClient.UpdateFolder(l.ctx, &borrower.DriveUpdateFolderRequest{
		Data: &borrower.DriveUpdateFolderRequestData{
			FolderName: req.FolderName,
		},
	})
	if err != nil {
		return nil, err
	}

	return &types.EmptyResponse{}, nil
}

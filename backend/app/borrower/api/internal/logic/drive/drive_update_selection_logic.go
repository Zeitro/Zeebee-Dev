package drive

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
)

type DriveUpdateSelectionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// update drive selection
func NewDriveUpdateSelectionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DriveUpdateSelectionLogic {
	return &DriveUpdateSelectionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DriveUpdateSelectionLogic) DriveUpdateSelection(req *types.DriveUpdateSelectionRequest) (resp *types.EmptyResponse, err error) {
	_, err = l.svcCtx.DriveRpcClient.UpdateSelection(l.ctx, &borrower.DriveUpdateSelectionRequest{
		Data: &borrower.DriveUpdateSelectionRequestData{
			DriveId:   req.DriveId,
			DriveName: req.DriveName,
		},
	})
	if err != nil {
		return nil, err
	}

	return &types.EmptyResponse{}, nil
}

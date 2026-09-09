package drive

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/proto"
)

type DriveListDrivesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// list shared drives
func NewDriveListDrivesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DriveListDrivesLogic {
	return &DriveListDrivesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DriveListDrivesLogic) DriveListDrives(req *types.EmptyRequest) (resp *types.DriveListDrivesResponse, err error) {
	rpcResp, err := l.svcCtx.DriveRpcClient.ListDrives(l.ctx, &proto.EmptyRequest{})
	if err != nil {
		return nil, err
	}

	lists, err := copierx.Trans[[]*borrower.DriveItemEntity, []types.DriveItemEntity](&rpcResp.Data)
	if err != nil {
		return nil, err
	}

	return &types.DriveListDrivesResponse{
		Lists: *lists,
	}, nil
}

package drive

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/domains/zeebee/proto"
)

type DriveDisconnectLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// disconnect drive
func NewDriveDisconnectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DriveDisconnectLogic {
	return &DriveDisconnectLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DriveDisconnectLogic) DriveDisconnect(req *types.EmptyRequest) (resp *types.EmptyResponse, err error) {
	_, err = l.svcCtx.DriveRpcClient.Disconnect(l.ctx, &proto.EmptyRequest{})
	if err != nil {
		return nil, err
	}

	return &types.EmptyResponse{}, nil
}

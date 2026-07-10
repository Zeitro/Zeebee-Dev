package drive

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/domains/zeebee/proto"
)

type DriveAuthUrlLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// drive auth url
func NewDriveAuthUrlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DriveAuthUrlLogic {
	return &DriveAuthUrlLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DriveAuthUrlLogic) DriveAuthUrl(req *types.EmptyRequest) (resp *types.DriveAuthUrlResponse, err error) {
	rpcResp, err := l.svcCtx.DriveRpcClient.AuthUrl(l.ctx, &proto.EmptyRequest{})
	if err != nil {
		return nil, err
	}

	return &types.DriveAuthUrlResponse{
		Url: rpcResp.Data.Url,
	}, nil
}

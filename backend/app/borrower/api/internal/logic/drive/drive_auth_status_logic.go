package drive

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
)

type DriveAuthStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// drive auth status
func NewDriveAuthStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DriveAuthStatusLogic {
	return &DriveAuthStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DriveAuthStatusLogic) DriveAuthStatus(req *types.EmptyRequest) (resp *types.DriveIntegrationEntity, err error) {
	rpcResp, err := l.svcCtx.DriveRpcClient.AuthStatus(l.ctx, &borrower.EmptyRequest{})
	if err != nil {
		return nil, err
	}

	return copierx.Trans[borrower.DriveIntegrationEntity, types.DriveIntegrationEntity](rpcResp.Data)
}

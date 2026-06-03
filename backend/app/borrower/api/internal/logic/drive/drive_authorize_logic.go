package drive

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
)

type DriveAuthorizeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// drive authorize
func NewDriveAuthorizeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DriveAuthorizeLogic {
	return &DriveAuthorizeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DriveAuthorizeLogic) DriveAuthorize(req *types.DriveAuthorizeRequest) (resp *types.DriveIntegrationEntity, err error) {
	rpcResp, err := l.svcCtx.DriveRpcClient.Authorize(l.ctx, &borrower.DriveAuthorizeRequest{
		Data: &borrower.DriveAuthorizeRequestData{
			Code:  req.Code,
			State: req.State,
		},
	})
	if err != nil {
		return nil, err
	}

	return copierx.Trans[borrower.DriveIntegrationEntity, types.DriveIntegrationEntity](rpcResp.Data)
}

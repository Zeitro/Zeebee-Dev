package drive

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"

	"github.com/zeromicro/go-zero/core/logx"
)

type DriveGetSyncTransactionStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// get latest drive sync job status for a transaction
func NewDriveGetSyncTransactionStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DriveGetSyncTransactionStatusLogic {
	return &DriveGetSyncTransactionStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DriveGetSyncTransactionStatusLogic) DriveGetSyncTransactionStatus(req *types.DriveSyncTransactionRequest) (resp *types.DriveSyncTransactionResponse, err error) {
	rpcResp, err := l.svcCtx.DriveRpcClient.GetSyncTransactionStatus(l.ctx, &borrower.DriveSyncTransactionRequest{
		Data: &borrower.DriveSyncTransactionRequestData{
			TransactionId: req.Id,
		},
	})
	if err != nil {
		return nil, err
	}

	return copierx.Trans[borrower.DriveSyncTransactionResponseData, types.DriveSyncTransactionResponse](rpcResp.Data)
}

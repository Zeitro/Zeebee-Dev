package drive

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
)

type DriveSyncTransactionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// sync transaction documents to drive
func NewDriveSyncTransactionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DriveSyncTransactionLogic {
	return &DriveSyncTransactionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DriveSyncTransactionLogic) DriveSyncTransaction(req *types.DriveSyncTransactionRequest) (resp *types.DriveSyncTransactionResponse, err error) {
	rpcResp, err := l.svcCtx.DriveRpcClient.SyncTransaction(l.ctx, &borrower.DriveSyncTransactionRequest{
		Data: &borrower.DriveSyncTransactionRequestData{
			TransactionId: req.Id,
		},
	})
	if err != nil {
		return nil, err
	}

	return copierx.Trans[borrower.DriveSyncTransactionResponseData, types.DriveSyncTransactionResponse](rpcResp.Data)
}

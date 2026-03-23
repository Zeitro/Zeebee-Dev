package loan

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
)

type SendPreApprovalLoanNotificationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// send pre_approval loan nofification
func NewSendPreApprovalLoanNotificationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendPreApprovalLoanNotificationLogic {
	return &SendPreApprovalLoanNotificationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendPreApprovalLoanNotificationLogic) SendPreApprovalLoanNotification(req *types.SendPreApprovalLoanNotificationRequest) (resp *types.EmptyResponse, err error) {
	rpcReq := &borrower.LoanSendPreApprovalNotificationRequest{
		Data: &borrower.LoanSendPreApprovalNotificationRequestData{},
	}
	rpcReq.Data, err = copierx.Trans[types.SendPreApprovalLoanNotificationRequest, borrower.LoanSendPreApprovalNotificationRequestData](req)
	if err != nil {
		return nil, err
	}

	_, err = l.svcCtx.LoanRpcClient.SendPreApprovalNotification(l.ctx, rpcReq)

	return
}

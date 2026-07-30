package loan_document

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
)

type SendLoanDocumentRequestsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// send loan document requests
func NewSendLoanDocumentRequestsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendLoanDocumentRequestsLogic {
	return &SendLoanDocumentRequestsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendLoanDocumentRequestsLogic) SendLoanDocumentRequests(req *types.LoanDocumentSendRequestsRequest) (resp *types.EmptyResponse, err error) {

	rpcReq := &lender.LoanDocumentSendRequestsRequest{
		Data: &lender.LoanDocumentSendRequestsRequestData{},
	}

	// Copy data from API request to RPC request
	rpcReq.Data, err = copierx.Trans[types.LoanDocumentSendRequestsRequest, lender.LoanDocumentSendRequestsRequestData](req)
	if err != nil {
		return nil, err
	}

	// Call RPC method
	_, err = l.svcCtx.LoanDocumentRpcClient.SendRequests(l.ctx, rpcReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	return &types.EmptyResponse{}, nil
}

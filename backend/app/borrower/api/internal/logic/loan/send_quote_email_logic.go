package loan

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
)

type SendQuoteEmailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// send quote email
func NewSendQuoteEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendQuoteEmailLogic {
	return &SendQuoteEmailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendQuoteEmailLogic) SendQuoteEmail(req *types.LoanSendQuoteEmailRequest) (resp *types.EmptyResponse, err error) {
	rpcReq := &borrower.LoanSendQuoteEmailRequest{}
	rpcReq.Data, err = copierx.Trans[types.LoanSendQuoteEmailRequest, borrower.LoanSendQuoteEmailRequestData](req)
	if err != nil {
		return nil, err
	}
	_, err = l.svcCtx.LoanRpcClient.SendQuoteEmail(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	return &types.EmptyResponse{}, nil
}

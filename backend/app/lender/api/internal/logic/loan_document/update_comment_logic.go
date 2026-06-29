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

type UpdateCommentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// update comment
func NewUpdateCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCommentLogic {
	return &UpdateCommentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateCommentLogic) UpdateComment(req *types.LoanDocumentUpdateCommentRequest) (resp *types.EmptyResponse, err error) {

	rpcReq := &lender.LoanDocumentUpdateCommentRequest{
		Data: &lender.LoanDocumentUpdateCommentRequestData{},
	}

	// Copy data from API request to RPC request
	rpcReq.Data, err = copierx.Trans[types.LoanDocumentUpdateCommentRequest, lender.LoanDocumentUpdateCommentRequestData](req)
	if err != nil {
		return nil, err
	}

	// Call RPC method
	_, err = l.svcCtx.LoanDocumentRpcClient.UpdateComment(l.ctx, rpcReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	return &types.EmptyResponse{}, nil
}

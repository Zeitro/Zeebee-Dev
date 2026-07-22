package borrower_property_transaction

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
)

type ParseFileExtractLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// manually trigger file content extract (manual_tree / manual_select)
func NewParseFileExtractLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ParseFileExtractLogic {
	return &ParseFileExtractLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ParseFileExtract forwards to RPC. A duplicate trigger returns a gRPC
// AlreadyExists error, which the handler maps to HTTP 409.
func (l *ParseFileExtractLogic) ParseFileExtract(req *types.ParseFileExtractRequest) (resp *types.ParseFileExtractResponse, err error) {
	rpcReq := &borrower.ParseFileExtractReq{
		TransactionId: req.Id,
		DocumentId:    req.DocumentId,
		Parser:        req.Parser,
	}
	rpcResp, err := l.svcCtx.BorrowerPropertyTransactionRpcClient.ParseFileExtract(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	return &types.ParseFileExtractResponse{RequestId: rpcResp.RequestId}, nil
}

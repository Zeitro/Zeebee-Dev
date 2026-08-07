package borrower_property_transaction

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
)

type ParseFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// manually trigger file content  (manual_tree / manual_select)
func NewParseFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ParseFileLogic {
	return &ParseFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ParseFile forwards to RPC. A duplicate trigger returns a gRPC
func (l *ParseFileLogic) ParseFile(req *types.ParseFileRequest) (resp *types.ParseFileResponse, err error) {
	rpcReq := &borrower.ParseFileRequest{}

	rpcReq.Data, err = copierx.Trans[types.ParseFileRequest, borrower.ParseFileRequestData](req)
	if err != nil {
		return nil, err
	}

	rpcResp, err := l.svcCtx.BorrowerPropertyTransactionRpcClient.ParseFile(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp, err = copierx.Trans[borrower.ParseFileResponseData, types.ParseFileResponse](rpcResp.Data)
	if err != nil {
		return nil, err
	}
	return
}

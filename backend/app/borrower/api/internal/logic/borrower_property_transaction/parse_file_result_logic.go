package borrower_property_transaction

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
)

type ParseFileResultLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// file extract progress (type=request_id uses request_id; type=params uses document_id+parser)
func NewParseFileResultLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ParseFileResultLogic {
	return &ParseFileResultLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ParseFileResultLogic) ParseFileResult(req *types.ParseFileResultRequest) (resp *types.ParseFileResultResponse, err error) {
	rpcReq := &borrower.ParseFileResultRequest{
		Data: &borrower.ParseFileResultRequestData{
			Type: req.Type,
		},
	}
	switch req.Type {
	case "request_id":
		rpcReq.Data.RequestId = req.RequestId
	case "params":
		rpcReq.Data.Params = &borrower.ProgressByParams{
			TransactionId: req.Id,
			DocumentId:    req.DocumentId,
			Parser:        req.Parser,
		}
	}
	rpcResp, err := l.svcCtx.BorrowerPropertyTransactionRpcClient.ParseFileResult(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	out, err := copierx.Trans[borrower.ParseFileResultResponseData, types.ParseFileResultResponse](rpcResp.Data)
	if err != nil {
		return nil, err
	}
	return out, nil
}

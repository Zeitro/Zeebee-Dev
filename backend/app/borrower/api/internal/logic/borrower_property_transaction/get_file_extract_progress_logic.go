package borrower_property_transaction

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
)

type GetFileExtractProgressLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// file extract progress (type=request_id uses request_id; type=params uses document_id+parser)
func NewGetFileExtractProgressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFileExtractProgressLogic {
	return &GetFileExtractProgressLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFileExtractProgressLogic) GetFileExtractProgress(req *types.GetFileExtractProgressRequest) (resp *types.GetFileExtractProgressResponse, err error) {
	rpcReq := &borrower.GetFileExtractProgressReq{Type: req.Type}
	switch req.Type {
	case "request_id":
		rpcReq.RequestId = req.RequestId
	case "params":
		rpcReq.Params = &borrower.ProgressByParams{
			TransactionId: req.Id,
			DocumentId:    req.DocumentId,
			Parser:        req.Parser,
		}
	}
	rpcResp, err := l.svcCtx.BorrowerPropertyTransactionRpcClient.GetFileExtractProgress(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	out, err := copierx.Trans[borrower.GetFileExtractProgressResp, types.GetFileExtractProgressResponse](rpcResp)
	if err != nil {
		return nil, err
	}
	return out, nil
}

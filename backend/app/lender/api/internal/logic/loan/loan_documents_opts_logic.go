package loan

import (
	"context"

	"github.com/jinzhu/copier"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/errorx"
)

type LoanDocumentsOptsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// loan documents opts
func NewLoanDocumentsOptsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoanDocumentsOptsLogic {
	return &LoanDocumentsOptsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoanDocumentsOptsLogic) LoanDocumentsOpts(req *types.PathIdRequest) (resp *types.LoanDocumentOptResponse, err error) {
	// Call RPC method
	rpcResp, err := l.svcCtx.LoanRpcClient.LoanDocumentOpts(l.ctx, &lender.LoanDocumentOptsRequest{
		Data: &lender.LoanDocumentOptsRequestData{
			Id: req.Id,
		},
	})
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	// Convert RPC response to API response
	resp = &types.LoanDocumentOptResponse{}
	err = copier.CopyWithOption(&resp.Options, rpcResp.Data, copier.Option{
		IgnoreEmpty:   false,
		DeepCopy:      true,
		CaseSensitive: true,
	})
	if err != nil {
		return nil, errorx.GetFailedTransErr(err)
	}

	return resp, nil
}

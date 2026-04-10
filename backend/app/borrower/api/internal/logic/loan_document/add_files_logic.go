package loan_document

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
)

type AddFilesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// add files
func NewAddFilesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddFilesLogic {
	return &AddFilesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddFilesLogic) AddFiles(req *types.LoanDocumentAddFilesRequest) (resp *types.EmptyResponse, err error) {
	rpcReq := &borrower.LoanDocumentAddFilesRequest{
		Data: &borrower.LoanDocumentAddFilesRequestData{},
	}

	// Copy data from API request to RPC request
	rpcReq.Data, err = copierx.Trans[types.LoanDocumentAddFilesRequest, borrower.LoanDocumentAddFilesRequestData](req)
	if err != nil {
		return nil, err
	}

	// Call RPC method
	_, err = l.svcCtx.LoanDocumentRpcClient.AddFiles(l.ctx, rpcReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	return &types.EmptyResponse{}, nil
}

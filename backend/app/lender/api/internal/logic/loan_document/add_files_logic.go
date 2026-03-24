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
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	rpcReq := &lender.LoanDocumentAddFilesRequest{
		Data: &lender.LoanDocumentAddFilesRequestData{},
	}

	// Copy data from API request to RPC request
	rpcReq.Data, err = copierx.Trans[types.LoanDocumentAddFilesRequest, lender.LoanDocumentAddFilesRequestData](req)
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

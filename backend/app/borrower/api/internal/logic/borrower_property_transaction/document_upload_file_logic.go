package borrower_property_transaction

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/proto"
)

type DocumentUploadFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDocumentUploadFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DocumentUploadFileLogic {
	return &DocumentUploadFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DocumentUploadFileLogic) DocumentUploadFile(req *types.BorrowerPropertyTransactionDocumentUploadFileRequest) (resp *types.DocumentUploadFileResponse, err error) {
	rpcReq := &borrower.BorrowerPropertyTransactionDocumentUploadFileRequest{}

	rpcReq.Data, err = copierx.Trans[types.BorrowerPropertyTransactionDocumentUploadFileRequest, borrower.BorrowerPropertyTransactionDocumentUploadFileRequestData](req)
	if err != nil {
		return nil, err
	}

	rpcResp, err := l.svcCtx.BorrowerPropertyTransactionRpcClient.UploadDocumentFile(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	files, err := copierx.Trans[[]*proto.FileEntity, []types.DocumentFile](&rpcResp.Data)
	if err != nil {
		return nil, err
	}

	return &types.DocumentUploadFileResponse{Files: *files, RequestId: rpcResp.RequestId}, nil
}

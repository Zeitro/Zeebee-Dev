package borrower_property_transaction

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
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
	rpcReq := &borrower.BorrowerPropertyTransactionDocumentUploadFileRequest{
		// Data: &borrower.BorrowerPropertyTransactionDocumentUploadFileRequestData{
		// 	TransactionId: req.TransactionId,
		// 	DocumentId:    req.DocumentId,
		// 	FileName:      req.FileName,
		// 	FileUrl:       req.FileUrl,
		// 	StorageType:   req.StorageType,
		// 	UploadType:    req.UploadType,
		// },
	}

	rpcReq.Data, err = copierx.Trans[types.BorrowerPropertyTransactionDocumentUploadFileRequest, borrower.BorrowerPropertyTransactionDocumentUploadFileRequestData](req)
	if err != nil {
		return nil, err
	}

	rpcResp, err := l.svcCtx.BorrowerPropertyTransactionRpcClient.UploadDocumentFile(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp, err = copierx.Trans[borrower.FileEntity, types.DocumentUploadFileResponse](rpcResp.Data)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

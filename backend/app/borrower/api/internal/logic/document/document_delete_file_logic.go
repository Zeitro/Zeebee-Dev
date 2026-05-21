package document

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
)

type DocumentDeleteFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDocumentDeleteFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DocumentDeleteFileLogic {
	return &DocumentDeleteFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DocumentDeleteFileLogic) DocumentDeleteFile(req *types.DocumentDeleteFileRequest) (resp *types.EmptyResponse, err error) {
	rpcReq := &borrower.DocumentDeleteFileRequest{
		Data: &borrower.DocumentDeleteFileRequestData{
			FileId: req.FileId,
		},
	}

	_, err = l.svcCtx.DocumentRpcClient.DeleteFile(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	return &types.EmptyResponse{}, nil
}

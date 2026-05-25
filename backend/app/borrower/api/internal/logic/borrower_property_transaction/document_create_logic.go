package borrower_property_transaction

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
)

type DocumentCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDocumentCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DocumentCreateLogic {
	return &DocumentCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DocumentCreateLogic) DocumentCreate(req *types.DocumentCreateRequest) (resp *types.DocumentCreateResponse, err error) {
	rpcReq := &borrower.DocumentCreateRequest{
		Data: &borrower.DocumentCreateRequestData{
			TransactionId: req.TransactionId,
			ParentId:      req.ParentId,
			Name:          req.Name,
		},
	}

	rpcResp, err := l.svcCtx.BorrowerPropertyTransactionRpcClient.CreateDocument(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	protoDoc := rpcResp.GetData()
	if protoDoc == nil {
		return &types.DocumentCreateResponse{}, nil
	}

	treeNode, err := copierx.Trans[borrower.DocumentEntity, types.DocumentTreeNode](protoDoc, copierx.StrToStrMap())
	if err != nil {
		return nil, err
	}

	return &types.DocumentCreateResponse{*treeNode}, nil
}

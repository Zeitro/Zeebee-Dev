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
		Data: &borrower.DocumentCreateRequestData{},
	}

	rpcReq.Data, err = copierx.Trans[types.DocumentCreateRequest, borrower.DocumentCreateRequestData](req)
	if err != nil {
		return nil, err
	}

	rpcResp, err := l.svcCtx.BorrowerPropertyTransactionRpcClient.CreateDocument(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	protoDocs := rpcResp.GetData()
	if len(protoDocs) == 0 {
		return &types.DocumentCreateResponse{}, nil
	}

	var treeNodes []types.DocumentTreeNode
	for _, protoDoc := range protoDocs {
		treeNode, transErr := copierx.Trans[borrower.DocumentEntity, types.DocumentTreeNode](protoDoc, copierx.StrToStrMap())
		if transErr != nil {
			return nil, transErr
		}
		treeNodes = append(treeNodes, *treeNode)
	}

	return &types.DocumentCreateResponse{List: treeNodes}, nil
}

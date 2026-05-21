package document

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
)

type DocumentTreeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDocumentTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DocumentTreeLogic {
	return &DocumentTreeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DocumentTreeLogic) DocumentTree(req *types.DocumentTreeRequest) (resp *types.DocumentTreeResponse, err error) {
	rpcReq := &borrower.DocumentDetailRequest{
		Data: &borrower.DocumentDetailRequestData{
			EntityType: req.EntityType,
			EntityId:   req.EntityId,
		},
	}

	rpcResp, err := l.svcCtx.DocumentRpcClient.Detail(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	protoNode := rpcResp.GetData()
	if protoNode == nil {
		return &types.DocumentTreeResponse{}, nil
	}

	treeNode, err := convertDocumentNode(protoNode)
	if err != nil {
		return nil, err
	}

	return &types.DocumentTreeResponse{
		Data: *treeNode,
	}, nil
}

func convertDocumentNode(protoDoc *borrower.DocumentEntity) (*types.DocumentTreeNode, error) {
	node, err := copierx.Trans[borrower.DocumentEntity, types.DocumentTreeNode](protoDoc)
	if err != nil {
		return nil, err
	}

	for _, child := range protoDoc.GetChildren() {
		childNode, err := convertDocumentNode(child)
		if err != nil {
			return nil, err
		}
		node.Children = append(node.Children, *childNode)
	}

	for _, f := range protoDoc.GetFiles() {
		fileEntity, err := copierx.Trans[borrower.FileEntity, types.DocumentFile](f)
		if err != nil {
			return nil, err
		}
		node.Files = append(node.Files, *fileEntity)
	}

	return node, nil
}

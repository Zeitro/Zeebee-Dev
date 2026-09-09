package loan_document

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/proto"
)

type DeleteLoanDocumentRequestLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// delete loan document request
func NewDeleteLoanDocumentRequestLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLoanDocumentRequestLogic {
	return &DeleteLoanDocumentRequestLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteLoanDocumentRequestLogic) DeleteLoanDocumentRequest(req *types.PathIdRequest) (resp *types.EmptyResponse, err error) {

	// Call RPC method
	_, err = l.svcCtx.LoanDocumentRpcClient.DeleteRequest(l.ctx, &proto.PathIdRequest{
		Data: &proto.PathIdRequestData{
			Id: req.Id,
		},
	})
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	return &types.EmptyResponse{}, nil
}

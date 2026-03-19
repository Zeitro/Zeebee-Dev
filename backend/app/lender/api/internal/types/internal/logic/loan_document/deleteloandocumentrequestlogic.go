// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package loan_document

import (
	"context"

	"zeebee-backend/app/lender/api/internal/types/internal/svc"
	"zeebee-backend/app/lender/api/internal/types/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
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
	// todo: add your logic here and delete this line

	return
}

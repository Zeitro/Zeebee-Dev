// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package loan_document

import (
	"context"

	"zeebee-backend/app/lender/api/internal/types/internal/svc"
	"zeebee-backend/app/lender/api/internal/types/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLoanDocumentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// create loan document
func NewCreateLoanDocumentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLoanDocumentLogic {
	return &CreateLoanDocumentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLoanDocumentLogic) CreateLoanDocument(req *types.LoanDocumentCreateRequest) (resp *types.LoanDocumentCreateResponse, err error) {
	// todo: add your logic here and delete this line

	return
}

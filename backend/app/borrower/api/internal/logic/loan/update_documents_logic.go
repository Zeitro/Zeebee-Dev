package loan

import (
	"context"

	"github.com/jinzhu/copier"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateDocumentsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// update documents
func NewUpdateDocumentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateDocumentsLogic {
	return &UpdateDocumentsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateDocumentsLogic) UpdateDocuments(req *types.UpdateDocumentsRequest) (resp *types.EmptyResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	borrowerId, err := utils.GetAccountIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	rpcReq := &borrower.LoanUpdateDocumentsRequest{
		Data: &borrower.LoanUpdateDocumentsRequestData{
			BorrowerId: borrowerId,
		},
	}

	err = copier.CopyWithOption(&rpcReq.Data, req, copier.Option{
		IgnoreEmpty:   true,
		CaseSensitive: true,
		DeepCopy:      true,
	})
	if err != nil {
		return nil, errorx.GetFailedTransErr(err)
	}

	_, err = l.svcCtx.LoanRpcClient.UpdateLoanDocuments(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	return &types.EmptyResponse{}, nil

}

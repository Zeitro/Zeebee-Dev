package loan_document

import (
	"context"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/utils"
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
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	// Get lender company id and employee id from context
	lenderCompanyId, lenderEmployeeId, err := utils.GetLenderCompanyAndEmployeeId(l.ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// Call RPC method
	_, err = l.svcCtx.LoanDocumentRpcClient.DeleteRequest(l.ctx, &lender.LoanDocumentDeleteRequestRequest{
		Data: &lender.LoanDocumentDeleteRequestRequestData{
			Id:               req.Id,
			LenderCompanyId:  lenderCompanyId,
			LenderEmployeeId: lenderEmployeeId,
		},
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &types.EmptyResponse{}, nil
}

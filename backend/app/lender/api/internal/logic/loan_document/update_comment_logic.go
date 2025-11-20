package loan_document

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/utils"
)

type UpdateCommentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// update comment
func NewUpdateCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCommentLogic {
	return &UpdateCommentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateCommentLogic) UpdateComment(req *types.LoanDocumentUpdateCommentRequest) (resp *types.EmptyResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	rpcReq := &lender.LoanDocumentUpdateCommentRequest{
		Data: &lender.LoanDocumentUpdateCommentRequestData{},
	}

	// Copy data from API request to RPC request
	err = copier.CopyWithOption(&rpcReq.Data, req, copier.Option{
		IgnoreEmpty:   false,
		DeepCopy:      true,
		CaseSensitive: true,
	})
	if err != nil {
		return nil, errorx.GetFailedTransReqErr(err)
	}

	// Get lender company id and employee id from context
	rpcReq.Data.LenderCompanyId, rpcReq.Data.LenderEmployeeId, err = utils.GetLenderCompanyAndEmployeeId(l.ctx)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	// Call RPC method
	_, err = l.svcCtx.LoanDocumentRpcClient.UpdateComment(l.ctx, rpcReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	return &types.EmptyResponse{}, nil
}

package loan_document

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/utils"
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

func (l *CreateLoanDocumentLogic) CreateLoanDocument(req *types.LoanDocumentCreateRequest) (resp *types.EmptyResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	rpcReq := &lender.LoanDocumentCreateExtendTypeRequest{
		Data: &lender.LoanDocumentCreateExtendTypeRequestData{},
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
		return nil, errors.WithStack(err)
	}

	// Call RPC method
	_, err = l.svcCtx.LoanDocumentRpcClient.Create(l.ctx, rpcReq)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &types.EmptyResponse{}, nil
}

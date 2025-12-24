package loan_document

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/errorx"
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
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	resp = &types.LoanDocumentCreateResponse{
		LoanDocuments: []*types.LoanDocumentOptEntity{},
	}
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

	// Call RPC method
	rpcResp, err := l.svcCtx.LoanDocumentRpcClient.Create(l.ctx, rpcReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	err = copier.CopyWithOption(&resp.LoanDocuments, rpcResp.Data, copier.Option{
		IgnoreEmpty:   false,
		DeepCopy:      true,
		CaseSensitive: true,
	})
	if err != nil {
		return nil, errorx.GetFailedTransErr(err)
	}

	return resp, nil
}

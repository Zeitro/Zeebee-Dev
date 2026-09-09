package loan_document

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/proto"
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

	resp = &types.LoanDocumentCreateResponse{
		LoanDocuments: []*types.LoanDocumentOptEntity{},
	}
	rpcReq := &lender.LoanDocumentCreateExtendTypeRequest{
		Data: &lender.LoanDocumentCreateExtendTypeRequestData{},
	}

	// Copy data from API request to RPC request
	rpcReq.Data, err = copierx.Trans[types.LoanDocumentCreateRequest, lender.LoanDocumentCreateExtendTypeRequestData](req)
	if err != nil {
		return nil, err
	}

	// Call RPC method
	rpcResp, err := l.svcCtx.LoanDocumentRpcClient.Create(l.ctx, rpcReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	documents, err := copierx.Trans[[]*proto.LoanDocumentOptEntity, []*types.LoanDocumentOptEntity](&rpcResp.Data)
	if err != nil {
		return nil, err
	}
	resp.LoanDocuments = *documents

	return resp, nil
}

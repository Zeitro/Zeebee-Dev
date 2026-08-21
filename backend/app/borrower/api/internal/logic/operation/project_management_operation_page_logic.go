package operation

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProjectManagementOperationPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// page project operation timeline
func NewProjectManagementOperationPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectManagementOperationPageLogic {
	return &ProjectManagementOperationPageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProjectManagementOperationPageLogic) ProjectManagementOperationPage(req *types.OperationPageRequest) (resp *types.OperationPageResponse, err error) {
	data, err := copierx.Trans[types.OperationPageRequest, borrower.OperationPageRequestData](req)
	if err != nil {
		return nil, err
	}
	rpcResp, err := l.svcCtx.OperationRpcClient.Page(l.ctx, &borrower.OperationPageRequest{Data: data})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	return copierx.Trans[borrower.OperationPageResponseData, types.OperationPageResponse](rpcResp.Data)
}

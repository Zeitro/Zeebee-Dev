package organization_email

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrganizationEmailAuthorizeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// email authorize
func NewOrganizationEmailAuthorizeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationEmailAuthorizeLogic {
	return &OrganizationEmailAuthorizeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationEmailAuthorizeLogic) OrganizationEmailAuthorize(req *types.OrganizationEmailAuthorizeRequest) (resp *types.EmptyResponse, err error) {
	rpcReq := &borrower.OrganizationEmailAuthorizeRequest{}

	rpcReq.Data, err = copierx.Trans[types.OrganizationEmailAuthorizeRequest, borrower.OrganizationEmailAuthorizeRequestData](req)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	_, err = l.svcCtx.OrganizationEmailRpcClient.Authorize(l.ctx, rpcReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	return &types.EmptyResponse{}, nil
}

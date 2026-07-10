package organization

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrganizationLenderPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// lender page
func NewOrganizationLenderPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationLenderPageLogic {
	return &OrganizationLenderPageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationLenderPageLogic) OrganizationLenderPage(req *types.OrganizationLenderPageRequest) (resp *types.OrganizationLenderPageResponse, err error) {
	rpcReq := &proto.CommonPageRequest{}

	rpcReq.Data, err = copierx.Trans[types.OrganizationLenderPageRequest, proto.CommonPageRequestData](req)
	if err != nil {
		return nil, err
	}

	rpcResp, err := l.svcCtx.OrganizationRpcClient.LenderPage(l.ctx, rpcReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	resp, err = copierx.Trans[borrower.OrganizationLenderPageResponseData, types.OrganizationLenderPageResponse](rpcResp.Data)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

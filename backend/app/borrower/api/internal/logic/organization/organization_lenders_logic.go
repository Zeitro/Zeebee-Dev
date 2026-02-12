// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package organization

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrganizationLendersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// lenders
func NewOrganizationLendersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationLendersLogic {
	return &OrganizationLendersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationLendersLogic) OrganizationLenders(req *types.EmptyRequest) (resp *types.OrganizationLendersResponse, err error) {
	rpcReq := &borrower.EmptyRequest{}
	resp = &types.OrganizationLendersResponse{}

	rpcResp, err := l.svcCtx.OrganizationRpcClient.AllLender(l.ctx, rpcReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	list, err := copierx.Trans[[]*borrower.OrganizationOptionsEntity, []*types.OrganizationOptionEntity](&rpcResp.Data)
	if err != nil {
		return nil, err
	}
	resp.List = *list

	return resp, nil
}

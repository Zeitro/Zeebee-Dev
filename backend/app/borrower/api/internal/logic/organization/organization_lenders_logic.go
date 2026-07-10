// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package organization

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/proto"

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
	rpcReq := &proto.EmptyRequest{}
	resp = &types.OrganizationLendersResponse{}

	rpcResp, err := l.svcCtx.OrganizationRpcClient.AllLender(l.ctx, rpcReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	list, err := copierx.Trans[[]*proto.OrganizationOptionsEntity, []*types.OrganizationOptionEntity](&rpcResp.Data)
	if err != nil {
		return nil, err
	}
	resp.List = *list

	return resp, nil
}

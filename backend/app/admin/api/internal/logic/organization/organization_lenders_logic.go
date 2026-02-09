package organization

import (
	"context"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/app/admin/rpc/client/organization"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrganizationLendersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// all lender
func NewOrganizationLendersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationLendersLogic {
	return &OrganizationLendersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationLendersLogic) OrganizationLenders(req *types.EmptyRequest) (resp *types.OrganizationOptionResponse, err error) {
	rpcReq := &organization.EmptyRequest{}
	resp = &types.OrganizationOptionResponse{}

	rpcResp, err := l.svcCtx.OrganizationRpcClient.AllLender(l.ctx, rpcReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	list, err := copierx.Trans[[]*admin.OrganizationOptionsEntity, []*types.OrganizationOptionEntity](&rpcResp.Data)
	if err != nil {
		return nil, err
	}
	resp.List = *list

	return resp, nil
}

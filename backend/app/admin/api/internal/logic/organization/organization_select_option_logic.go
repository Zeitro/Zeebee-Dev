package organization

import (
	"context"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrganizationSelectOptionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// select option
func NewOrganizationSelectOptionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationSelectOptionLogic {
	return &OrganizationSelectOptionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationSelectOptionLogic) OrganizationSelectOption(req *types.OrganizationSelectOptionRequest) (resp *types.OrganizationOptResponse, err error) {
	rpcReq := &admin.OrganizationSelectOptionsRequest{}
	rpcReq.Data, err = copierx.Trans[types.OrganizationSelectOptionRequest, admin.OrganizationSelectOptionsRequestData](req)

	if err != nil {
		return nil, err
	}

	rpcResp, err := l.svcCtx.OrganizationRpcClient.SelectOptions(l.ctx, rpcReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	resp = &types.OrganizationOptResponse{}
	options, err := copierx.Trans[[]*admin.OrganizationSelectOptionEntity, []*types.OrganizationSelectOptionEntity](&rpcResp.Data)
	if err != nil {
		return nil, err
	}
	resp.Options = *options

	return resp, nil
}

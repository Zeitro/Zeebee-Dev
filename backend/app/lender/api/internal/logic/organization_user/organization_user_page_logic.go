package organization_user

import (
	"context"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrganizationUserPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// page
func NewOrganizationUserPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationUserPageLogic {
	return &OrganizationUserPageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationUserPageLogic) OrganizationUserPage(req *types.OrganizationUserPageRequest) (resp *types.OrganizationUserPageResponse, err error) {
	data, err := copierx.Trans[types.OrganizationUserPageRequest, lender.OrganizationUserPageRequestData](req)
	if err != nil {
		return nil, err
	}

	rpcReq := &lender.OrganizationUserPageRequest{
		Data: data,
	}

	rpcResp, err := l.svcCtx.OrganizationUserRpcClient.Page(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	var list []*types.OrganizationUserPageEntity
	for _, v := range rpcResp.Data.List {
		item, err := copierx.Trans[lender.OrganizationUserPageEntity, types.OrganizationUserPageEntity](v, copierx.TimeToRFC3339Str())
		if err != nil {
			return nil, err
		}
		list = append(list, item)
	}

	resp = &types.OrganizationUserPageResponse{
		Page:      rpcResp.Data.Page,
		Limit:     rpcResp.Data.Limit,
		Total:     rpcResp.Data.Total,
		TotalPage: rpcResp.Data.TotalPage,
		List:      list,
	}

	return
}

package organization_user

import (
	"context"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrganizationUserAllLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// all
func NewOrganizationUserAllLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationUserAllLogic {
	return &OrganizationUserAllLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationUserAllLogic) OrganizationUserAll(req *types.EmptyRequest) (resp *types.OrganizationUserAllResponse, err error) {
	rpcReq := &lender.EmptyRequest{}

	rpcResp, err := l.svcCtx.OrganizationUserRpcClient.All(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	var list []*types.OrganizationUserEntity
	for _, v := range rpcResp.Data {
		item, err := copierx.Trans[lender.OrganizationUserEntity, types.OrganizationUserEntity](v, copierx.TimeToRFC3339Str())
		if err != nil {
			return nil, err
		}
		list = append(list, item)
	}

	resp = &types.OrganizationUserAllResponse{
		List: list,
	}

	return
}

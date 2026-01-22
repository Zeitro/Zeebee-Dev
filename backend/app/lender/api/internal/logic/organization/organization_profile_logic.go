package organization

import (
	"context"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrganizationProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// profile
func NewOrganizationProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationProfileLogic {
	return &OrganizationProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationProfileLogic) OrganizationProfile(req *types.EmptyRequest) (resp *types.OrganizationEntity, err error) {
	// Call RPC client to get organization profile
	rpcResp, err := l.svcCtx.OrganizationRpcClient.Profile(l.ctx, &lender.EmptyRequest{})
	if err != nil {
		return nil, err
	}

	// Convert proto response to API types
	resp, err = copierx.Trans[lender.OrganizationEntity, types.OrganizationEntity](
		rpcResp.Data,
		copierx.TimeToRFC3339Str(),
	)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

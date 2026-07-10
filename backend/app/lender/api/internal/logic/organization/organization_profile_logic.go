package organization

import (
	"context"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/proto"

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

func (l *OrganizationProfileLogic) OrganizationProfile(req *types.EmptyRequest) (resp *types.OrganizationProfileEntity, err error) {
	// Call RPC client to get organization profile
	rpcResp, err := l.svcCtx.OrganizationRpcClient.Profile(l.ctx, &proto.EmptyRequest{})
	if err != nil {
		return nil, err
	}

	// Convert proto response to API types
	resp, err = copierx.Trans[proto.OrganizationProfileEntity, types.OrganizationProfileEntity](
		rpcResp.Data,
		copierx.TimeToRFC3339Str(),
		copierx.BytesToStrMap(),
	)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

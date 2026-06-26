package organization_subscription

import (
	"context"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrganizationSubscriptionCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// create
func NewOrganizationSubscriptionCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationSubscriptionCreateLogic {
	return &OrganizationSubscriptionCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationSubscriptionCreateLogic) OrganizationSubscriptionCreate(req *types.OrganizationSubscriptionCreateRequest) (resp *types.OrganizationSubscriptionDetailResponse, err error) {

	// Convert types to pb request
	rpcReq := &admin.OrganizationSubscriptionCreateRequest{}
	rpcReq.Data, err = copierx.Trans[types.OrganizationSubscriptionCreateRequest, admin.OrganizationSubscriptionCreateRequestData](req)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	// Call RPC server
	rpcResp, err := l.svcCtx.OrganizationSubscriptionRpcClient.Create(l.ctx, rpcReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	// Convert pb response to types
	resp, err = copierx.Trans[admin.OrganizationSubscriptionEntity, types.OrganizationSubscriptionDetailResponse](rpcResp.Data, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	return resp, nil
}

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

type OrganizationSubscriptionDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// detail
func NewOrganizationSubscriptionDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationSubscriptionDetailLogic {
	return &OrganizationSubscriptionDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationSubscriptionDetailLogic) OrganizationSubscriptionDetail(req *types.PathIdRequest) (resp *types.OrganizationSubscriptionDetailResponse, err error) {
	rpcReq := &admin.PathIdRequest{
		Data: &admin.PathIdRequestData{
			Id: req.Id,
		},
	}

	rpcResp, err := l.svcCtx.OrganizationSubscriptionRpcClient.Detail(l.ctx, rpcReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	resp, err = copierx.Trans[admin.OrganizationSubscriptionEntity, types.OrganizationSubscriptionDetailResponse](rpcResp.Data, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	return resp, nil
}

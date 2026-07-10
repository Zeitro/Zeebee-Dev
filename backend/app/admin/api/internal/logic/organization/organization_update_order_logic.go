package organization

import (
	"context"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrganizationUpdateOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOrganizationUpdateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationUpdateOrderLogic {
	return &OrganizationUpdateOrderLogic{
		Logger: svcCtx.Logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationUpdateOrderLogic) OrganizationUpdateOrder(req *types.OrganizationUpdateOrderRequest) (resp *types.EmptyResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	_, err = l.svcCtx.OrganizationRpcClient.UpdateOrder(l.ctx, &admin.OrganizationUpdateOrderRequest{
		Data: &admin.OrganizationUpdateOrderRequestData{
			Id:    req.Id,
			Order: int32(req.Order),
		},
	})
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	return &types.EmptyResponse{}, nil
}

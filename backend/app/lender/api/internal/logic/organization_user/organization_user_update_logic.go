package organization_user

import (
	"context"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrganizationUserUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// update
func NewOrganizationUserUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationUserUpdateLogic {
	return &OrganizationUserUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationUserUpdateLogic) OrganizationUserUpdate(req *types.OrganizationUserUpdateRequest) (resp *types.OrganizationUserEntity, err error) {
	data, err := copierx.Trans[types.OrganizationUserUpdateRequest, lender.OrganizationUserUpdateRequestData](req)
	if err != nil {
		return nil, err
	}

	rpcReq := &lender.OrganizationUserUpdateRequest{
		Data: data,
	}

	rpcResp, err := l.svcCtx.OrganizationUserRpcClient.Update(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp, err = copierx.Trans[lender.OrganizationUserEntity, types.OrganizationUserEntity](rpcResp.Data, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}

	return
}

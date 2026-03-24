package organization_user

import (
	"context"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"

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

func (l *OrganizationUserUpdateLogic) OrganizationUserUpdate(req *types.OrganizationUserUpdateRequest) (resp *types.OrganizationUserDetailResponse, err error) {
	resp = &types.OrganizationUserDetailResponse{}
	rpcReq := &admin.OrganizationUserUpdateRequest{
		Data: &admin.OrganizationUserUpdateRequestData{},
	}

	rpcReq.Data, err = copierx.Trans[types.OrganizationUserUpdateRequest, admin.OrganizationUserUpdateRequestData](req)
	if err != nil {
		err = errorx.GetFailedTransReqErr(err)
		return nil, err
	}

	rpcResp, err := l.svcCtx.OrganizationUserRpcClient.Update(l.ctx, rpcReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	// trans response
	var entity *types.OrganizationUserEntity
	entity, err = copierx.Trans[admin.OrganizationUserEntity, types.OrganizationUserEntity](rpcResp.Data)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	resp = &types.OrganizationUserDetailResponse{
		OrganizationUserEntity: *entity,
	}
	return resp, nil
}

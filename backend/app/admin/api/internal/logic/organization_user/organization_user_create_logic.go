package organization_user

import (
	"context"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrganizationUserCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// create
func NewOrganizationUserCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationUserCreateLogic {
	return &OrganizationUserCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationUserCreateLogic) OrganizationUserCreate(req *types.OrganizationUserCreateRequest) (resp *types.OrganizationUserDetailResponse, err error) {
	resp = &types.OrganizationUserDetailResponse{}
	rpcReq := &admin.OrganizationUserCreateRequest{
		Data: &admin.OrganizationUserCreateRequestData{},
	}

	rpcReq.Data, err = copierx.Trans[types.OrganizationUserCreateRequest, admin.OrganizationUserCreateRequestData](req)
	if err != nil {
		err = errorx.GetFailedTransReqErr(err)
		return nil, err
	}

	rpcResp, err := l.svcCtx.OrganizationUserRpcClient.Create(l.ctx, rpcReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	// trans response
	var entity *types.OrganizationUserEntity
	entity, err = copierx.Trans[proto.OrganizationUserEntity, types.OrganizationUserEntity](rpcResp.Data)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	resp = &types.OrganizationUserDetailResponse{
		OrganizationUserEntity: *entity,
	}
	return resp, nil
}

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

type OrganizationUserPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// page
func NewOrganizationUserPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationUserPageLogic {
	return &OrganizationUserPageLogic{
		Logger: svcCtx.Logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationUserPageLogic) OrganizationUserPage(req *types.OrganizationUserPageRequest) (resp *types.OrganizationUserPageResponse, err error) {
	resp = &types.OrganizationUserPageResponse{}
	rpcReq := &admin.OrganizationUserPageRequest{
		Data: &admin.OrganizationUserPageRequestData{},
	}

	rpcReq.Data, err = copierx.Trans[types.OrganizationUserPageRequest, admin.OrganizationUserPageRequestData](req)
	if err != nil {
		err = errorx.GetFailedTransReqErr(err)
		return nil, err
	}

	rpcResp, err := l.svcCtx.OrganizationUserRpcClient.Page(l.ctx, rpcReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	// trans response
	resp, err = copierx.Trans[admin.OrganizationUserPageResponseData, types.OrganizationUserPageResponse](rpcResp.Data)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	return resp, err
}

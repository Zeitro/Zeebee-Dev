package organization_user

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
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
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationUserPageLogic) OrganizationUserPage(req *types.OrganizationUserPageRequest) (resp *types.OrganizationUserPageResponse, err error) {
	rpcReq := &borrower.OrganizationUserPageRequest{
		Data: &borrower.OrganizationUserPageRequestData{},
	}
	rpcReq.Data, err = copierx.Trans[types.OrganizationUserPageRequest, borrower.OrganizationUserPageRequestData](req)
	if err != nil {
		return nil, err
	}

	rpcResp, err := l.svcCtx.OrganizationUserRpcClient.Page(l.ctx, rpcReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	resp = &types.OrganizationUserPageResponse{}
	resp, err = copierx.Trans[borrower.OrganizationUserPageResponseData, types.OrganizationUserPageResponse](rpcResp.Data, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}

	return resp, nil
}

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

type OrganizationUserDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// detail
func NewOrganizationUserDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationUserDetailLogic {
	return &OrganizationUserDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationUserDetailLogic) OrganizationUserDetail(req *types.PathIdRequest) (resp *types.OrganizationUserEntity, err error) {
	rpcReq := &borrower.OrganizationUserDetailRequest{
		Data: &borrower.OrganizationUserDetailRequestData{},
	}
	rpcReq.Data, err = copierx.Trans[types.PathIdRequest, borrower.OrganizationUserDetailRequestData](req)
	if err != nil {
		return nil, err
	}

	rpcResp, err := l.svcCtx.OrganizationUserRpcClient.Detail(l.ctx, rpcReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	resp = &types.OrganizationUserEntity{}
	resp, err = copierx.Trans[borrower.OrganizationUserEntity, types.OrganizationUserEntity](rpcResp.Data, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}

	return resp, nil
}

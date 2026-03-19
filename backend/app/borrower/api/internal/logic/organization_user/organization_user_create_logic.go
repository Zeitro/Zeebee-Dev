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

func (l *OrganizationUserCreateLogic) OrganizationUserCreate(req *types.OrganizationUserCreateRequest) (resp *types.OrganizationUserEntity, err error) {
	rpcReq := &borrower.OrganizationUserCreateRequest{
		Data: &borrower.OrganizationUserCreateRequestData{},
	}
	rpcReq.Data, err = copierx.Trans[types.OrganizationUserCreateRequest, borrower.OrganizationUserCreateRequestData](req)
	if err != nil {
		return nil, err
	}

	rpcResp, err := l.svcCtx.OrganizationUserRpcClient.Create(l.ctx, rpcReq)
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

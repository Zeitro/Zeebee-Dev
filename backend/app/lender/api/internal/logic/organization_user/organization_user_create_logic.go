package organization_user

import (
	"context"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"

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
	// Transform API request to RPC request
	data, err := copierx.Trans[types.OrganizationUserCreateRequest, lender.OrganizationUserCreateRequestData](req)
	if err != nil {
		return nil, err
	}

	rpcReq := &lender.OrganizationUserCreateRequest{
		Data: data,
	}

	// Call RPC service
	rpcResp, err := l.svcCtx.OrganizationUserRpcClient.Create(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// Transform RPC response to API response
	resp, err = copierx.Trans[lender.OrganizationUserEntity, types.OrganizationUserEntity](rpcResp.Data, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}

	return
}

package organization_user

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/proto"

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
	rpcReq := &borrower.OrganizationUserUpdateRequest{
		Data: &borrower.OrganizationUserUpdateRequestData{},
	}
	rpcReq.Data, err = copierx.Trans[types.OrganizationUserUpdateRequest, borrower.OrganizationUserUpdateRequestData](req)
	if err != nil {
		return nil, err
	}

	rpcResp, err := l.svcCtx.OrganizationUserRpcClient.Update(l.ctx, rpcReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	resp = &types.OrganizationUserEntity{}
	resp, err = copierx.Trans[proto.OrganizationUserEntity, types.OrganizationUserEntity](rpcResp.Data, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}

	return resp, nil
}

package organization_user

import (
	"context"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"

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
	rpcReq := &lender.OrganizationUserDetailRequest{
		Data: &lender.OrganizationUserDetailRequestData{
			Id: req.Id,
		},
	}

	rpcResp, err := l.svcCtx.OrganizationUserRpcClient.Detail(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp, err = copierx.Trans[lender.OrganizationUserEntity, types.OrganizationUserEntity](rpcResp.Data, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}

	return
}

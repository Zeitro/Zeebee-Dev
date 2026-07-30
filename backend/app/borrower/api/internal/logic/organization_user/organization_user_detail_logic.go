package organization_user

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/proto"

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
	rpcReq := &proto.PathIdRequest{
		Data: &proto.PathIdRequestData{},
	}
	rpcReq.Data, err = copierx.Trans[types.PathIdRequest, proto.PathIdRequestData](req)
	if err != nil {
		return nil, err
	}

	rpcResp, err := l.svcCtx.OrganizationUserRpcClient.Detail(l.ctx, rpcReq)
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

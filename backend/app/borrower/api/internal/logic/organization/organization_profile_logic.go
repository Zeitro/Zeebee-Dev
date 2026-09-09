package organization

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrganizationProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// info
func NewOrganizationProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationProfileLogic {
	return &OrganizationProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationProfileLogic) OrganizationProfile(req *types.EmptyRequest) (resp *types.OrganizationEntity, err error) {
	rpcReq := &proto.EmptyRequest{}

	rpcResp, err := l.svcCtx.OrganizationRpcClient.Profile(l.ctx, rpcReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	resp, err = copierx.Trans[proto.OrganizationEntity, types.OrganizationEntity](rpcResp.Data)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

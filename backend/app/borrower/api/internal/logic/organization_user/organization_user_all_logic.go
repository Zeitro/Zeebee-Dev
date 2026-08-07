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

type OrganizationUserAllLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// all
func NewOrganizationUserAllLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationUserAllLogic {
	return &OrganizationUserAllLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationUserAllLogic) OrganizationUserAll(req *types.EmptyRequest) (resp *types.OrganizationUserAllResponse, err error) {
	rpcReq := &proto.EmptyRequest{}
	rpcResp, err := l.svcCtx.OrganizationUserRpcClient.All(l.ctx, rpcReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	resp = &types.OrganizationUserAllResponse{
		List: []*types.OrganizationUserEntity{},
	}
	listPtr, err := copierx.Trans[[]*proto.OrganizationUserEntity, []*types.OrganizationUserEntity](&rpcResp.Data, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}
	resp.List = *listPtr
	return resp, nil
}

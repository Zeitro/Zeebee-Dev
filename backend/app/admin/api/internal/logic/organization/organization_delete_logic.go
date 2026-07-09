package organization

import (
	"context"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrganizationDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// delete organization
func NewOrganizationDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationDeleteLogic {
	return &OrganizationDeleteLogic{
		Logger: svcCtx.Logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationDeleteLogic) OrganizationDelete(req *types.PathIdRequest) (resp *types.EmptyResponse, err error) {

	// Call RPC server
	_, err = l.svcCtx.OrganizationRpcClient.Delete(l.ctx, &proto.IdRequest{
		Data: req.Id,
	})
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	return &types.EmptyResponse{}, nil
}

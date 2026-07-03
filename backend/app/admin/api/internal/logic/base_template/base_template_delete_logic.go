package base_template

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/proto"
)

type BaseTemplateDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Delete base template
func NewBaseTemplateDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BaseTemplateDeleteLogic {
	return &BaseTemplateDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BaseTemplateDeleteLogic) BaseTemplateDelete(req *types.PathIdRequest) (resp *types.EmptyResponse, err error) {

	// Convert types to pb request
	rpcReq := &admin.DeleteBaseTemplateRequest{
		Data: &proto.IdRequest{Data: req.Id},
	}

	// Call RPC server
	_, err = l.svcCtx.BaseTemplateRpcClient.Delete(l.ctx, rpcReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	return &types.EmptyResponse{}, nil
}

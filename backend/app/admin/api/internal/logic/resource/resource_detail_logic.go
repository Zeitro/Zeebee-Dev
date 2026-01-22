package resource

import (
	"context"

	"github.com/jinzhu/copier"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type ResourceDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// detail
func NewResourceDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResourceDetailLogic {
	return &ResourceDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ResourceDetailLogic) ResourceDetail(req *types.ResourceDetailRequest) (resp *types.ResourceEntity, err error) {
	resp = &types.ResourceEntity{}
	rpcReq := &admin.ResourceDetailRequest{
		Data: &admin.ResourceDetailRequestData{},
	}
	err = copier.CopyWithOption(rpcReq.Data, req, copier.Option{DeepCopy: true})
	if err != nil {
		return nil, errorx.GetFailedTransReqErr(err)
	}

	rpcResp, err := l.svcCtx.ResourceRpcClient.Detail(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// trans response
	err = copier.CopyWithOption(resp, &rpcResp.Data, copier.Option{DeepCopy: true})
	if err != nil {
		return nil, errorx.GetFailedTransReqErr(err)
	}
	return resp, err
}

package organization

import (
	"context"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrganizationDiscoverPreviewLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOrganizationDiscoverPreviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationDiscoverPreviewLogic {
	return &OrganizationDiscoverPreviewLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationDiscoverPreviewLogic) LenderDiscoverPreview(_ *types.EmptyRequest) (resp *types.LenderDiscoverPreviewResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	rpcResp, err := l.svcCtx.OrganizationRpcClient.DiscoverPreview(context.WithoutCancel(l.ctx), &proto.EmptyRequest{})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	var data types.LenderDiscoverPreviewData
	if rpcResp.GetData() != nil {
		data, err = discoverProtoPreviewDataToTypes(rpcResp.GetData())
		if err != nil {
			return nil, err
		}
	}
	return &types.LenderDiscoverPreviewResponse{Data: data}, nil
}

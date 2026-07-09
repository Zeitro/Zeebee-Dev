package organization

import (
	"context"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrganizationDiscoverUpdateSettingsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOrganizationDiscoverUpdateSettingsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationDiscoverUpdateSettingsLogic {
	return &OrganizationDiscoverUpdateSettingsLogic{
		Logger: svcCtx.Logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationDiscoverUpdateSettingsLogic) LenderDiscoverUpdateSettings(req *types.LenderDiscoverUpdateSettingsRequest) (resp *types.LenderDiscoverSettingsResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	if len(req.OpenaiModel) == 0 {
		req.OpenaiModel = "gpt-4o-mini"
	}
	rpcData, err := copierx.Trans[types.LenderDiscoverUpdateSettingsRequest, admin.LenderDiscoverUpdateSettingsRequestData](req)
	if err != nil {
		return nil, err
	}
	rpcData.UpdatedByUserId = 0
	rpcResp, err := l.svcCtx.OrganizationRpcClient.DiscoverUpdateSettings(l.ctx, &admin.LenderDiscoverUpdateSettingsRequest{
		Data: rpcData,
	})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	data, err := discoverProtoSettingsToTypes(rpcResp.GetData())
	if err != nil {
		return nil, err
	}
	return &types.LenderDiscoverSettingsResponse{Data: data}, nil
}

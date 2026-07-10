package organization

import (
	"context"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrganizationDiscoverGetSettingsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOrganizationDiscoverGetSettingsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationDiscoverGetSettingsLogic {
	return &OrganizationDiscoverGetSettingsLogic{
		Logger: svcCtx.Logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationDiscoverGetSettingsLogic) LenderDiscoverGetSettings(_ *types.EmptyRequest) (resp *types.LenderDiscoverSettingsResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	rpcResp, err := l.svcCtx.OrganizationRpcClient.DiscoverGetSettings(l.ctx, &proto.EmptyRequest{})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	data, err := discoverProtoSettingsToTypes(rpcResp.GetData())
	if err != nil {
		return nil, err
	}
	return &types.LenderDiscoverSettingsResponse{Data: data}, nil
}

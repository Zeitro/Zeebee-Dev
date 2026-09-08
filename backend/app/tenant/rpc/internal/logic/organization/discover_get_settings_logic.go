package organizationlogic

import (
	"context"

	"zeebee-backend/app/tenant/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/domains/zeebee/lenderdiscover"
	"zeebee-backend/domains/zeebee/proto"
)

type DiscoverGetSettingsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDiscoverGetSettingsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DiscoverGetSettingsLogic {
	return &DiscoverGetSettingsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DiscoverGetSettingsLogic) DiscoverGetSettings(in *proto.EmptyRequest) (resp *tenant.LenderDiscoverSettingsResponse, err error) {
	cfg, _, err := lenderdiscover.LoadConfig(l.ctx, l.svcCtx.OrganizationSettingRepo)
	if err != nil {
		return nil, err
	}
	cfg.Model = l.svcCtx.Config.LlmConf.LenderDiscoverConf.Model
	data, err := discoverConfigToProto(cfg)
	if err != nil {
		return nil, err
	}
	return &tenant.LenderDiscoverSettingsResponse{Data: data}, nil
}

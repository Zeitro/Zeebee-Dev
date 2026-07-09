package organizationlogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/lenderdiscover"
	"zeebee-backend/domains/zeebee/proto"
)

type DiscoverGetSettingsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDiscoverGetSettingsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DiscoverGetSettingsLogic {
	return &DiscoverGetSettingsLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *DiscoverGetSettingsLogic) DiscoverGetSettings(_ *proto.EmptyRequest) (resp *admin.LenderDiscoverSettingsResponse, err error) {
	defer func() {
		if err != nil {
			l.svcCtx.Logger.WithContext(l.ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
			).Error(err)
		}
	}()
	cfg, _, err := lenderdiscover.LoadConfig(l.ctx, l.svcCtx.OrganizationSettingRepo)
	if err != nil {
		return nil, err
	}
	data, err := discoverConfigToProto(cfg)
	if err != nil {
		return nil, err
	}
	return &admin.LenderDiscoverSettingsResponse{Data: data}, nil
}

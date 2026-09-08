package organizationlogic

import (
	"context"

	"zeebee-backend/app/tenant/rpc/internal/schedule"
	"zeebee-backend/app/tenant/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/lenderdiscover"
)

type DiscoverUpdateSettingsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDiscoverUpdateSettingsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DiscoverUpdateSettingsLogic {
	return &DiscoverUpdateSettingsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DiscoverUpdateSettingsLogic) DiscoverUpdateSettings(in *tenant.LenderDiscoverUpdateSettingsRequest) (resp *tenant.LenderDiscoverSettingsResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	d := in.GetData()
	if d.IntervalDays < 1 {
		return nil, errorx.New("interval_days must be at least 1")
	}
	_, rowID, err := lenderdiscover.LoadConfig(l.ctx, l.svcCtx.OrganizationSettingRepo)
	if err != nil {
		return nil, err
	}
	cfg, err := copierx.Trans[tenant.LenderDiscoverUpdateSettingsRequestData, lenderdiscover.Config](d)
	if err != nil {
		return nil, err
	}
	if err = lenderdiscover.SaveConfig(l.ctx, l.svcCtx.OrganizationSettingRepo, rowID, cfg); err != nil {
		return nil, err
	}
	syncDiscoverCron(l.svcCtx, cfg)
	data, err := discoverConfigToProto(cfg)
	if err != nil {
		return nil, err
	}
	return &tenant.LenderDiscoverSettingsResponse{Data: data}, nil
}

func syncDiscoverCron(svcCtx *svc.ServiceContext, cfg *lenderdiscover.Config) {
	if svcCtx.Cron == nil {
		return
	}
	if cfg.ScheduleEnabled {
		_ = svcCtx.Cron.Replace(schedule.TaskName, schedule.SpecFromInterval(cfg.IntervalDays), func() {
			schedule.RunScheduled(context.Background(), svcCtx)
		})
	} else {
		_ = svcCtx.Cron.Unregister(schedule.TaskName)
	}
}

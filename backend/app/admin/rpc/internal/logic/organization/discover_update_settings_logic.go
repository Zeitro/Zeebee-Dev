package organizationlogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/schedule"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/cron"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/lenderdiscover"
)

type DiscoverUpdateSettingsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDiscoverUpdateSettingsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DiscoverUpdateSettingsLogic {
	return &DiscoverUpdateSettingsLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *DiscoverUpdateSettingsLogic) DiscoverUpdateSettings(in *admin.LenderDiscoverUpdateSettingsRequest) (resp *admin.LenderDiscoverSettingsResponse, err error) {

	if in.GetData() == nil {
		return nil, errorx.New("Request data is required")
	}
	d := in.GetData()
	if d.IntervalDays < 1 {
		return nil, errorx.New("interval_days must be at least 1")
	}
	_, rowId, err := lenderdiscover.LoadConfig(l.ctx, l.svcCtx.OrganizationSettingRepo)
	if err != nil {
		return nil, err
	}
	cfg, err := copierx.Trans[admin.LenderDiscoverUpdateSettingsRequestData, lenderdiscover.Config](d)
	if err != nil {
		return nil, err
	}
	if err = lenderdiscover.SaveConfig(l.ctx, l.svcCtx.OrganizationSettingRepo, rowId, cfg); err != nil {
		return nil, err
	}
	syncDiscoverCron(l.svcCtx, cfg)
	data, err := discoverConfigToProto(cfg)
	if err != nil {
		return nil, err
	}
	return &admin.LenderDiscoverSettingsResponse{Data: data}, nil
}

// syncDiscoverCron re-arms the lender discover cron task to match the just-saved
// config: Replace with the new interval when enabled, Unregister when disabled.
func syncDiscoverCron(svcCtx *svc.ServiceContext, cfg *lenderdiscover.Config) {
	if cfg.ScheduleEnabled {
		_ = cron.Replace(schedule.TaskName, schedule.SpecFromInterval(cfg.IntervalDays), func() {
			schedule.RunScheduled(context.Background(), svcCtx)
		})
	} else {
		_ = cron.Unregister(schedule.TaskName)
	}
}

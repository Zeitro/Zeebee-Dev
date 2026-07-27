package schedule

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	domainlenderdiscover "zeebee-backend/domains/zeebee/lenderdiscover"
)

// TaskName is the cron entry name for lender discover.
const TaskName = "lender_discover"

// SpecFromInterval builds an "@every" cron spec from IntervalDays. The spec
// carries the cadence so a config change + cron.Replace actually re-arms the
// timer. Min 1 day; default 7.
func SpecFromInterval(days int32) string {
	if days < 1 {
		days = 7
	}
	return fmt.Sprintf("@every %dh", int(days)*24)
}

// NewRunner builds a discover Runner from svcCtx (repos + logic + injected LLM).
// Shared by the scheduled run and the manual Preview/Commit/Run RPCs.
func NewRunner(svcCtx *svc.ServiceContext) domainlenderdiscover.Runner {
	bl := logic.GetBaseLogic(svcCtx)
	return domainlenderdiscover.Runner{
		Repos:              svcCtx.AllRepos,
		TemplateLogic:      bl.OrganizationTemplateLogic,
		SettingLogic:       bl.OrganizationSettingLogic,
		LLM:                svcCtx.DiscoverLLM,
		PlaceholderLogoURL: domainlenderdiscover.DefaultPlaceholderLogoURL,
	}
}

// RunScheduled is the cron task body: insert all non-duplicate lenders. No gate,
// no lock — the @every spec owns the cadence and admin runs single-instance.
// ponytail: if admin scales horizontally, wrap this body in a Redis SetNX lock.
func RunScheduled(ctx context.Context, svcCtx *svc.ServiceContext) {
	rn := NewRunner(svcCtx)
	if _, err := rn.Run(ctx, domainlenderdiscover.TriggerScheduled); err != nil {
		logx.WithContext(ctx).Errorf("lender discover scheduled run failed: %v", err)
	}
}

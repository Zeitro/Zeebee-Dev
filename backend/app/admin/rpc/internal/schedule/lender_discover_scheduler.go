package schedule

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/admin/rpc/internal/global"
	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	domainlenderdiscover "zeebee-backend/domains/zeebee/lenderdiscover"
)

const schedulerLockKey = "lender_discover:scheduler"
const schedulerLockTTL = 45 * time.Minute
const tickInterval = 5 * time.Minute

const lastScheduledSuccessKey = "lender_discover:last_scheduled_success"

// StartLenderDiscoverScheduler runs a background ticker that may trigger scheduled discover jobs.
func StartLenderDiscoverScheduler(svcCtx *svc.ServiceContext) {
	go func() {
		ticker := time.NewTicker(tickInterval)
		defer ticker.Stop()
		for range ticker.C {
			tryRunScheduled(context.Background(), svcCtx)
		}
	}()
}

func tryRunScheduled(ctx context.Context, svcCtx *svc.ServiceContext) {
	cfg, _, err := domainlenderdiscover.LoadConfig(ctx, svcCtx.OrganizationSettingRepo)
	if err != nil {
		logx.WithContext(ctx).Errorf("lender discover scheduler: load settings: %v", err)
		return
	}
	if !cfg.ScheduleEnabled {
		return
	}
	apiKey := strings.TrimSpace(svcCtx.Config.OpenAIConf.APIKey)
	if apiKey == "" {
		apiKey = os.Getenv("OPENAI_API_KEY")
	}
	if apiKey == "" {
		return
	}

	rds := global.GetRedis()
	if rds == nil {
		return
	}

	now := time.Now().UTC()
	lastStr, err := rds.Get(ctx, lastScheduledSuccessKey)
	if err != nil {
		logx.WithContext(ctx).Errorf("lender discover scheduler: redis get last run: %v", err)
		return
	}
	if lastStr != "" {
		lastTime, parseErr := time.Parse(time.RFC3339, lastStr)
		if parseErr == nil {
			due := lastTime.AddDate(0, 0, int(cfg.IntervalDays))
			if now.Before(due) {
				return
			}
		}
	}

	ok, err := rds.SetNX(ctx, schedulerLockKey, "1", schedulerLockTTL)
	if err != nil || !ok {
		return
	}

	bl := logic.GetBaseLogic(svcCtx)
	rn := domainlenderdiscover.Runner{
		Repos:              svcCtx.AllRepos,
		TemplateLogic:      bl.OrganizationTemplateLogic,
		SettingLogic:       bl.OrganizationSettingLogic,
		OpenAIKey:          apiKey,
		PlaceholderLogoURL: domainlenderdiscover.DefaultPlaceholderLogoURL,
	}
	if _, err := rn.Run(ctx, domainlenderdiscover.TriggerScheduled); err != nil {
		logx.WithContext(ctx).Errorf("lender discover scheduled run failed: %v", err)
		return
	}

	ttlSeconds := int(cfg.IntervalDays) * 2 * 86400
	if setErr := rds.SetExpire(ctx, lastScheduledSuccessKey, now.Format(time.RFC3339), ttlSeconds); setErr != nil {
		logx.WithContext(ctx).Errorf("lender discover scheduler: redis set last run: %v", setErr)
	}
}

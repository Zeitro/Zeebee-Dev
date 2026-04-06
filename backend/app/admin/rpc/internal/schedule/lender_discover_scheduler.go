package schedule

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/global"
	"zeebee-backend/app/admin/rpc/internal/svc"
	domainlenderdiscover "zeebee-backend/domains/zeebee/lenderdiscover"
)

const schedulerLockKey = "lender_discover:scheduler"
const schedulerLockTTL = 45 * time.Minute
const tickInterval = 5 * time.Minute

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
	settings, err := svcCtx.LenderDiscoverSettingsRepo.GetSingleton(ctx)
	if err != nil {
		logx.WithContext(ctx).Errorf("lender discover scheduler: load settings: %v", err)
		return
	}
	if !settings.ScheduleEnabled {
		return
	}
	apiKey := strings.TrimSpace(svcCtx.Config.OpenAIConf.APIKey)
	if apiKey == "" {
		apiKey = os.Getenv("OPENAI_API_KEY")
	}
	if apiKey == "" {
		return
	}

	lastFinPtr, err := svcCtx.LenderDiscoverRunRepo.LastSuccessfulScheduledFinished(ctx)
	if err != nil {
		logx.WithContext(ctx).Errorf("lender discover scheduler: last run: %v", err)
		return
	}
	now := time.Now().UTC()
	if lastFinPtr != nil {
		due := lastFinPtr.AddDate(0, 0, int(settings.IntervalDays))
		if now.Before(due) {
			return
		}
	}

	rds := global.GetRedis()
	if rds == nil {
		return
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
	}
}

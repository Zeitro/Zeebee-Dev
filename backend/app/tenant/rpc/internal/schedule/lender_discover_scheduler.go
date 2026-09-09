package schedule

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/tenant/rpc/internal/logic"
	"zeebee-backend/app/tenant/rpc/internal/svc"
	domainlenderdiscover "zeebee-backend/domains/zeebee/lenderdiscover"
)

const TaskName = "lender_discover"

func SpecFromInterval(days int32) string {
	if days < 1 {
		days = 7
	}
	return fmt.Sprintf("@every %dh", int(days)*24)
}

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

func RunScheduled(ctx context.Context, svcCtx *svc.ServiceContext) {
	rn := NewRunner(svcCtx)
	if _, err := rn.Run(ctx, domainlenderdiscover.TriggerScheduled); err != nil {
		logx.WithContext(ctx).Errorf("lender discover scheduled run failed: %v", err)
	}
}

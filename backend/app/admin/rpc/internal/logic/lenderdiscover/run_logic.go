package lenderdiscover

import (
	"context"
	"os"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/errorx"
	domainlenderdiscover "zeebee-backend/domains/zeebee/lenderdiscover"
)

type RunLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRunLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RunLogic {
	return &RunLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *RunLogic) Run(in *admin.LenderDiscoverRunRequest) (resp *admin.LenderDiscoverRunResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	trigger := domainlenderdiscover.TriggerManual
	if in.GetData() != nil && in.GetData().GetTrigger() == domainlenderdiscover.TriggerScheduled {
		trigger = domainlenderdiscover.TriggerScheduled
	}
	if trigger != domainlenderdiscover.TriggerScheduled {
		return nil, errorx.New("manual discover: use Preview then Commit from the admin UI")
	}

	bl := logic.GetBaseLogic(l.svcCtx)
	apiKey := strings.TrimSpace(l.svcCtx.Config.OpenAIConf.APIKey)
	if apiKey == "" {
		apiKey = os.Getenv("OPENAI_API_KEY")
	}
	rn := domainlenderdiscover.Runner{
		Repos:              l.svcCtx.AllRepos,
		TemplateLogic:      bl.OrganizationTemplateLogic,
		SettingLogic:       bl.OrganizationSettingLogic,
		OpenAIKey:          apiKey,
		PlaceholderLogoURL: domainlenderdiscover.DefaultPlaceholderLogoURL,
	}
	runRow, err := rn.Run(context.WithoutCancel(l.ctx), trigger)
	if err != nil {
		return nil, err
	}
	// Reload for latest stats/error persisted by runner
	fresh, err := l.svcCtx.LenderDiscoverRunRepo.GetById(l.ctx, runRow.Id)
	if err != nil {
		return &admin.LenderDiscoverRunResponse{Data: runToProto(runRow)}, nil
	}
	return &admin.LenderDiscoverRunResponse{Data: runToProto(fresh)}, nil
}

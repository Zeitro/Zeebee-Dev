package organizationlogic

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

type DiscoverCommitLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDiscoverCommitLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DiscoverCommitLogic {
	return &DiscoverCommitLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DiscoverCommitLogic) DiscoverCommit(in *admin.LenderDiscoverCommitRequest) (resp *admin.LenderDiscoverRunResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	var lenders []*admin.LenderDiscoverCommitLender
	if in.GetData() != nil {
		lenders = in.GetData().GetLenders()
	}
	raws := make([]domainlenderdiscover.RawLender, 0, len(lenders))
	for _, p := range lenders {
		raw, transErr := discoverRawLenderFromCommitProto(p)
		if transErr != nil {
			return nil, transErr
		}
		raws = append(raws, raw)
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
	result, err := rn.Commit(context.WithoutCancel(l.ctx), raws)
	if err != nil {
		return nil, err
	}
	return &admin.LenderDiscoverRunResponse{Data: discoverRunResultToProto(result)}, nil
}

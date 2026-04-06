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

type PreviewLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPreviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PreviewLogic {
	return &PreviewLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PreviewLogic) Preview(_ *admin.EmptyRequest) (resp *admin.LenderDiscoverPreviewResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
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
	cands, err := rn.Preview(context.WithoutCancel(l.ctx))
	if err != nil {
		return nil, err
	}
	var protoList []*admin.LenderDiscoverCandidateEntity
	var newN, existedN int32
	for i := range cands {
		c := &cands[i]
		protoList = append(protoList, candidateToProto(c))
		if c.DedupeStatus == domainlenderdiscover.DedupeNew {
			newN++
		} else {
			existedN++
		}
	}
	return &admin.LenderDiscoverPreviewResponse{
		Data: &admin.LenderDiscoverPreviewResponseData{
			Candidates:    protoList,
			NewCount:      newN,
			ExistedCount:  existedN,
		},
	}, nil
}

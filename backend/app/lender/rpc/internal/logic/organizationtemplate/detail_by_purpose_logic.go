package organizationtemplatelogic

import (
	"context"

	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/data/bo"
	"zeebee-backend/domains/zeebee/repo"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type DetailByPurposeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
	repo      *repo.AllRepos
}

func NewDetailByPurposeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailByPurposeLogic {
	return &DetailByPurposeLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.NewBaseLogic(ctx, svcCtx),
	}
}

func (l *DetailByPurposeLogic) DetailByPurpose(in *lender.OrganizationTemplateDetailByPurposeRequest) (resp *lender.OrganizationTemplateDetailResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	claims, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	detail, err := l.baseLogic.OrganizationTemplateLogic.DetailByPurpose(l.ctx, claims.OrganizationId, in.Data.Purpose)
	if err != nil {
		return nil, err
	}
	data, err := copierx.Trans[bo.OrganizationTemplateBo, lender.OrganizationTemplateEntity](
		detail,
		copierx.TimeToRFC3339Str(),
	)
	if err != nil {
		return nil, err
	}

	return &lender.OrganizationTemplateDetailResponse{
		Data: data,
	}, nil
}

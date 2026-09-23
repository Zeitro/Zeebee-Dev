package projectmanagementtemplatelogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/common/errorx"
)

type StatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StatusLogic {
	return &StatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StatusLogic) Status(in *admin.ProjectManagementTemplateStatusRequest) (resp *admin.ProjectManagementTemplateDetailResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	base := logic.GetBaseLogic(l.svcCtx)
	claim, err := base.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	template, err := base.ProjectManagementTemplateLogic.SetStatus(l.ctx, claim, in.Data.Id, in.Data.Status)
	if err != nil {
		return nil, err
	}
	return templateResponse(template)
}

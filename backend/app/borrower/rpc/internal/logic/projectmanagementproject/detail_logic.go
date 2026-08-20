package projectmanagementprojectlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"

	"zeebee-backend/common/errorx"
)

type DetailLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *DetailLogic) Detail(in *borrower.ProjectManagementProjectDetailRequest) (resp *borrower.ProjectManagementProjectDetailResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	claimValue, err := l.baseLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	claim := &claimValue
	project, err := l.baseLogic.ProjectManagementAccessLogic.RequireProject(l.ctx, in.Data.Id, claim.OrganizationId)
	if err != nil {
		return nil, err
	}
	if _, err = l.baseLogic.ProjectManagementAccessLogic.RequireMember(l.ctx, project.Id, claim.OrganizationUserId); err != nil {
		return nil, err
	}
	data, err := projectEntity(project)
	if err != nil {
		return nil, err
	}
	return &borrower.ProjectManagementProjectDetailResponse{Data: data}, nil
}

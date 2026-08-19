package projectmanagementlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
)

type DependencyUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDependencyUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DependencyUpdateLogic {
	return &DependencyUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DependencyUpdateLogic) DependencyUpdate(in *borrower.ProjectManagementDependencyUpdateRequest) (resp *borrower.ProjectManagementDependencyDetailResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	base := logic.GetBaseLogic(l.svcCtx)
	actor, err := requireProjectManagementActor(l.ctx, base)
	if err != nil {
		return nil, err
	}
	dependency, err := base.ProjectManagementLogic.UpdateTaskDependency(l.ctx, actor, in.Data.TaskId, in.Data.DependencyId, constants.ProjectManagementDependencyType(in.Data.DependencyType), constants.ProjectManagementDependencyRule(in.Data.DependencyRule))
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	data, err := dependencyEntity(dependency)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	return &borrower.ProjectManagementDependencyDetailResponse{Data: data}, nil
}

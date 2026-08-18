package projectmanagementlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
)

type DependencyCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDependencyCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DependencyCreateLogic {
	return &DependencyCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DependencyCreateLogic) DependencyCreate(in *borrower.ProjectManagementDependencyCreateRequest) (resp *borrower.ProjectManagementDependencyDetailResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	base := logic.GetBaseLogic(l.svcCtx)
	actor, err := requireProjectManagementActor(l.ctx, base)
	if err != nil {
		return nil, err
	}
	dependency, err := base.ProjectManagementLogic.CreateTaskDependency(l.ctx, actor, dto.ProjectManagementTaskDependencyDto{SuccessorTaskId: in.Data.SuccessorTaskId, PredecessorTaskId: in.Data.PredecessorTaskId, DependencyType: constants.ProjectManagementDependencyType(in.Data.DependencyType)})
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	data, err := dependencyEntity(dependency)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	return &borrower.ProjectManagementDependencyDetailResponse{Data: data}, nil
}

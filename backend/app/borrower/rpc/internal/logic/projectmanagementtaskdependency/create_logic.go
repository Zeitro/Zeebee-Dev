package projectmanagementtaskdependencylogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
)

type CreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLogic) Create(in *borrower.ProjectManagementDependencyCreateRequest) (resp *borrower.ProjectManagementDependencyDetailResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	base := logic.GetBaseLogic(l.svcCtx)
	claimValue, err := base.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	claim := &claimValue
	dependency, err := base.ProjectManagementDependencyLogic.Create(l.ctx, claim, dto.ProjectManagementTaskDependencyDto{TaskId: in.Data.TaskId, DependencyTaskId: in.Data.DependencyTaskId, DependencyType: constants.ProjectManagementDependencyType(in.Data.DependencyType), DependencyRule: constants.ProjectManagementDependencyRule(in.Data.DependencyRule)})
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	data, err := dependencyEntity(dependency)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	return &borrower.ProjectManagementDependencyDetailResponse{Data: data}, nil
}

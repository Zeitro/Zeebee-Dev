package projectmanagementtaskdependencylogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/proto"
)

type UpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLogic) Update(in *borrower.ProjectManagementDependencyUpdateRequest) (resp *borrower.ProjectManagementDependencyDetailResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	base := logic.GetBaseLogic(l.svcCtx)
	claimValue, err := base.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	claim := claimValue
	dependency, err := base.ProjectManagementDependencyLogic.Update(l.ctx, claim, dto.ProjectManagementTaskDependencyDto{Id: in.Data.Id, TaskId: in.Data.TaskId, DependencyType: constants.ProjectManagementDependencyType(in.Data.DependencyType), DependencyRule: constants.ProjectManagementDependencyRule(in.Data.DependencyRule)})
	if err != nil {
		return nil, err
	}
	data, err := copierx.Trans[model.ProjectManagementTaskDependency, proto.ProjectManagementTaskDependencyEntity](dependency, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}
	return &borrower.ProjectManagementDependencyDetailResponse{Data: data}, nil
}

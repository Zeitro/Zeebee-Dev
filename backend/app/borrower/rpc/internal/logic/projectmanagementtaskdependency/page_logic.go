package projectmanagementtaskdependencylogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/proto"
)

type PageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PageLogic {
	return &PageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PageLogic) Page(in *borrower.ProjectManagementDependencyPageRequest) (resp *borrower.ProjectManagementDependencyPageResponse, err error) {
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
	page, err := base.ProjectManagementDependencyLogic.Page(l.ctx, claim, dto.ProjectManagementTaskDependencyPageRequestDto{TaskId: in.Data.TaskId, Page: in.Data.Page, Limit: in.Data.Limit, Predecessor: in.Data.Predecessor})
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	data := &borrower.ProjectManagementDependencyPageResponseData{Page: page.Page, Limit: page.Limit, Total: page.Total, TotalPage: page.TotalPage, List: make([]*proto.ProjectManagementTaskDependencyEntity, 0, len(page.List))}
	for _, dependency := range page.List {
		entity, err := copierx.Trans[model.ProjectManagementTaskDependency, proto.ProjectManagementTaskDependencyEntity](dependency, copierx.TimeToRFC3339Str())
		if err != nil {
			return nil, errorx.GetFailedTransRespErr(err)
		}
		data.List = append(data.List, entity)
	}
	return &borrower.ProjectManagementDependencyPageResponse{Data: data}, nil
}

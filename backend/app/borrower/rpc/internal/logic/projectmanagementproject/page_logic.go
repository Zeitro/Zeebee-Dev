package projectmanagementprojectlogic

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
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PageLogic {
	return &PageLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *PageLogic) Page(in *borrower.ProjectManagementProjectPageRequest) (resp *borrower.ProjectManagementProjectPageResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	claim, err := l.baseLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	data, err := copierx.Trans[borrower.ProjectManagementProjectPageRequestData, dto.ProjectManagementProjectPageRequestDto](in.Data)
	if err != nil {
		return nil, err
	}
	page, err := l.baseLogic.ProjectManagementProjectLogic.Page(l.ctx, claim, *data)
	if err != nil {
		return nil, err
	}
	result := &borrower.ProjectManagementProjectPageResponseData{Page: page.Page, Limit: page.Limit, Total: page.Total, TotalPage: page.TotalPage, List: make([]*proto.ProjectManagementProjectEntity, 0, len(page.List))}
	for _, project := range page.List {
		entity, err := copierx.Trans[model.ProjectManagementProject, proto.ProjectManagementProjectEntity](project, copierx.TimeToRFC3339Str())
		if err != nil {
			return nil, err
		}
		result.List = append(result.List, entity)
	}
	return &borrower.ProjectManagementProjectPageResponse{Data: result}, nil
}

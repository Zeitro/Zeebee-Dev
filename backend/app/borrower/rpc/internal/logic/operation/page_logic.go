package operationlogic

import (
	"context"
	"encoding/json"
	"time"

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

func (l *PageLogic) Page(in *borrower.OperationPageRequest) (resp *borrower.OperationPageResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	base := logic.GetBaseLogic(l.svcCtx)
	claimValue, err := base.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	filter := dto.OperationPageFilterDto{}
	if in.Data.Filter != nil {
		filter.TaskId = in.Data.Filter.TaskId
		filter.OperationType = constants.OperationAction(in.Data.Filter.OperationType)
		filter.EntityType = constants.OperationEntityType(in.Data.Filter.EntityType)
		filter.EntityId = in.Data.Filter.EntityId
	}
	page, err := base.OperationLogic.Page(l.ctx, &claimValue, dto.OperationPageRequestDto{Page: in.Data.Page, Limit: in.Data.Limit, Filter: filter})
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	data := &borrower.OperationPageResponseData{Page: page.Page, Limit: page.Limit, Total: page.Total, TotalPage: page.TotalPage, List: make([]*proto.OperationEntity, 0, len(page.List))}
	for _, operation := range page.List {
		entity, err := operationEntity(operation)
		if err != nil {
			return nil, errorx.GetFailedTransRespErr(err)
		}
		data.List = append(data.List, entity)
	}
	return &borrower.OperationPageResponse{Data: data}, nil
}

func operationEntity(operation *model.Operation) (*proto.OperationEntity, error) {
	entity, err := copierx.Trans[model.Operation, proto.OperationEntity](operation)
	if err != nil {
		return nil, err
	}
	entity.EntityType = string(operation.EntityType)
	entity.EntityId = operation.EntityId
	if operation.EntityType == constants.OperationEntityTypeProjectManagementProject {
		entity.ProjectId = operation.EntityId
	}
	entity.CreatedAt, entity.UpdatedAt = formatEntityTime(operation.CreatedAt), formatEntityTime(operation.UpdatedAt)
	var metadata map[string]int64
	if json.Unmarshal(operation.Metadata, &metadata) == nil {
		entity.TaskId = metadata["task_id"]
	}
	return entity, nil
}

func formatEntityTime(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.Format(time.RFC3339)
}

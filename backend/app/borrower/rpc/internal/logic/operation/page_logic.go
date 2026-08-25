package operationlogic

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
		filterData, err := copierx.Trans[borrower.OperationPageRequestFilter, dto.OperationPageFilterDto](in.Data.Filter)
		if err != nil {
			return nil, err
		}
		filter = *filterData
		filter.OperationType = constants.OperationAction(in.Data.Filter.OperationType)
		filter.EntityType = constants.OperationEntityType(in.Data.Filter.EntityType)
	}
	page, err := base.OperationLogic.Page(l.ctx, &claimValue, dto.OperationPageRequestDto{Page: in.Data.Page, Limit: in.Data.Limit, Filter: filter})
	if err != nil {
		return nil, err
	}
	entities, err := copierx.Trans[[]*model.Operation, []*proto.OperationEntity](&page.List, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}
	if entities == nil {
		entities = &[]*proto.OperationEntity{}
	}
	for i, operation := range page.List {
		(*entities)[i].EntityType = string(operation.EntityType)
		(*entities)[i].Metadata = string(operation.Metadata)
	}
	data := &borrower.OperationPageResponseData{Page: page.Page, Limit: page.Limit, Total: page.Total, TotalPage: page.TotalPage, List: *entities}
	return &borrower.OperationPageResponse{Data: data}, nil
}

func operationEntity(operation *model.Operation) (*proto.OperationEntity, error) {
	entity, err := copierx.Trans[model.Operation, proto.OperationEntity](operation, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}
	entity.EntityType = string(operation.EntityType)
	entity.Metadata = string(operation.Metadata)
	return entity, nil
}

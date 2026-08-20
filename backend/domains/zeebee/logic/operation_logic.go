package logic

import (
	"context"

	"gorm.io/gorm"

	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
)

type OperationLogic struct {
	*BaseLogic
}

func NewOperationLogic(baseLogic *BaseLogic) *OperationLogic {
	return &OperationLogic{
		BaseLogic: baseLogic,
	}
}

func (l *OperationLogic) Page(ctx context.Context, claim *utils.JwtUser, input dto.OperationPageRequestDto) (gormx.PageResponse[*model.Operation], error) {
	if input.Filter.EntityType == "" {
		return gormx.PageResponse[*model.Operation]{}, errorx.New("entity type is required")
	}
	switch input.Filter.EntityType {
	case constants.OperationEntityTypeProjectManagementProject:
		if input.Filter.EntityId == 0 {
			return gormx.PageResponse[*model.Operation]{}, errorx.New("project id is required")
		}
		if _, err := l.projectManagementAccess().RequireProject(ctx, input.Filter.EntityId, claim.OrganizationId); err != nil {
			return gormx.PageResponse[*model.Operation]{}, err
		}
		if _, err := l.projectManagementAccess().RequireMember(ctx, input.Filter.EntityId, claim.OrganizationUserId); err != nil {
			return gormx.PageResponse[*model.Operation]{}, err
		}
	case constants.OperationEntityTypeProjectManagementOrganization:
		if input.Filter.EntityId == 0 {
			return gormx.PageResponse[*model.Operation]{}, errorx.New("entity id is required")
		}
	case constants.OperationEntityTypeBorrowerPropertyTransaction:
		if input.Filter.EntityId == 0 {
			return gormx.PageResponse[*model.Operation]{}, errorx.New("transaction id is required")
		}
		if _, err := l.repos.BorrowerPropertyTransactionRepo.GetOrganizationTransaction(ctx, input.Filter.EntityId, claim.OrganizationId); err != nil {
			return gormx.PageResponse[*model.Operation]{}, projectManagementNotFound(err, "transaction not found")
		}
	default:
		return gormx.PageResponse[*model.Operation]{}, errorx.New("unsupported operation entity type")
	}
	return l.repos.OperationRepo.PageOperations(ctx, claim.OrganizationId, input)
}

func (l *OperationLogic) CreateOperationLog(ctx context.Context, tx *gorm.DB, input dto.OperationCreateDto) (err error) {
	// operator, _ := l.repos.OrganizationUserRepo.FindOneByIdWithTx(ctx, tx, input.OperatorId)
	//
	// operatorName := ""
	// if operator != nil {
	// 	operatorName = operator.FullName
	// }

	operation := &model.Operation{
		EntityId:       input.EntityId,
		EntityType:     input.EntityType,
		OperatorId:     input.OperatorId,
		OperatorName:   input.OperatorName,
		OperatorType:   input.OperatorType,
		Operation:      input.Operation,
		Note:           input.Note,
		Metadata:       input.Metadata,
		OrganizationId: input.OrganizationId,
	}
	_, err = l.repos.OperationRepo.CreateWithTx(ctx, tx, operation, false)
	if err != nil {
		return errorx.Wrap(err, "Failed to create operation log")
	}

	return nil
}

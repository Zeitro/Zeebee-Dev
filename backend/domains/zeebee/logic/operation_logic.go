package logic

import (
	"context"

	"gorm.io/gorm"

	"zeebee-backend/common/errorx"
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

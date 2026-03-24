package logic

import (
	"context"

	"gorm.io/gorm"

	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
)

// LoanOperationLogic provides loan operation business logic
type LoanOperationLogic struct {
	*BaseLogic
}

// NewLoanOperationLogic creates a new loan operation logic instance
func NewLoanOperationLogic(baseLogic *BaseLogic) *LoanOperationLogic {
	return &LoanOperationLogic{
		BaseLogic: baseLogic,
	}
}

func (l *LoanOperationLogic) CreateLoanOperationLog(ctx context.Context, tx *gorm.DB, input dto.LoanOperationCreateDto) (err error) {

	operator, _ := l.repos.OrganizationUserRepo.FindOneByIdWithTx(ctx, tx, input.OperatorId)

	operatorName := ""
	if operator != nil {
		operatorName = operator.FullName
	}

	operation := &model.LoanOperation{
		LoanId:       input.LoanId,
		OperatorId:   input.OperatorId,
		OperatorName: operatorName,
		OperatorType: input.OperationType,
		Action:       input.Action,
		From:         input.From,
		To:           input.To,
		Metadata:     input.Metadata,
	}
	_, err = l.repos.LoanOperationRepo.CreateWithTx(ctx, tx, operation, false)
	if err != nil {
		return errorx.Wrap(err, "Failed to create operation log")
	}

	return nil
}

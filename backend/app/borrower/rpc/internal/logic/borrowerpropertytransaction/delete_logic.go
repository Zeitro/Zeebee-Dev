package borrowerpropertytransactionlogic

import (
	"context"

	"gorm.io/gorm"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
)

type DeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *DeleteLogic) Delete(in *borrower.BorrowerPropertyTransactionDeleteRequest) (resp *borrower.EmptyResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	existing, err := l.svcCtx.BorrowerPropertyTransactionRepo.GetOne(l.ctx, in.Data.Id)
	if err != nil {
		return nil, errorx.Wrap(err, "transaction not found")
	}
	if existing.OrganizationId != claim.OrganizationId {
		return nil, errorx.New("transaction does not belong to your organization")
	}

	// validate no loan binding
	isBound, err := l.svcCtx.BorrowerPropertyTransactionRepo.IsBoundToLoan(l.ctx, in.Data.Id)
	if err != nil {
		return nil, err
	}
	if isBound {
		return nil, errorx.New("cannot delete transaction bound to a loan")
	}

	err = l.svcCtx.BorrowerPropertyTransactionRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		// delete transaction
		_, innerErr := l.svcCtx.BorrowerPropertyTransactionRepo.DeleteByIDWithTx(l.ctx, tx, in.Data.Id)
		if innerErr != nil {
			return innerErr
		}

		// find and delete entity_document
		entityDoc, innerErr := l.svcCtx.EntityDocumentRepo.GetByEntityWithTx(l.ctx, tx, constants.EntityDocumentTypeBorrowerPropertyTransaction, in.Data.Id)
		if innerErr != nil {
			return innerErr
		}
		_, innerErr = l.svcCtx.EntityDocumentRepo.DeleteByIDWithTx(l.ctx, tx, entityDoc.Id)
		if innerErr != nil {
			return innerErr
		}

		// get all document ids and delete them
		var docIds []int64
		innerErr = l.svcCtx.DocumentRepo.PluckByColumnsWithTx(l.ctx, tx, "id", &docIds, map[string]interface{}{
			"entity_document_id": entityDoc.Id,
		})
		if innerErr != nil {
			return innerErr
		}

		_, innerErr = l.svcCtx.DocumentRepo.DeleteByIDsWithTx(l.ctx, tx, docIds)
		if innerErr != nil {
			return innerErr
		}
		// delete document files
		innerErr = l.svcCtx.FileRepo.DeleteByEntityIdWithTx(l.ctx, tx, docIds, constants.FileEntityTypeDocument)

		if innerErr != nil {
			return innerErr
		}

		// create operation log
		innerErr = l.baseLogic.OperationLogic.CreateOperationLog(l.ctx, tx, dto.OperationCreateDto{
			EntityId:       in.Data.Id,
			EntityType:     constants.OperationEntityTypeBorrowerPropertyTransaction,
			OperatorId:     claim.OrganizationUserId,
			OperatorName:   claim.Username,
			OperatorType:   claim.UserType,
			Operation:      constants.OperationActionBorrowerPropertyTransactionDelete,
			OrganizationId: claim.OrganizationId,
		})
		return innerErr
	})
	if err != nil {
		return nil, errorx.Wrap(err, "failed to delete transaction")
	}

	return &borrower.EmptyResponse{}, nil
}

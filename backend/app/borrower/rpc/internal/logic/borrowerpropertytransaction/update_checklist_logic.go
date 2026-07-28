package borrowerpropertytransactionlogic

import (
	"context"
	"encoding/json"

	"gorm.io/gorm"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/proto"

	"zeebee-backend/common/errorx"
)

type UpdateChecklistLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewUpdateChecklistLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateChecklistLogic {
	return &UpdateChecklistLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *UpdateChecklistLogic) UpdateChecklist(in *borrower.BorrowerPropertyTransactionChecklistRequest) (resp *proto.EmptyResponse, err error) {
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

	// check if status is changed
	if existing.Checklist.AllChecked == in.Data.AllChecked {
		return nil, nil
	}

	checklist := existing.Checklist
	oldChecklist := existing.Checklist
	checklist.AllChecked = in.Data.AllChecked

	checklistJson, _ := json.Marshal(checklist)
	err = l.svcCtx.BorrowerPropertyTransactionRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		innerErr := l.svcCtx.BorrowerPropertyTransactionRepo.UpdateColumnsById(l.ctx, in.Data.Id, map[string]interface{}{
			"checklist": string(checklistJson),
		})
		if innerErr != nil {
			return innerErr
		}

		metadata, _ := json.Marshal(map[string]interface{}{
			"before": map[string]interface{}{
				"checklist": oldChecklist,
			},
			"after": map[string]interface{}{
				"checklist": checklist,
			},
		})

		innerErr = l.baseLogic.OperationLogic.CreateOperationLog(l.ctx, tx, dto.OperationCreateDto{
			EntityId:       in.Data.Id,
			EntityType:     constants.OperationEntityTypeBorrowerPropertyTransaction,
			OperatorId:     claim.OrganizationUserId,
			OperatorName:   claim.Username,
			OperatorType:   claim.UserType,
			Operation:      constants.OperationActionBorrowerPropertyTransactionUpdateChecklist,
			OrganizationId: claim.OrganizationId,
			Metadata:       metadata,
		})
		return innerErr
	})
	if err != nil {
		return nil, errorx.Wrap(err, "failed to update checklist")
	}

	return &proto.EmptyResponse{}, nil
}

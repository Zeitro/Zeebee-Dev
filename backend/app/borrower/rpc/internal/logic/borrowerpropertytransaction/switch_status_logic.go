package borrowerpropertytransactionlogic

import (
	"context"
	"encoding/json"

	"gorm.io/gorm"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
)

type SwitchStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewSwitchStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SwitchStatusLogic {
	return &SwitchStatusLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *SwitchStatusLogic) SwitchStatus(in *borrower.BorrowerPropertyTransactionSwitchStatusRequest) (resp *borrower.EmptyResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	// validate ownership
	existing, err := l.svcCtx.BorrowerPropertyTransactionRepo.GetOne(l.ctx, in.Data.Id)
	if err != nil {
		return nil, errorx.Wrap(err, "transaction not found")
	}
	if existing.OrganizationId != claim.OrganizationId {
		return nil, errorx.New("transaction does not belong to your organization")
	}

	newStatus := constants.BorrowerPropertyTransactionStatus(in.Data.Status)
	if !newStatus.IsValid() {
		return nil, errorx.New("invalid status")
	}

	oldStatus := existing.Status

	if oldStatus == newStatus {
		return &borrower.EmptyResponse{}, nil
	}

	err = l.svcCtx.BorrowerPropertyTransactionRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		innerErr := l.svcCtx.BorrowerPropertyTransactionRepo.UpdateColumnsById(l.ctx, in.Data.Id, map[string]interface{}{
			"status": newStatus,
		})
		if innerErr != nil {
			return innerErr
		}

		metadata, _ := json.Marshal(map[string]interface{}{
			"before": map[string]interface{}{
				"status": oldStatus,
			},
			"after": map[string]interface{}{
				"status": newStatus,
			},
		})

		innerErr = l.baseLogic.OperationLogic.CreateOperationLog(l.ctx, tx, dto.OperationCreateDto{
			EntityId:       in.Data.Id,
			EntityType:     constants.OperationEntityTypeBorrowerPropertyTransaction,
			OperatorId:     claim.OrganizationUserId,
			OperatorName:   claim.Username,
			OperatorType:   claim.UserType,
			Operation:      constants.OperationActionBorrowerPropertyTransactionUpdateStatus,
			OrganizationId: claim.OrganizationId,
			Metadata:       metadata,
		})
		return innerErr
	})
	if err != nil {
		return nil, errorx.Wrap(err, "failed to switch status")
	}

	return &borrower.EmptyResponse{}, nil
}

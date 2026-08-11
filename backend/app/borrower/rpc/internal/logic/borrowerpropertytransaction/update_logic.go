package borrowerpropertytransactionlogic

import (
	"context"
	"encoding/json"
	"time"

	"gorm.io/gorm"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/timex"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/proto"
)

type UpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *UpdateLogic) Update(in *borrower.BorrowerPropertyTransactionUpdateRequest) (resp *borrower.BorrowerPropertyTransactionDetailResponse, err error) {
	claim, err := l.baseLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	// validate ownership
	existing, err := l.svcCtx.BorrowerPropertyTransactionRepo.GetOrganizationTransaction(l.ctx, in.Data.Id, claim.OrganizationId)
	if err != nil {
		return nil, errorx.Wrap(err, "transaction not found")
	}
	if existing.OrganizationId != claim.OrganizationId {
		return nil, errorx.New("transaction does not belong to your organization")
	}

	updateData := &model.BorrowerPropertyTransaction{}

	err = copierx.Assign(in.Data, updateData, copierx.StrToSqlTimePtr(timex.LayoutUsStyle))
	if err != nil {
		return nil, err
	}

	err = l.svcCtx.BorrowerPropertyTransactionRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		_, innerErr := l.svcCtx.BorrowerPropertyTransactionRepo.UpdateOne(l.ctx, in.Data.Id, updateData)
		if innerErr != nil {
			return innerErr
		}

		// update property
		propertyReq := in.Data.Property
		if propertyReq != nil {
			if existing.Property == nil {
				existing.Property = &model.PropertyTransactionProperty{}
			}
			if propertyReq.Id > 0 && propertyReq.Id != existing.Property.Id {
				return errorx.New("property does not belong to transaction")
			}
			if propertyReq.Id == 0 {
				propertyReq.Id = existing.Property.Id
			}

			innerErr = copierx.Assign[borrower.BorrowerPropertyTransactionPropertyEditEntity, model.PropertyTransactionProperty](propertyReq, existing.Property)
			if innerErr != nil {
				return innerErr
			}
			_, innerErr = l.svcCtx.PropertyTransactionPropertyRepo.UpdateOneByIdWithTx(l.ctx, tx, propertyReq.Id, existing.Property)
			if innerErr != nil {
				return innerErr
			}
		}

		// sync property transaction escrows (full-replace by id); conversions via copierx
		keepIds := make([]int64, 0, len(in.Data.Escrows))
		for _, e := range in.Data.Escrows {
			if e.Id > 0 {
				keepIds = append(keepIds, e.Id)
				existingEscrow, innerErr := l.svcCtx.PropertyTransactionEscrowRepo.FindOneByIdWithTx(l.ctx, tx, e.Id)
				if innerErr != nil {
					return innerErr
				}

				if existingEscrow.BorrowerPropertyTransactionId != in.Data.Id {
					return errorx.New("escrow does not belong to transaction")
				}

				// Assign onto the loaded row so timestamps are preserved.
				if innerErr = copierx.Assign(e, existingEscrow); innerErr != nil {
					return innerErr
				}
				if _, innerErr = l.svcCtx.PropertyTransactionEscrowRepo.UpdateOneByIdWithTx(l.ctx, tx, e.Id, existingEscrow); innerErr != nil {
					return innerErr
				}
			} else {
				esc, innerErr := copierx.Trans[borrower.PropertyTransactionEscrowData, model.PropertyTransactionEscrow](e)
				if innerErr != nil {
					return innerErr
				}
				esc.BorrowerPropertyTransactionId = in.Data.Id
				created, innerErr := l.svcCtx.PropertyTransactionEscrowRepo.CreateWithTx(l.ctx, tx, esc, false)
				if innerErr != nil {
					return innerErr
				}
				keepIds = append(keepIds, created.Id)
			}
		}
		if innerErr = l.svcCtx.PropertyTransactionEscrowRepo.DeleteByIdsNotInWithTx(l.ctx, tx, in.Data.Id, keepIds); innerErr != nil {
			return innerErr
		}

		type transactionMetadata struct {
			*model.BorrowerPropertyTransaction
			Id        int64             `json:"-"`
			Code      string            `json:"-"`
			Checklist dto.ChecklistRoot `json:"-"`
			CreatedAt time.Time         `json:"-"`
			UpdatedAt time.Time         `json:"-"`
		}

		metadata, _ := json.Marshal(map[string]interface{}{
			"before": transactionMetadata{BorrowerPropertyTransaction: existing},
			"after":  transactionMetadata{BorrowerPropertyTransaction: updateData},
		})

		innerErr = l.baseLogic.OperationLogic.CreateOperationLog(l.ctx, tx, dto.OperationCreateDto{
			EntityId:       in.Data.Id,
			EntityType:     constants.OperationEntityTypeBorrowerPropertyTransaction,
			OperatorId:     claim.OrganizationUserId,
			OperatorName:   claim.Username,
			OperatorType:   claim.UserType,
			Operation:      constants.OperationActionBorrowerPropertyTransactionUpdateBaseInfo,
			OrganizationId: claim.OrganizationId,
			Metadata:       metadata,
		})
		return innerErr
	})
	if err != nil {
		return nil, errorx.Wrap(err, "fail to update transaction")
	}

	detailLogic := NewDetailLogic(l.ctx, l.svcCtx)
	detail, err := detailLogic.Detail(&proto.PathIdRequest{
		Data: &proto.PathIdRequestData{
			Id: in.Data.Id,
		},
	})
	return detail, err
}

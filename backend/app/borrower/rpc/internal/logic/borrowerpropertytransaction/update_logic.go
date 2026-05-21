package borrowerpropertytransactionlogic

import (
	"context"
	"database/sql"
	"time"

	"gorm.io/gorm"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/timex"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
)

type UpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *UpdateLogic) Update(in *borrower.BorrowerPropertyTransactionUpdateRequest) (resp *borrower.EmptyResponse, err error) {
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
		return nil, errorx.Wrap(err, "Transaction not found")
	}
	if existing.OrganizationId != claim.OrganizationId {
		return nil, errorx.New("Transaction does not belong to your organization")
	}

	updateData := &model.BorrowerPropertyTransaction{
		PurchaseAmount: in.Data.PurchaseAmount,
		EmdAmount:      in.Data.EmdAmount,
		PropertyType:   constants.PropertyType(in.Data.PropertyType),
		Pm:             in.Data.Pm,
		Strategy:       in.Data.Strategy,
		Note:           in.Data.Note,
	}

	if in.Data.CoeDate != "" {
		t, parseErr := time.Parse(timex.LayoutUsStyle, in.Data.CoeDate)
		if parseErr == nil {
			updateData.CoeDate = &sql.NullTime{Time: t, Valid: true}
		}
	}

	err = l.svcCtx.BorrowerPropertyTransactionRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		_, innerErr := l.svcCtx.BorrowerPropertyTransactionRepo.UpdateOne(l.ctx, in.Data.Id, updateData)
		if innerErr != nil {
			return innerErr
		}

		innerErr = l.baseLogic.OperationLogic.CreateOperationLog(l.ctx, tx, dto.OperationCreateDto{
			EntityId:       in.Data.Id,
			EntityType:     constants.OperationEntityTypeBorrowerPropertyTransaction,
			OperatorId:     claim.OrganizationUserId,
			OperatorName:   claim.Username,
			Operation:      constants.OperationActionBorrowerPropertyTransactionUpdateBaseInfo,
			OrganizationId: claim.OrganizationId,
		})
		return innerErr
	})
	if err != nil {
		return nil, errorx.Wrap(err, "Failed to update transaction")
	}

	return &borrower.EmptyResponse{}, nil
}

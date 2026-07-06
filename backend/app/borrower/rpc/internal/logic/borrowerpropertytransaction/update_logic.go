package borrowerpropertytransactionlogic

import (
	"context"
	"encoding/json"
	"time"

	"gorm.io/gorm"

	"github.com/zeromicro/go-zero/core/logx"

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

func (l *UpdateLogic) Update(in *borrower.BorrowerPropertyTransactionUpdateRequest) (resp *borrower.BorrowerPropertyTransactionDetailResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
				logx.Field("service_name", l.svcCtx.Config.Name),
			).Error(err)
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
		return nil, errorx.Wrap(err, "failed to update transaction")
	}

	detailLogic := NewDetailLogic(l.ctx, l.svcCtx)
	detail, err := detailLogic.Detail(&proto.PathIdRequest{
		Data: &proto.PathIdRequestData{
			Id: in.Data.Id,
		},
	})
	return detail, err
}

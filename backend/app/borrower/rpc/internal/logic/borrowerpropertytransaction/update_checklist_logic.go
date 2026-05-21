package borrowerpropertytransactionlogic

import (
	"context"
	"encoding/json"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/domains/zeebee/data/dto"

	"github.com/zeromicro/go-zero/core/logx"
	"zeebee-backend/common/errorx"
)

type UpdateChecklistLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewUpdateChecklistLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateChecklistLogic {
	return &UpdateChecklistLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *UpdateChecklistLogic) UpdateChecklist(in *borrower.BorrowerPropertyTransactionChecklistRequest) (resp *borrower.EmptyResponse, err error) {
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
		return nil, errorx.Wrap(err, "Transaction not found")
	}
	if existing.OrganizationId != claim.OrganizationId {
		return nil, errorx.New("Transaction does not belong to your organization")
	}

	// update all_checked on existing checklist, then save
	var checklist dto.ChecklistRoot
	if existing.Checklist != nil {
		_ = json.Unmarshal(existing.Checklist, &checklist)
	}
	checklist.AllChecked = in.Data.AllChecked

	updatedChecklist, _ := json.Marshal(checklist)

	err = l.svcCtx.BorrowerPropertyTransactionRepo.UpdateColumnsById(l.ctx, in.Data.Id, map[string]interface{}{
		"checklist": string(updatedChecklist),
	})
	if err != nil {
		return nil, errorx.Wrap(err, "Failed to update checklist")
	}

	return &borrower.EmptyResponse{}, nil
}

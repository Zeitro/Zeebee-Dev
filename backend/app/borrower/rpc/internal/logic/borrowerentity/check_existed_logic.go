package borrowerentitylogic

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
)

type CheckExistedLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewCheckExistedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckExistedLogic {
	return &CheckExistedLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *CheckExistedLogic) CheckExisted(in *borrower.BorrowerEntityCheckExistedRequest) (resp *borrower.BorrowerEntityCheckExistedResponse, err error) {

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	condition := map[string]interface{}{
		// "user_id":         in.Data.UserId,
		"organization_id": claim.OrganizationId,
		"name":            in.Data.Name,
	}

	entity, err := l.svcCtx.BorrowerEntityRepo.FindOneByColumns(l.ctx, condition)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &borrower.BorrowerEntityCheckExistedResponse{
				Data: &borrower.BorrowerEntityCheckExistedResponseData{
					IsExisted: false,
				},
			}, nil
		}

		return nil, err
	}

	isExisted := true
	if in.Data.Id != 0 {
		isExisted = entity.Id != in.Data.Id
	}

	return &borrower.BorrowerEntityCheckExistedResponse{
		Data: &borrower.BorrowerEntityCheckExistedResponseData{
			IsExisted: isExisted,
		},
	}, nil
}

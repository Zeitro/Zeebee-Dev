package borrowerapplicantlogic

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

func (l *CheckExistedLogic) CheckExisted(in *borrower.BorrowerApplicantCheckExistedRequest) (resp *borrower.BorrowerApplicantCheckExistedResponse, err error) {

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	// todo check organization
	condition := map[string]interface{}{
		// "user_id":         in.Data.UserId,
		"organization_id": claim.OrganizationId,
		"first_name":      in.Data.FirstName,
		"mid_name":        in.Data.MidName,
		"last_name":       in.Data.LastName,
	}

	entity, err := l.svcCtx.BorrowerApplicantRepo.FindOneByColumns(l.ctx, condition)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &borrower.BorrowerApplicantCheckExistedResponse{
				Data: &borrower.BorrowerApplicantCheckExistedResponseData{
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

	return &borrower.BorrowerApplicantCheckExistedResponse{
		Data: &borrower.BorrowerApplicantCheckExistedResponseData{
			IsExisted: isExisted,
		},
	}, nil
}

package borrowerapplicantlogic

import (
	"context"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/errorx"
)

type CheckExistedLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewCheckExistedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckExistedLogic {
	return &CheckExistedLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *CheckExistedLogic) CheckExisted(in *borrower.BorrowerApplicantCheckExistedRequest) (resp *borrower.BorrowerApplicantCheckExistedResponse, err error) {
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

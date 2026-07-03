package borrowerentitylogic

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/zeromicro/go-zero/core/logx"

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

func (l *CheckExistedLogic) CheckExisted(in *borrower.BorrowerEntityCheckExistedRequest) (resp *borrower.BorrowerEntityCheckExistedResponse, err error) {
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

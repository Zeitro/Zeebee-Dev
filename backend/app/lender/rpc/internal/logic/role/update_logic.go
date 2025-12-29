package rolelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
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
		baseLogic: logic.NewBaseLogic(ctx, svcCtx),
	}
}

func (l *UpdateLogic) Update(in *lender.RoleUpdateRequest) (resp *lender.RoleUpdateResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	roleModel, err := copierx.Trans[lender.RoleUpdateRequestData, model.Role](in.Data)
	if err != nil {
		return nil, err
	}
	roleModel.Type = constants.UserTypeLender

	role, err := l.baseLogic.RoleLogic.Update(l.ctx, in.Data.Id, roleModel)
	if err != nil {
		return nil, err
	}

	data, err := copierx.Trans[model.Role, lender.RoleEntity](role, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}

	return &lender.RoleUpdateResponse{
		Data: data,
	}, nil
}

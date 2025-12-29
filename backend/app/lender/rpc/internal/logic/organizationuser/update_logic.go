package organizationuserlogic

import (
	"context"

	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
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

func (l *UpdateLogic) Update(in *lender.OrganizationUserUpdateRequest) (resp *lender.OrganizationUserUpdateResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	input := dto.OrganizationUserUpdateRequestDto{
		Id:             in.Data.Id,
		OrganizationId: claim.OrganizationId,
		Remark:         in.Data.Remark,
		RoleIds:        in.Data.RoleIds,
	}
	result, err := l.baseLogic.OrganizationUserLogic.Update(l.ctx, input)
	if err != nil {
		return nil, err
	}
	data, err := copierx.Trans[model.OrganizationUser, lender.OrganizationUserEntity](result, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	return &lender.OrganizationUserUpdateResponse{
		Data: data,
	}, nil
}

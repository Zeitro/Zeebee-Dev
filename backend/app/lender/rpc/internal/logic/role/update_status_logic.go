package rolelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
)

type UpdateStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewUpdateStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateStatusLogic {
	return &UpdateStatusLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *UpdateStatusLogic) UpdateStatus(in *lender.RoleUpdateStatusRequest) (resp *lender.EmptyResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	err = l.baseLogic.RoleLogic.SwitchStatus(l.ctx, in.Data.Id, claim.OrganizationId, constants.UserTypeLender, constants.CommonStatus(in.Data.Status))
	if err != nil {
		return nil, err
	}

	return &lender.EmptyResponse{}, nil
}

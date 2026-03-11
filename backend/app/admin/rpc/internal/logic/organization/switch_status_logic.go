package organizationlogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/data/dto"

	"github.com/zeromicro/go-zero/core/logx"
)

type SwitchStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	*logic.BaseLogic
}

func NewSwitchStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SwitchStatusLogic {
	return &SwitchStatusLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		BaseLogic: logic.GetBaseLogic(svcCtx),
	}
}

// SwitchStatus switches organization status between Active and Disabled.
// It updates both the organization status and all associated organization users' status
// in a single transaction for data consistency.
func (l *SwitchStatusLogic) SwitchStatus(in *admin.OrganizationSwitchStatusRequest) (resp *admin.EmptyResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	input, err := copierx.Trans[admin.OrganizationSwitchStatusRequestData, dto.OrganizationSwitchStatusRequestDto](in.Data)
	if err != nil {
		return nil, err
	}
	err = l.BaseLogic.OrganizationLogic.SwitchStatus(l.ctx, *input)
	if err != nil {
		return nil, err
	}
	return &admin.EmptyResponse{}, nil
}

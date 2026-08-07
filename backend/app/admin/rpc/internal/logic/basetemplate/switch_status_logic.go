package basetemplatelogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/proto"
)

type SwitchStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewSwitchStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SwitchStatusLogic {
	return &SwitchStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

// SwitchStatus switches the status of a base template with cascade to organization templates
func (l *SwitchStatusLogic) SwitchStatus(in *admin.SwitchBaseTemplateStatusRequest) (resp *proto.EmptyResponse, err error) {

	// Validate status
	newStatus := constants.CommonStatus(in.Data.Status)
	if err = newStatus.Validate(); err != nil {
		return nil, errorx.WithStack(err)
	}

	// Perform switch with cascade
	err = l.baseLogic.BaseTemplateLogic.SwitchStatus(l.ctx, in.Data.Id, newStatus)
	if err != nil {
		return nil, err
	}

	return &proto.EmptyResponse{}, nil
}

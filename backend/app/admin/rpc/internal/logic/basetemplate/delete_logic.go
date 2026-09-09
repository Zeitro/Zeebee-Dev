package basetemplatelogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/domains/zeebee/proto"
)

type DeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

// Delete deletes a base template by ID with cascade to organization templates
func (l *DeleteLogic) Delete(in *admin.DeleteBaseTemplateRequest) (resp *proto.EmptyResponse, err error) {

	// Perform delete with cascade
	err = l.baseLogic.BaseTemplateLogic.Delete(l.ctx, in.Data.Data)
	if err != nil {
		return nil, err
	}

	return &proto.EmptyResponse{}, nil
}

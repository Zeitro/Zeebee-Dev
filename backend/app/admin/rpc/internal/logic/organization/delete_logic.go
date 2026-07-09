package organizationlogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/common/errorx"
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

// Delete deletes an organization and all its associated data.
// Deletion order (respecting foreign key constraints):
// 1. organization_users (via DeleteByOrganizationIdWithTx)
// 2. roles (via GORM clause.Associations cascade delete)
// 3. organization
// Note: Users are NOT deleted as they may belong to other organizations
func (l *DeleteLogic) Delete(in *proto.IdRequest) (resp *proto.EmptyResponse, err error) {
	defer func() {
		if err != nil {
			l.svcCtx.Logger.WithContext(l.ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
			).Error(err)
		}
	}()
	err = l.baseLogic.OrganizationLogic.Delete(l.ctx, in.Data)
	return &proto.EmptyResponse{}, err
}

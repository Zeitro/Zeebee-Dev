package organizationlogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/errorx"
)

type DeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.NewBaseLogic(ctx, svcCtx),
	}
}

// Delete deletes an organization and all its associated data.
// Deletion order (respecting foreign key constraints):
// 1. organization_users (via DeleteByOrganizationIdWithTx)
// 2. roles (via GORM clause.Associations cascade delete)
// 3. organization
// Note: Users are NOT deleted as they may belong to other organizations
func (l *DeleteLogic) Delete(in *admin.IdRequest) (resp *admin.EmptyResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	err = l.baseLogic.OrganizationLogic.Delete(l.ctx, in.Data)
	return &admin.EmptyResponse{}, err
}

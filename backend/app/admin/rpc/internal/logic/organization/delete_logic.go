package organizationlogic

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	*logic.BaseLogic
}

func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		BaseLogic: logic.NewBaseLogic(ctx, svcCtx),
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

	// Verify organization exists
	organization, err := l.svcCtx.AllRepos.OrganizationRepo.GetOrganization(l.ctx, in.Data, constants.UserTypeLender, false)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	// Use transaction for atomic deletion
	err = l.svcCtx.AllRepos.UserRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		var innerErr error

		// 1. Delete organization_users (using our custom method)
		innerErr = l.svcCtx.AllRepos.OrganizationUserRepo.DeleteByOrganizationIdWithTx(l.ctx, tx, organization.Id)
		if innerErr != nil {
			return innerErr
		}

		// 2. Delete organization and cascade delete roles via GORM Select
		// This will soft delete the organization and its roles due to is_del field
		innerErr = tx.WithContext(l.ctx).
			Select(clause.Associations).
			Delete(&organization).
			Error
		if innerErr != nil {
			return innerErr
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &admin.EmptyResponse{}, nil
}

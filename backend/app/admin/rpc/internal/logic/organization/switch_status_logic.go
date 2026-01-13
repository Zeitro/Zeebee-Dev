package organizationlogic

import (
	"context"

	"gorm.io/gorm"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"

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
		BaseLogic: logic.NewBaseLogic(ctx, svcCtx),
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

	// Validate status
	newStatus := constants.OrganizationStatus(in.Data.Status)
	err = newStatus.Validate()
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	// Get existing organization
	organization, err := l.svcCtx.AllRepos.OrganizationRepo.GetOrganizationWithId(l.ctx, in.Data.Id, true)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	// Use transaction to update both organization and organization users
	err = l.svcCtx.AllRepos.UserRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		var innerErr error

		// Update organization status
		organization.Status = newStatus
		_, innerErr = l.svcCtx.AllRepos.OrganizationRepo.UpdateOneByIdWithTx(l.ctx, tx, organization.Id, organization)
		if innerErr != nil {
			return innerErr
		}

		// Update organization owner status
		owner := organization.Owner
		if owner == nil {
			return nil
		}
		if owner.Status == constants.UserStatusInactive && newStatus == constants.OrganizationStatusActive {
			owner.Status = constants.UserStatusActive
			innerErr = l.svcCtx.UserRepo.UpdateColumnsByIdWithTx(l.ctx, tx, owner.Id, map[string]interface{}{
				"status": constants.UserStatusActive,
			})
		}
		return innerErr
	})

	if err != nil {
		return nil, err
	}

	return &admin.EmptyResponse{}, nil
}

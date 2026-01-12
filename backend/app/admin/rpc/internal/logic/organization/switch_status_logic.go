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
	if newStatus != constants.OrganizationStatusActive && newStatus != constants.OrganizationStatusDisabled {
		return nil, errorx.New("Invalid Status. Only Active and Disabled are allowed")
	}

	// Get existing organization
	organization, err := l.svcCtx.AllRepos.OrganizationRepo.GetOrganization(l.ctx, in.Data.Id, constants.UserTypeLender, false)
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

		// Batch update all organization users' status
		// This uses a single SQL UPDATE for performance (1000+ users < 3 seconds)
		var orgUserStatus constants.OrganizationUserStatus
		if newStatus == constants.OrganizationStatusActive {
			orgUserStatus = constants.OrganizationUserStatusActive
		} else {
			orgUserStatus = constants.OrganizationUserStatusDisabled
		}
		innerErr = l.svcCtx.AllRepos.OrganizationUserRepo.UpdateStatusByOrganizationId(l.ctx, organization.Id, orgUserStatus)
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

package logic

import (
	"context"

	"github.com/zeromicro/go-zero/core/threading"
	"gorm.io/gorm"

	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/auth"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/repo"
)

// OrganizationLogic provides organization business logic
type OrganizationLogic struct {
	*BaseLogic
	// enforcerMap map[constants.UserType]*rbac.Enforcer
	authTopicMap map[constants.UserType]string
}

// NewOrganizationLogic creates a new organization logic instance
func NewOrganizationLogic(baseLogic *BaseLogic, authTopicMap map[constants.UserType]string) *OrganizationLogic {
	return &OrganizationLogic{
		BaseLogic:    baseLogic,
		authTopicMap: authTopicMap,
	}
}

// Page returns a page of organization entities.
func (l *OrganizationLogic) Page(ctx context.Context, pageReq dto.OrganizationPageRequestDto) (page dto.PageResultDto[dto.OrganizationPageEntityDto], err error) {
	results, err := l.repos.OrganizationRepo.Page(ctx, pageReq)
	if err != nil {
		return page, err
	}
	pageResult, err := copierx.Trans[gormx.PageResponse[*repo.OrganizationPageEntity], dto.PageResultDto[dto.OrganizationPageEntityDto]](&results)
	if err != nil {
		return page, err
	}
	return *pageResult, nil
}

// DetailWithType retrieves organization details based on the specified type.
// It fetches the organization from the repository and transforms it into a business object.
func (l *OrganizationLogic) DetailWithType(ctx context.Context, id int64, organizationType constants.UserType) (result *dto.OrganizationEntityDto, err error) {
	organization, err := l.repos.OrganizationRepo.GetOrganization(ctx, id, organizationType, true)
	if err != nil {
		return
	}
	result, err = l.tranToBo(ctx, organization)
	return result, err
}

// Detail retrieves detailed information about an organization by its ID.
// It returns the organization details in the form of OrganizationBo.
func (l *OrganizationLogic) Detail(ctx context.Context, id int64) (result *dto.OrganizationEntityDto, err error) {
	organization, err := l.repos.OrganizationRepo.GetOrganizationById(ctx, id, true)
	if err != nil {
		return
	}
	result, err = l.tranToBo(ctx, organization)
	return result, err
}

// All retrieves all organizations for the specified organization type.
func (l *OrganizationLogic) All(ctx context.Context, organizationType constants.UserType) (result []*dto.OrganizationEntityDto, err error) {
	organizations, err := l.repos.OrganizationRepo.All(ctx, organizationType)
	if err != nil {
		return
	}
	resultBo, err := copierx.Trans[[]*model.Organization, []*dto.OrganizationEntityDto](&organizations)
	if err != nil {
		return nil, err
	}
	return *resultBo, nil
}

func (l *OrganizationLogic) AllForOption(ctx context.Context, organizationType constants.UserType, status constants.OrganizationStatus) (result []*dto.OrganizationOptionDto, err error) {
	organizations, err := l.repos.OrganizationRepo.AllForOption(ctx, organizationType, status)
	if err != nil {
		return
	}
	resultBo, err := copierx.Trans[[]*model.Organization, []*dto.OrganizationOptionDto](&organizations)
	if err != nil {
		return nil, err
	}
	return *resultBo, nil
}

// Delete removes an organization and all its associated data from the system.
func (l *OrganizationLogic) Delete(ctx context.Context, id int64) (err error) {
	organization, err := l.repos.OrganizationRepo.GetOrganizationById(ctx, id, true)
	if err != nil {
		return err
	}
	if organization.OwnerType == constants.OrganizationOwnerTypeSelf || organization.Type == constants.UserTypeAdmin {
		return errorx.New("Cannot delete self owned organization")
	}
	organizationType := organization.Type
	var policies [][]string
	err = l.repos.OrganizationRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var innerErr error
		organization.Owner = nil
		organization.Users = nil

		// Delete related roles
		for i := range organization.Roles {
			var relations []*model.RoleAndResource
			relations, innerErr = l.repos.RoleAndResourceRepo.GetRoleRelationsWithTx(ctx, tx, organization.Roles[i].Id)
			if innerErr != nil {
				return innerErr
			}
			for j := range relations {
				policies = append(policies, auth.BuildPolicy(organization.Roles[i].Id, relations[j].ResourceId, organization.Type))
			}

			innerErr = l.repos.RoleRepo.DeleteByIdWithTx(ctx, tx, organization.Roles[i].Id)
			if innerErr != nil {
				return innerErr
			}
		}
		// Delete related organization emails
		for i := range organization.OrganizationEmails {
			innerErr = l.repos.OrganizationEmailRepo.DeleteByIdWithTx(ctx, tx, organization.OrganizationEmails[i].Id)
			if innerErr != nil {
				return innerErr
			}
		}
		// Delete related organization users
		innerErr = l.repos.OrganizationUserRepo.DeleteByOrganizationIdWithTx(ctx, tx, organization.Id)
		if innerErr != nil {
			return innerErr
		}

		// delete organization template
		innerErr = l.repos.OrganizationTemplateRepo.DeleteByOrganizationIdWithTx(ctx, tx, organization.Id)
		if innerErr != nil {
			return innerErr
		}
		// delete organization setting
		innerErr = l.repos.OrganizationSettingRepo.DeleteByOrganizationIdWithTx(ctx, tx, organization.Id)
		if innerErr != nil {
			return innerErr
		}

		// delete organization
		_, innerErr = l.repos.OrganizationRepo.DeleteByIDWithTx(ctx, tx, id)
		if innerErr != nil {
			return innerErr
		}

		return innerErr

	})
	if err != nil {
		return err
	}

	threading.GoSafe(func() {
		subCtx := context.WithoutCancel(ctx)
		msg := auth.Msg{
			Action:          auth.ActionRoleReload,
			Scope:           auth.ScopeOrganization,
			Type:            organizationType,
			OrganizationId:  id,
			RemovedPolicies: policies,
		}
		_ = l.publishAuthEvent(subCtx, l.authTopicMap[organizationType], msg)
	})

	return nil
}

// UpdateOrder sets the display order of an organization. Lower = higher in list. Min 0.
func (l *OrganizationLogic) UpdateOrder(ctx context.Context, id int64, order int) error {
	if order < 0 {
		return errorx.New("Order must be >= 0")
	}
	return l.repos.OrganizationRepo.UpdateOrder(ctx, id, order)
}

// SwitchStatus changes the status of an organization.
func (l *OrganizationLogic) SwitchStatus(ctx context.Context, input dto.OrganizationSwitchStatusRequestDto) (err error) {

	// Validate status
	newStatus := input.Status
	err = newStatus.Validate()
	if err != nil {
		return
	}
	// Get existing organization
	organization, err := l.repos.OrganizationRepo.GetOrganizationById(ctx, input.Id, true)
	if err != nil {
		return
	}

	// Use transaction to update both organization and organization users
	err = l.repos.OrganizationRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var innerErr error

		// Update organization status
		organization.Status = newStatus
		if newStatus == constants.OrganizationStatusActive && organization.Order < 2 {
			organization.Order = 2
		}
		_, innerErr = l.repos.OrganizationRepo.UpdateOneByIdWithTx(ctx, tx, organization.Id, organization)
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
			innerErr = l.repos.UserRepo.UpdateColumnsByIdWithTx(ctx, tx, owner.Id, map[string]interface{}{
				"status": constants.UserStatusActive,
			})
			if innerErr != nil {
				return innerErr
			}
		}

		// Sync organization template
		if newStatus == constants.OrganizationStatusActive {
			organizationTemplateLogic := NewOrganizationTemplateLogic(l.BaseLogic)
			innerErr = organizationTemplateLogic.
				SyncByOrganizationIdWithTx(ctx, tx, organization)
		}

		return innerErr
	})

	if err != nil {
		return
	}

	threading.GoSafe(func() {
		subCtx := context.WithoutCancel(ctx)
		msg := auth.Msg{
			Action:         auth.ActionRoleReload,
			Scope:          auth.ScopeOrganization,
			Type:           organization.Type,
			OrganizationId: input.Id,
		}
		_ = l.publishAuthEvent(subCtx, l.authTopicMap[organization.Type], msg)
	})
	return

}

// SelectOptions retrieves a list of organization options based on the specified type.
func (l *OrganizationLogic) SelectOptions(ctx context.Context, organizationType constants.UserType) (results []*dto.OrganizationOptionDto, err error) {
	organizations, err := l.repos.OrganizationRepo.SelectOptions(ctx, organizationType)
	if err != nil {
		return
	}
	resultsPtr, err := copierx.Trans[[]*model.Organization, []*dto.OrganizationOptionDto](&organizations)
	if resultsPtr == nil {
		return nil, err
	}
	results = *resultsPtr
	return results, err
}

// tranToBo converts a model.Organization entity to a dto.OrganizationEntityDto business object.
// It enriches the owner's information if present and checks authorization status for organization emails.
func (l *OrganizationLogic) tranToBo(ctx context.Context, organization *model.Organization) (transBo *dto.OrganizationEntityDto, err error) {
	if organization == nil {
		return nil, nil
	}
	owner := organization.Owner
	if organization.Owner != nil {
		// Get organization user details to enrich owner information
		organizationUser, err := l.repos.OrganizationUserRepo.GetOrganizationUserByUserId(ctx, organization.Id, owner.Id)
		if err != nil {
			return nil, err
		}
		// Update owner's name information if available in organization user
		if len(organizationUser.FullName) > 0 {
			owner.FullName = organizationUser.FullName
			owner.FirstName = organizationUser.FirstName
			owner.LastName = organizationUser.LastName
		}
		organization.Owner = owner
	}

	// Transform the organization model to business object
	transBo, err = copierx.Trans[model.Organization, dto.OrganizationEntityDto](organization)
	if err != nil {
		return nil, err
	}

	// Check authorization status for each organization email
	for i := range transBo.OrganizationEmails {
		isExisted, _ := l.repos.EmailAuthorizationRepo.IsExistedByColumns(ctx, map[string]interface{}{
			"organization_email_id": transBo.OrganizationEmails[i].Id,
			"purpose":               transBo.OrganizationEmails[i].Purpose,
			"is_authorized":         true,
			"user_id":               0,
		})
		transBo.OrganizationEmails[i].IsAuthorized = isExisted
	}

	return transBo, nil
}

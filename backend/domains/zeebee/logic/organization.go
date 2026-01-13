package logic

import (
	"context"

	"gorm.io/gorm"

	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/bo"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/repo"
)

type OrganizationLogic struct {
	repos *repo.AllRepos
}

// NewOrganizationLogic creates a new instance of OrganizationLogic.
// It initializes the logic with the provided repository dependencies.
func NewOrganizationLogic(repos *repo.AllRepos) *OrganizationLogic {
	return &OrganizationLogic{
		repos: repos,
	}
}

// Page
func (l *OrganizationLogic) Page(ctx context.Context, pageReq dto.OrganizationPageRequestDto) (page bo.PageResultBo[bo.OrganizationPageEntityBo], err error) {
	results, err := l.repos.OrganizationRepo.Page(ctx, pageReq)
	if err != nil {
		return page, err
	}
	pageResult, err := copierx.Trans[gormx.PageResponse[*repo.OrganizationPageEntity], bo.PageResultBo[bo.OrganizationPageEntityBo]](&results)
	if err != nil {
		return page, err
	}
	return *pageResult, nil
}

// DetailWithType retrieves organization details based on the specified type.
// It fetches the organization from the repository and transforms it into a business object.
func (l *OrganizationLogic) DetailWithType(ctx context.Context, id int64, organizationType constants.UserType) (result *bo.OrganizationBo, err error) {
	organization, err := l.repos.OrganizationRepo.GetOrganization(ctx, id, organizationType, true)
	if err != nil {
		return
	}
	result, err = l.tranToBo(ctx, organization)
	return result, err
}

// Detail retrieves detailed information about an organization by its ID.
// It returns the organization details in the form of OrganizationBo.
func (l *OrganizationLogic) Detail(ctx context.Context, id int64) (result *bo.OrganizationBo, err error) {
	organization, err := l.repos.OrganizationRepo.GetOrganizationById(ctx, id, true)
	if err != nil {
		return
	}
	result, err = l.tranToBo(ctx, organization)
	return result, err
}

// All
func (l *OrganizationLogic) All(ctx context.Context, organizationType constants.UserType) (result []*bo.OrganizationBo, err error) {
	organizations, err := l.repos.OrganizationRepo.All(ctx, organizationType)
	if err != nil {
		return
	}
	resultBo, err := copierx.Trans[[]*model.Organization, []*bo.OrganizationBo](&organizations)
	if err != nil {
		return nil, err
	}
	return *resultBo, nil

}

// Delete removes an organization and all its associated data from the system.
// It performs the following operations within a database transaction:
// 1. Retrieves the organization by ID
// 2. Removes all associated roles
// 3. Removes all associated organization emails
// 4. Removes all organization-user relationships
// 5. Finally deletes the organization itself
func (l *OrganizationLogic) Delete(ctx context.Context, id int64) (err error) {
	organization, err := l.repos.OrganizationRepo.GetOrganizationById(ctx, id, true)
	if err != nil {
		return err
	}
	if organization.OwnerType == constants.OrganizationOwnerTypeSelf || organization.Type == constants.UserTypeAdmin {
		return errorx.New("Cannot delete self-owned organization")
	}

	err = l.repos.OrganizationRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var innerErr error
		organization.Owner = nil
		organization.Users = nil
		// delete related roles
		for i := range organization.Roles {
			innerErr = l.repos.RoleRepo.DeleteByIdWithTx(ctx, tx, organization.Roles[i].Id)
			if innerErr != nil {
				return innerErr
			}
		}
		// delete related organization_emails
		for i := range organization.OrganizationEmails {
			innerErr = l.repos.OrganizationEmailRepo.DeleteByIdWithTx(ctx, tx, organization.OrganizationEmails[i].Id)
			if innerErr != nil {
				return innerErr
			}
		}
		// delete related organization_users
		innerErr = l.repos.OrganizationUserRepo.DeleteByOrganizationIdWithTx(ctx, tx, organization.Id)
		if innerErr != nil {
			return innerErr
		}
		_, innerErr = l.repos.OrganizationRepo.DeleteByIDWithTx(ctx, tx, id)

		return innerErr

	})
	if err != nil {
		return errorx.WithStack(err)
	}

	return nil
}

// tranToBo converts a model.Organization entity to a bo.OrganizationBo business object.
// It enriches the owner's information if present and checks authorization status for organization emails.
func (l *OrganizationLogic) tranToBo(ctx context.Context, organization *model.Organization) (transBo *bo.OrganizationBo, err error) {
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
	transBo, err = copierx.Trans[model.Organization, bo.OrganizationBo](organization)
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

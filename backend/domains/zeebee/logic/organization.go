package logic

import (
	"context"

	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/bo"
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

// Detail retrieves organization details by ID and type.
// It fetches the organization from the repository and transforms it to the business object format.
func (l *OrganizationLogic) Detail(ctx context.Context, id int64, organizationType constants.UserType) (result *bo.OrganizationBo, err error) {
	organization, err := l.repos.OrganizationRepo.GetOrganization(ctx, id, organizationType, true)
	if err != nil {
		return
	}
	result, err = copierx.Trans[model.Organization, bo.OrganizationBo](organization)
	if err != nil {
		return
	}

	for i := range result.OrganizationEmails {
		isExisted, _ := l.repos.EmailAuthorizationRepo.IsExistedByColumns(ctx, map[string]interface{}{
			"organization_email_id": result.OrganizationEmails[i].Id,
			"purpose":               result.OrganizationEmails[i].Purpose,
			"is_authorized":         true,
			"user_id":               0,
		})
		result.OrganizationEmails[i].IsAuthorized = isExisted
	}
	return result, nil
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

package logic

import "zeebee-backend/domains/zeebee/repo"

type OrganizationSettingLogic struct {
	repos *repo.AllRepos
}

// NewOrganizationSettingLogic creates a new instance of OrganizationSettingLogic.
// It initializes the logic with the provided repository dependencies.
func NewOrganizationSettingLogic(repos *repo.AllRepos) *OrganizationSettingLogic {
	return &OrganizationSettingLogic{
		repos: repos,
	}
}

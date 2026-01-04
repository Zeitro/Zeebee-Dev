package logic

import (
	"zeebee-backend/common/email"
	"zeebee-backend/common/redisx"
	"zeebee-backend/domains/zeebee/repo"
)

type AllLogics struct {
	ApiLogic               *ApiLogic
	ResourceLogic          *ResourceLogic
	RoleLogic              *RoleLogic
	UserLogic              *UserLogic
	OrganizationUserLogic  *OrganizationUserLogic
	OrganizationLogic      *OrganizationLogic
	OrganizationEmailLogic *OrganizationEmailLogic
}

func NewAllLogics(repos *repo.AllRepos, redisCli *redisx.Redisx, jwtConfig JwtConfig, emailClient *email.Email) *AllLogics {
	return &AllLogics{
		ApiLogic:               NewApiLogic(repos),
		ResourceLogic:          NewResourceLogic(repos),
		RoleLogic:              NewRoleLogic(repos),
		UserLogic:              NewUserLogic(repos, redisCli, jwtConfig),
		OrganizationUserLogic:  NewOrganizationUserLogic(repos, redisCli, emailClient),
		OrganizationLogic:      NewOrganizationLogic(repos),
		OrganizationEmailLogic: NewOrganizationEmailLogic(repos),
	}
}

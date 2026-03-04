package logic

import (
	"zeebee-backend/common/email"
	"zeebee-backend/common/rbac"
	"zeebee-backend/common/redisx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/repo"
)

type AllLogics struct {
	ApiLogic                  *ApiLogic
	ResourceLogic             *ResourceLogic
	RoleLogic                 *RoleLogic
	UserLogic                 *UserLogic
	OrganizationUserLogic     *OrganizationUserLogic
	OrganizationLogic         *OrganizationLogic
	OrganizationEmailLogic    *OrganizationEmailLogic
	BaseTemplateLogic         *BaseTemplateLogic
	OrganizationTemplateLogic *OrganizationTemplateLogic
	OrganizationSettingLogic  *OrganizationSettingLogic
	LoanLogic                 *LoanLogic
	LoanOperationLogic        *LoanOperationLogic
	AuthLogic                 *AuthLogic
}

func NewAllLogics(
	repos *repo.AllRepos,
	redisCli *redisx.Redisx,
	jwtConfig JwtConfig,
	emailClient *email.Email,
	currentEnforcer *rbac.Enforcer,
	enforcerMap map[constants.UserType]*rbac.Enforcer,
) *AllLogics {
	return &AllLogics{
		ApiLogic:                  NewApiLogic(repos),
		ResourceLogic:             NewResourceLogic(repos),
		RoleLogic:                 NewRoleLogic(repos, currentEnforcer),
		UserLogic:                 NewUserLogic(repos, redisCli, jwtConfig),
		OrganizationUserLogic:     NewOrganizationUserLogic(repos, redisCli, emailClient),
		OrganizationLogic:         NewOrganizationLogic(repos, enforcerMap),
		OrganizationEmailLogic:    NewOrganizationEmailLogic(repos),
		BaseTemplateLogic:         NewBaseTemplateLogic(repos),
		OrganizationTemplateLogic: NewOrganizationTemplateLogic(repos),
		OrganizationSettingLogic:  NewOrganizationSettingLogic(repos),
		LoanLogic:                 NewLoanLogic(repos),
		LoanOperationLogic:        NewLoanOperationLogic(repos),
		AuthLogic:                 NewAuthLogic(repos, redisCli, jwtConfig),
	}
}

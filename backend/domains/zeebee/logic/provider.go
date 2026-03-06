package logic

import (
	"zeebee-backend/common/email"
	"zeebee-backend/common/rbac"
	"zeebee-backend/common/redisx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
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

type AllLogicConfig struct {
	Repos                 *repo.AllRepos
	RedisCli              *redisx.Redisx
	JwtConfig             dto.JwtConfig
	EmailClient           *email.Email
	CurrentEnforcer       *rbac.Enforcer
	EnforcerMap           map[constants.UserType]*rbac.Enforcer
	CurrentAuthEventTopic string
	AuthEventTopicMap     map[constants.UserType]string
}

func NewAllLogics(config AllLogicConfig) *AllLogics {
	return &AllLogics{
		ApiLogic:                  NewApiLogic(config.Repos, config.RedisCli),
		ResourceLogic:             NewResourceLogic(config.Repos),
		RoleLogic:                 NewRoleLogic(config.Repos, config.CurrentEnforcer),
		UserLogic:                 NewUserLogic(config.Repos, config.RedisCli, config.JwtConfig),
		OrganizationUserLogic:     NewOrganizationUserLogic(config.Repos, config.RedisCli, config.EmailClient),
		OrganizationLogic:         NewOrganizationLogic(config.Repos, config.EnforcerMap),
		OrganizationEmailLogic:    NewOrganizationEmailLogic(config.Repos),
		BaseTemplateLogic:         NewBaseTemplateLogic(config.Repos),
		OrganizationTemplateLogic: NewOrganizationTemplateLogic(config.Repos),
		OrganizationSettingLogic:  NewOrganizationSettingLogic(config.Repos),
		LoanLogic:                 NewLoanLogic(config.Repos),
		LoanOperationLogic:        NewLoanOperationLogic(config.Repos),
		AuthLogic:                 NewAuthLogic(config.Repos, config.RedisCli, config.JwtConfig, config.CurrentEnforcer, config.CurrentAuthEventTopic),
	}
}

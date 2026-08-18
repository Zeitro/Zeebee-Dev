package logic

import (
	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/email"
	"zeebee-backend/common/rbac"
	"zeebee-backend/common/redisx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/repo"
)

type AllLogics struct {
	ApiLogic                      *ApiLogic
	ResourceLogic                 *ResourceLogic
	RoleLogic                     *RoleLogic
	UserLogic                     *UserLogic
	OrganizationUserLogic         *OrganizationUserLogic
	OrganizationLogic             *OrganizationLogic
	OrganizationEmailLogic        *OrganizationEmailLogic
	DriveLogic                    *DriveLogic
	BaseTemplateLogic             *BaseTemplateLogic
	OrganizationTemplateLogic     *OrganizationTemplateLogic
	OrganizationSettingLogic      *OrganizationSettingLogic
	LoanLogic                     *LoanLogic
	LoanOperationLogic            *LoanOperationLogic
	OperationLogic                *OperationLogic
	AuthLogic                     *AuthLogic
	PricingPlanLogic              *PricingPlanLogic
	PricingPlanPolicyLogic        *PricingPlanPolicyLogic
	OrganizationSubscriptionLogic *OrganizationSubscriptionLogic
	ProjectManagementLogic        *ProjectManagementLogic
}

type AllLogicConfig struct {
	Repos                  *repo.AllRepos
	RedisCli               *redisx.Client
	JwtConfig              dto.JwtConfig
	EmailClient            *email.Email
	EmailClients           map[constants.UserType]*email.Email
	CurrentEnforcer        *rbac.Enforcer
	CurrentAuthEventTopics []string
	AuthEventTopicMap      map[constants.UserType]string
	Logger                 logx.Logger
}

func NewAllLogics(config AllLogicConfig) *AllLogics {
	baseLogic := NewBaseLogic(config.Repos, config.RedisCli, config.Logger, config.JwtConfig)
	return &AllLogics{
		ApiLogic:                      NewApiLogic(baseLogic, config.AuthEventTopicMap),
		ResourceLogic:                 NewResourceLogic(baseLogic, config.AuthEventTopicMap),
		RoleLogic:                     NewRoleLogic(baseLogic, config.AuthEventTopicMap),
		UserLogic:                     NewUserLogic(baseLogic, config.AuthEventTopicMap),
		OrganizationUserLogic:         NewOrganizationUserLogic(baseLogic, config.EmailClient, config.EmailClients, config.AuthEventTopicMap),
		OrganizationLogic:             NewOrganizationLogic(baseLogic, config.AuthEventTopicMap),
		OrganizationEmailLogic:        NewOrganizationEmailLogic(baseLogic),
		DriveLogic:                    NewDriveLogic(baseLogic),
		BaseTemplateLogic:             NewBaseTemplateLogic(baseLogic),
		OrganizationTemplateLogic:     NewOrganizationTemplateLogic(baseLogic),
		OrganizationSettingLogic:      NewOrganizationSettingLogic(baseLogic),
		LoanLogic:                     NewLoanLogic(baseLogic),
		LoanOperationLogic:            NewLoanOperationLogic(baseLogic),
		OperationLogic:                NewOperationLogic(baseLogic),
		AuthLogic:                     NewAuthLogic(baseLogic, config.CurrentEnforcer, config.CurrentAuthEventTopics, config.AuthEventTopicMap),
		PricingPlanLogic:              NewPricingPlanLogic(baseLogic, config.AuthEventTopicMap),
		PricingPlanPolicyLogic:        NewPricingPlanPolicyLogic(baseLogic),
		OrganizationSubscriptionLogic: NewOrganizationSubscriptionLogic(baseLogic, config.AuthEventTopicMap),
		ProjectManagementLogic:        NewProjectManagementLogic(baseLogic),
	}
}

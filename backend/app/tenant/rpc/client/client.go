// Package client bundles tenant's per-domain RPC clients (user, auth,
// organization, role, resource, ...) into one composite so each consuming
// backend wires a single tenantclient.New(cli) and reaches any sub-client via
// .User, .Auth, .Organization, etc. The sub-clients are the generated
// interfaces under app/tenant/rpc/client/<domain>; this package only composes
// them — it re-implements nothing.
package client

import (
	"github.com/zeromicro/go-zero/zrpc"

	"zeebee-backend/app/tenant/rpc/client/api"
	"zeebee-backend/app/tenant/rpc/client/auth"
	"zeebee-backend/app/tenant/rpc/client/organization"
	"zeebee-backend/app/tenant/rpc/client/organizationemail"
	"zeebee-backend/app/tenant/rpc/client/organizationsetting"
	"zeebee-backend/app/tenant/rpc/client/organizationsubscription"
	"zeebee-backend/app/tenant/rpc/client/organizationtemplate"
	"zeebee-backend/app/tenant/rpc/client/organizationuser"
	"zeebee-backend/app/tenant/rpc/client/pricingplan"
	"zeebee-backend/app/tenant/rpc/client/pricingplanpolicy"
	"zeebee-backend/app/tenant/rpc/client/resource"
	"zeebee-backend/app/tenant/rpc/client/role"
	"zeebee-backend/app/tenant/rpc/client/user"
)

// Client is the composite tenant RPC client. Each field is the generated
// interface for one tenant domain; all are constructed from the same
// zrpc.Client by New. Field name mirrors the interface type (cf. repo.AllRepos).
type Client struct {
	Api                      api.Api
	Auth                     auth.Auth
	Organization             organization.Organization
	OrganizationEmail        organizationemail.OrganizationEmail
	OrganizationSetting      organizationsetting.OrganizationSetting
	OrganizationSubscription organizationsubscription.OrganizationSubscription
	OrganizationTemplate     organizationtemplate.OrganizationTemplate
	OrganizationUser         organizationuser.OrganizationUser
	PricingPlan              pricingplan.PricingPlan
	PricingPlanPolicy        pricingplanpolicy.PricingPlanPolicy
	Resource                 resource.Resource
	Role                     role.Role
	User                     user.User
}

// New builds a Client holding every tenant sub-client from a single
// zrpc.Client. Sub-clients are thin wrappers that defer dialing, so New is safe
// to call before the tenant server is reachable.
func New(cli zrpc.Client) *Client {
	return &Client{
		Api:                      api.NewApi(cli),
		Auth:                     auth.NewAuth(cli),
		Organization:             organization.NewOrganization(cli),
		OrganizationEmail:        organizationemail.NewOrganizationEmail(cli),
		OrganizationSetting:      organizationsetting.NewOrganizationSetting(cli),
		OrganizationSubscription: organizationsubscription.NewOrganizationSubscription(cli),
		OrganizationTemplate:     organizationtemplate.NewOrganizationTemplate(cli),
		OrganizationUser:         organizationuser.NewOrganizationUser(cli),
		PricingPlan:              pricingplan.NewPricingPlan(cli),
		PricingPlanPolicy:        pricingplanpolicy.NewPricingPlanPolicy(cli),
		Resource:                 resource.NewResource(cli),
		Role:                     role.NewRole(cli),
		User:                     user.NewUser(cli),
	}
}

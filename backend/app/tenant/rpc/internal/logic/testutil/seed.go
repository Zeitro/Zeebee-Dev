package testutil

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
)

// SeedUser inserts a user and returns it with the generated id.
func SeedUser(t *testing.T, db *gorm.DB, u *model.User) *model.User {
	t.Helper()
	require.NoError(t, db.WithContext(context.Background()).Create(u).Error)
	return u
}

// SeedOrganization inserts an organization and returns it with the generated id.
func SeedOrganization(t *testing.T, db *gorm.DB, o *model.Organization) *model.Organization {
	t.Helper()
	require.NoError(t, db.WithContext(context.Background()).Create(o).Error)
	return o
}

// SeedOrganizationUser inserts an org-user binding.
func SeedOrganizationUser(t *testing.T, db *gorm.DB, ou *model.OrganizationUser) *model.OrganizationUser {
	t.Helper()
	require.NoError(t, db.WithContext(context.Background()).Create(ou).Error)
	return ou
}

// SeedRole inserts a role.
func SeedRole(t *testing.T, db *gorm.DB, r *model.Role) *model.Role {
	t.Helper()
	require.NoError(t, db.WithContext(context.Background()).Create(r).Error)
	return r
}

// SeedResource inserts a resource.
func SeedResource(t *testing.T, db *gorm.DB, r *model.Resource) *model.Resource {
	t.Helper()
	require.NoError(t, db.WithContext(context.Background()).Create(r).Error)
	return r
}

// SeedRoleResource binds a resource to a role.
func SeedRoleResource(t *testing.T, db *gorm.DB, roleId, resourceId int64) {
	t.Helper()
	require.NoError(t, db.WithContext(context.Background()).Create(&model.RoleAndResource{
		RoleId:     roleId,
		ResourceId: resourceId,
		DataScope:  constants.DataScopeAll,
	}).Error)
}

// SeedApi inserts an api registry entry.
func SeedApi(t *testing.T, db *gorm.DB, a *model.Api) *model.Api {
	t.Helper()
	require.NoError(t, db.WithContext(context.Background()).Create(a).Error)
	return a
}

// SeedPricingPlan inserts a pricing plan.
func SeedPricingPlan(t *testing.T, db *gorm.DB, p *model.PricingPlan) *model.PricingPlan {
	t.Helper()
	require.NoError(t, db.WithContext(context.Background()).Create(p).Error)
	return p
}

// SeedSubscription inserts an organization subscription.
func SeedSubscription(t *testing.T, db *gorm.DB, s *model.OrganizationSubscription) *model.OrganizationSubscription {
	t.Helper()
	require.NoError(t, db.WithContext(context.Background()).Create(s).Error)
	return s
}

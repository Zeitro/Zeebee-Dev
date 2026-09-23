package projectmanagementworkflowlogic

import (
	"context"
	"fmt"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"gorm.io/gorm"

	"zeebee-backend/app/borrower/rpc/internal/svc"
	tenantclient "zeebee-backend/app/tenant/rpc/client"
	tenantuser "zeebee-backend/app/tenant/rpc/client/user"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/proto"
	"zeebee-backend/domains/zeebee/repo"
)

type workflowSourceTenantUser struct{ tenantuser.User }

func (workflowSourceTenantUser) GetClaim(context.Context, *tenantuser.GetClaimRequest, ...grpc.CallOption) (*tenantuser.GetClaimResponse, error) {
	return &tenant.GetClaimResponse{Data: &tenant.GetClaimResponseData{
		UserId: 1, OrganizationId: 7, OrganizationUserId: 11, Username: "Owner", UserType: string(constants.UserTypeBorrower),
	}}, nil
}

func TestAllScopesWorkflowSourcesToClaimOrganization(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())), &gorm.Config{})
	require.NoError(t, err)
	for _, statement := range []string{
		"CREATE TABLE project_management_templates (id integer primary key, created_at datetime, updated_at datetime, is_del integer not null default 0, template_type text not null, organization_id integer not null default 0, name text not null default '', description text not null default '', status text not null default 'Active', is_default boolean not null default false, sort_order integer not null default 0, project_type text)",
		"CREATE TABLE project_management_workflows (id integer primary key, created_at datetime, updated_at datetime, is_del integer not null default 0, entity_type text not null, entity_id integer not null, source_workflow_id integer not null default 0, source_workflow_name text not null default '', name text not null default '', description text not null default '', is_default boolean not null default false)",
		"CREATE TABLE project_management_stages (id integer primary key, created_at datetime, updated_at datetime, is_del integer not null default 0, entity_type text not null, entity_id integer not null, name text not null default '', description text not null default '', phase text not null, color text not null default '', sort_order integer not null default 0)",
		"CREATE TABLE project_management_workflow_stages (id integer primary key, workflow_id integer not null, stage_id integer not null)",
		"CREATE TABLE project_management_workflow_transitions (id integer primary key, created_at datetime, updated_at datetime, is_del integer not null default 0, workflow_id integer not null, name text not null default '', from_stage_id integer not null, to_stage_id integer not null, allowed_actor_types json not null default '[]')",
		"INSERT INTO project_management_templates (id, template_type, organization_id, name, status, is_default) VALUES (1, 'Organization', 7, 'Default', 'Active', true), (2, 'Organization', 8, 'Foreign', 'Active', true)",
		"INSERT INTO project_management_workflows (id, entity_type, entity_id, name, is_default) VALUES (1, 'Template', 1, 'Default', true), (2, 'Template', 2, 'Foreign', true)",
	} {
		require.NoError(t, db.Exec(statement).Error)
	}
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(utils.JwtTokenCtxKey, "tenant-issued-token", utils.UserTypeCtxKey, string(constants.UserTypeBorrower)))
	svcCtx := &svc.ServiceContext{AllRepos: repo.NewAllRepos(db), TenantClient: &tenantclient.Client{User: workflowSourceTenantUser{}}}
	resp, err := NewAllLogic(ctx, svcCtx).All(&proto.EmptyRequest{})
	require.NoError(t, err)
	require.Len(t, resp.Data, 1)
	require.Equal(t, int64(1), resp.Data[0].Id)
}

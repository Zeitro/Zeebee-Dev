package borrowerpropertylogic

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
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	tenantclient "zeebee-backend/app/tenant/rpc/client"
	tenantuser "zeebee-backend/app/tenant/rpc/client/user"
	tenant "zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/repo"
)

type propertyDeleteTenantUser struct {
	tenantuser.User
}

func (propertyDeleteTenantUser) GetClaim(context.Context, *tenantuser.GetClaimRequest, ...grpc.CallOption) (*tenantuser.GetClaimResponse, error) {
	return &tenant.GetClaimResponse{Data: &tenant.GetClaimResponseData{
		UserId:             1,
		OrganizationId:     1,
		OrganizationUserId: 1,
		Username:           "Owner",
		UserType:           string(constants.UserTypeBorrower),
	}}, nil
}

func TestDeleteRejectsPropertyWithUnfinishedProjectTask(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())), &gorm.Config{})
	require.NoError(t, err)
	for _, statement := range []string{
		"CREATE TABLE borrower_properties (id integer primary key, is_del integer not null default 0, organization_id integer not null default 0, loan_count integer not null default 0)",
		"CREATE TABLE project_management_projects (id integer primary key, is_del integer not null default 0, entity_id integer not null default 0, entity_type text not null default '')",
		"CREATE TABLE project_management_workflow_stages (id integer primary key, is_del integer not null default 0, phase text not null default '')",
		"CREATE TABLE project_management_tasks (id integer primary key, is_del integer not null default 0, project_id integer not null default 0, stage_id integer not null default 0)",
		"INSERT INTO borrower_properties (id, organization_id) VALUES (1, 1)",
		"INSERT INTO project_management_projects (id, entity_id, entity_type) VALUES (1, 1, 'Property')",
		"INSERT INTO project_management_workflow_stages (id, phase) VALUES (1, 'InProgress')",
		"INSERT INTO project_management_tasks (id, project_id, stage_id) VALUES (1, 1, 1)",
	} {
		require.NoError(t, db.Exec(statement).Error)
	}

	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(
		utils.JwtTokenCtxKey, "tenant-issued-token",
		utils.UserTypeCtxKey, string(constants.UserTypeBorrower),
	))
	svcCtx := &svc.ServiceContext{
		AllRepos: repo.NewAllRepos(db),
		TenantClient: &tenantclient.Client{
			User: propertyDeleteTenantUser{},
		},
	}

	_, err = NewDeleteLogic(ctx, svcCtx).Delete(&borrower.PropertyDeleteRequest{Data: &borrower.PropertyDeleteRequestData{Id: 1}})

	require.ErrorContains(t, err, "cannot delete a property with unfinished project tasks")
	var property model.BorrowerProperty
	require.NoError(t, db.First(&property, 1).Error)
}

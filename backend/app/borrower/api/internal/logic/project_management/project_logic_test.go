package project_management

import (
	"testing"

	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
)

func TestProjectConvertersKeepPropertyOwnerAndMemberRole(t *testing.T) {
	data, err := copierx.Trans[types.ProjectManagementProjectCreateRequest, borrower.ProjectManagementProjectCreateRequestData](&types.ProjectManagementProjectCreateRequest{
		PropertyId: 7, Name: "Renovation", Description: "Kitchen", CoverUrl: "cover", WorkflowTemplateId: 9,
	})
	if err != nil {
		t.Fatal(err)
	}
	if data.PropertyId != 7 || data.WorkflowTemplateId != 9 || data.Name != "Renovation" {
		t.Fatalf("project request lost data: %#v", data)
	}

	member, err := copierx.Trans[borrower.ProjectManagementMemberEntity, types.ProjectManagementMemberEntity](&borrower.ProjectManagementMemberEntity{
		ProjectId: 1, OrganizationUserId: 10, Role: "ProjectOwner",
	})
	if err != nil {
		t.Fatal(err)
	}
	if member.ProjectId != 1 || member.OrganizationUserId != 10 || member.Role != "ProjectOwner" {
		t.Fatalf("member role lost: %#v", member)
	}
}

package projectmanagementlogic

import (
	"testing"

	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
)

func TestProjectConvertersKeepEntityOwnerAndMemberRole(t *testing.T) {
	data, err := copierx.Trans[borrower.ProjectManagementProjectCreateRequestData, dto.ProjectManagementProjectDto](&borrower.ProjectManagementProjectCreateRequestData{
		EntityId: 7, EntityType: "Transaction", Name: "Renovation", Description: "Kitchen", CoverUrl: "cover", WorkflowTemplateId: 9,
	})
	if err != nil {
		t.Fatal(err)
	}
	if data.EntityId != 7 || data.EntityType != "Transaction" || data.Name != "Renovation" || data.CoverUrl != "cover" {
		t.Fatalf("project create data lost: %#v", data)
	}

	entity, err := copierx.Trans[model.ProjectManagementProject, borrower.ProjectManagementProjectEntity](&model.ProjectManagementProject{
		Base: model.Base{Id: 1}, EntityId: 7, EntityType: "Transaction", Name: "Renovation", OwnerId: 10, OwnerName: "Owner",
	})
	if err != nil {
		t.Fatal(err)
	}
	if entity.Id != 1 || entity.EntityId != 7 || entity.EntityType != "Transaction" || entity.OwnerId != 10 || entity.OwnerName != "Owner" {
		t.Fatalf("project response lost identity: %#v", entity)
	}
}

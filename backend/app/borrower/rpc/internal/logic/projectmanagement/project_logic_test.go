package projectmanagementlogic

import (
	"testing"

	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
)

func TestProjectConvertersKeepPropertyOwnerAndMemberRole(t *testing.T) {
	data, err := copierx.Trans[borrower.ProjectManagementProjectCreateRequestData, dto.ProjectManagementProjectDto](&borrower.ProjectManagementProjectCreateRequestData{
		PropertyId: 7, Name: "Renovation", Description: "Kitchen", CoverUrl: "cover", WorkflowTemplateId: 9,
	})
	if err != nil {
		t.Fatal(err)
	}
	if data.PropertyId != 7 || data.Name != "Renovation" || data.CoverUrl != "cover" {
		t.Fatalf("project create data lost: %#v", data)
	}

	entity, err := copierx.Trans[model.ProjectManagementProject, borrower.ProjectManagementProjectEntity](&model.ProjectManagementProject{
		Base: model.Base{Id: 1}, PropertyId: 7, Name: "Renovation", OwnerId: 10,
	})
	if err != nil {
		t.Fatal(err)
	}
	if entity.Id != 1 || entity.PropertyId != 7 || entity.OwnerId != 10 {
		t.Fatalf("project response lost identity: %#v", entity)
	}
}

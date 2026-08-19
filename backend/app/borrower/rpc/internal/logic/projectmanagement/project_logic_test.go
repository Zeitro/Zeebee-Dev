package projectmanagementlogic

import (
	"testing"
	"time"

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

	now := time.Date(2026, 8, 19, 12, 0, 0, 0, time.UTC)
	entity, err := projectEntity(&model.ProjectManagementProject{
		Base: model.Base{Id: 1, CreatedAt: now, UpdatedAt: now.Add(time.Hour)}, EntityId: 7, EntityType: "Transaction", Name: "Renovation", OwnerId: 10, OwnerName: "Owner",
	})
	if err != nil {
		t.Fatal(err)
	}
	if entity.Id != 1 || entity.EntityId != 7 || entity.EntityType != "Transaction" || entity.OwnerId != 10 || entity.OwnerName != "Owner" || entity.CreatedAt != now.Format(time.RFC3339) || entity.UpdatedAt != now.Add(time.Hour).Format(time.RFC3339) {
		t.Fatalf("project response lost identity: %#v", entity)
	}
}

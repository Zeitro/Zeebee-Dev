package projectmanagementlogic

import (
	"database/sql"
	"testing"
	"time"

	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
)

func TestTaskEntityKeepsTypeAndAssigneeSnapshots(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	updatedAt := now.Add(time.Hour)
	task := &model.ProjectManagementTask{
		Base:      model.Base{CreatedAt: now, UpdatedAt: updatedAt},
		ProjectId: 1, CreatorId: 2, CreatorName: "Former creator", Title: "Task",
		StartAt:   &sql.NullTime{Time: now, Valid: true},
		TaskType:  &dto.ProjectManagementTaskTypeSnapshotDto{Id: 3, Name: "Inspection", Description: "Historical type", Color: "#fff"},
		Assignees: []*model.ProjectManagementTaskAssignee{{Base: model.Base{CreatedAt: now, UpdatedAt: updatedAt}, OrganizationUserId: 4, OrganizationUserName: "Former executor"}},
	}
	entity, err := taskEntity(task, &dto.ProjectManagementTaskOverdueDto{Overdue: true, OverdueSeconds: 3600})
	if err != nil {
		t.Fatal(err)
	}
	if entity.TaskType == nil || entity.TaskType.Id != 3 || entity.TaskType.Name != "Inspection" || len(entity.Assignees) != 1 || entity.Assignees[0].OrganizationUserId != 4 || entity.Assignees[0].OrganizationUserName != "Former executor" || entity.StartAt != now.Format(time.RFC3339) || entity.CreatedAt != now.Format(time.RFC3339) || entity.UpdatedAt != updatedAt.Format(time.RFC3339) || entity.Assignees[0].CreatedAt != now.Format(time.RFC3339) || entity.Assignees[0].UpdatedAt != updatedAt.Format(time.RFC3339) || !entity.Overdue || entity.OverdueSeconds != 3600 {
		t.Fatalf("task snapshots were not preserved: %#v", entity)
	}
	attachment, err := attachmentEntity(&model.File{Base: model.Base{Id: 5, CreatedAt: now, UpdatedAt: updatedAt}, FileName: "task.pdf", FileUrl: "https://example.test/task.pdf", StorageType: constants.StorageTypeAwsS3})
	if err != nil || attachment.FileId != 5 || attachment.FileName != "task.pdf" || attachment.CreatedAt != now.Format(time.RFC3339) || attachment.UpdatedAt != updatedAt.Format(time.RFC3339) {
		t.Fatalf("attachment conversion failed: %#v, %v", attachment, err)
	}
}

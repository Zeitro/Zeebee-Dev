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
	task := &model.ProjectManagementTask{
		ProjectId: 1, CreatorId: 2, CreatorName: "Former creator", Title: "Task",
		StartAt:   &sql.NullTime{Time: now, Valid: true},
		TaskType:  &dto.ProjectManagementTaskTypeSnapshotDto{Id: 3, Name: "Inspection", Description: "Historical type", Color: "#fff"},
		Assignees: []*model.ProjectManagementTaskAssignee{{OrganizationUserId: 4, OrganizationUserName: "Former executor"}},
		Overdue:   true, OverdueSeconds: 3600,
	}
	entity, err := taskEntity(task)
	if err != nil {
		t.Fatal(err)
	}
	if entity.TaskType == nil || entity.TaskType.Id != 3 || entity.TaskType.Name != "Inspection" || len(entity.Assignees) != 1 || entity.Assignees[0].OrganizationUserId != 4 || entity.Assignees[0].OrganizationUserName != "Former executor" || entity.StartAt != now.Format(time.RFC3339) || !entity.Overdue || entity.OverdueSeconds != 3600 {
		t.Fatalf("task snapshots were not preserved: %#v", entity)
	}
	attachment, err := attachmentEntity(&model.File{Base: model.Base{Id: 5}, FileName: "task.pdf", FileUrl: "https://example.test/task.pdf", StorageType: constants.StorageTypeAwsS3})
	if err != nil || attachment.FileId != 5 || attachment.FileName != "task.pdf" {
		t.Fatalf("attachment conversion failed: %#v, %v", attachment, err)
	}
}

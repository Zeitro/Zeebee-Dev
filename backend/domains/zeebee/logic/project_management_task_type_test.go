package logic

import (
	"testing"

	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
)

func TestProjectManagementTaskTypeSnapshot(t *testing.T) {
	taskType := &model.ProjectManagementTaskType{Base: model.Base{Id: 9}, Name: "Inspection", Description: "Initial visit", Color: "#123456"}
	snapshot := ProjectManagementTaskTypeSnapshot(taskType)
	taskType.Name = "Changed"

	if snapshot.Id != 9 || snapshot.Name != "Inspection" || snapshot.Description != "Initial visit" || snapshot.Color != "#123456" {
		t.Fatalf("snapshot changed: %#v", snapshot)
	}
}

func TestValidateProjectManagementTaskType(t *testing.T) {
	if err := validateProjectManagementTaskType(dto.ProjectManagementTaskTypeDto{Name: "Inspection"}); err != nil {
		t.Fatalf("valid task type rejected: %v", err)
	}
	if err := validateProjectManagementTaskType(dto.ProjectManagementTaskTypeDto{Name: " \t"}); err == nil {
		t.Fatal("blank task type accepted")
	}
}

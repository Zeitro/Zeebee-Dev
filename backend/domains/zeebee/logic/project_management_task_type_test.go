package logic

import (
	"encoding/json"
	"testing"

	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
)

func TestProjectManagementTaskTypeSnapshot(t *testing.T) {
	taskType := &model.ProjectManagementTaskType{Base: model.Base{Id: 9}, Name: "Inspection", Description: "Initial visit", Color: "#123456"}
	snapshot, err := ProjectManagementTaskTypeSnapshot(taskType)
	if err != nil {
		t.Fatalf("snapshot: %v", err)
	}
	taskType.Name = "Changed"

	var got map[string]any
	if err := json.Unmarshal(snapshot, &got); err != nil {
		t.Fatalf("unmarshal snapshot: %v", err)
	}
	if got["id"] != float64(9) || got["name"] != "Inspection" || got["description"] != "Initial visit" || got["color"] != "#123456" {
		t.Fatalf("snapshot changed: %#v", got)
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

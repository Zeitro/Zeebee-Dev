package constants

import (
	"testing"

	"zeebee-backend/common/enums"
)

func TestProjectManagementPriority(t *testing.T) {
	for _, priority := range []ProjectManagementPriority{
		ProjectManagementPriorityLow,
		ProjectManagementPriorityMedium,
		ProjectManagementPriorityHigh,
		ProjectManagementPriorityCritical,
	} {
		if err := priority.Validate(); err != nil {
			t.Errorf("%q should be valid: %v", priority, err)
		}
	}
	if err := ProjectManagementPriority("invalid").Validate(); err == nil {
		t.Error("invalid priority should be rejected")
	}
	if opts := enums.Get(ProjectManagementPriorityName); len(opts) != 4 {
		t.Errorf("registered priorities = %d, want 4", len(opts))
	}
}

func TestProjectManagementEnums(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		valid   func() error
		invalid func() error
		count   int
	}{
		{"workflow kind", ProjectManagementWorkflowKindName, func() error { return ProjectManagementWorkflowKindProject.Validate() }, func() error { return ProjectManagementWorkflowKind("invalid").Validate() }, 2},
		{"owner type", ProjectManagementOwnerTypeName, func() error { return ProjectManagementOwnerTypeProject.Validate() }, func() error { return ProjectManagementOwnerType("invalid").Validate() }, 2},
		{"template type", ProjectManagementTemplateTypeName, func() error { return ProjectManagementTemplateTypePublic.Validate() }, func() error { return ProjectManagementTemplateType("invalid").Validate() }, 2},
		{"project type", ProjectManagementProjectTypeName, func() error { return ProjectManagementProjectTypePlanned.Validate() }, func() error { return ProjectManagementProjectType("invalid").Validate() }, 2},
		{"period unit", ProjectManagementPeriodUnitName, func() error { return ProjectManagementPeriodUnitMonth.Validate() }, func() error { return ProjectManagementPeriodUnit("invalid").Validate() }, 3},
		{"stage phase", ProjectManagementStagePhaseName, func() error { return ProjectManagementStagePhasePending.Validate() }, func() error { return ProjectManagementStagePhase("invalid").Validate() }, 3},
		{"member role", ProjectManagementMemberRoleName, func() error { return ProjectManagementMemberRoleProjectOwner.Validate() }, func() error { return ProjectManagementMemberRole("invalid").Validate() }, 4},
		{"candidate purpose", ProjectManagementMemberCandidatePurposeName, func() error { return ProjectManagementMemberCandidatePurposeAddMember.Validate() }, func() error { return ProjectManagementMemberCandidatePurpose("invalid").Validate() }, 3},
		{"dependency type", ProjectManagementDependencyTypeName, func() error { return ProjectManagementDependencyTypePreTask.Validate() }, func() error { return ProjectManagementDependencyType("invalid").Validate() }, 2},
		{"dependency rule", ProjectManagementDependencyRuleName, func() error { return ProjectManagementDependencyRuleFS.Validate() }, func() error { return ProjectManagementDependencyRule("invalid").Validate() }, 4},
		{"project entity type", ProjectManagementProjectEntityTypeName, func() error { return ProjectManagementEntityTypeBorrowerProperty.Validate() }, func() error { return ProjectManagementProjectEntityType("invalid").Validate() }, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.valid(); err != nil {
				t.Errorf("valid enum rejected: %v", err)
			}
			if err := tt.invalid(); err == nil {
				t.Error("invalid enum accepted")
			}
			if opts := enums.Get(tt.key); len(opts) != tt.count {
				t.Errorf("registered values = %d, want %d", len(opts), tt.count)
			}
		})
	}
}

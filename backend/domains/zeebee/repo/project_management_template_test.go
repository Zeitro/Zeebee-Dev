package repo

import (
	"testing"

	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
)

func TestProjectManagementTemplatePreloadsStages(t *testing.T) {
	db := newRepoTestDB(t)
	execTestSQL(t, db, []string{
		"CREATE TABLE project_management_templates (id integer primary key autoincrement, created_at datetime, updated_at datetime, is_del integer not null default 0, template_type text not null, organization_id integer not null default 0, name text not null default '', description text not null default '', status text not null default 'Active', is_default boolean not null default false, sort_order integer not null default 0, project_type text, start_at datetime, end_at datetime, period_value integer, period_unit text, source_project_id integer not null default 0)",
		"CREATE TABLE project_management_stages (id integer primary key autoincrement, created_at datetime, updated_at datetime, is_del integer not null default 0, entity_type text not null, entity_id integer not null default 0, name text not null default '', description text not null default '', phase text not null, color text not null default '', sort_order integer not null default 0)",
	})
	template := &model.ProjectManagementTemplate{TemplateType: constants.ProjectManagementTemplateTypeOrganization, Name: "Template"}
	stages := []*model.ProjectManagementStage{
		{EntityType: constants.ProjectManagementOwnerTypeTemplate, EntityId: 1, Name: "Pending", Phase: constants.ProjectManagementStagePhasePending},
		{EntityType: constants.ProjectManagementOwnerTypeProject, EntityId: 1, Name: "Project stage", Phase: constants.ProjectManagementStagePhasePending},
	}
	if err := db.Create(template).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Create(stages).Error; err != nil {
		t.Fatal(err)
	}

	got := new(model.ProjectManagementTemplate)
	if err := db.Preload("Stages").First(got, template.Id).Error; err != nil {
		t.Fatal(err)
	}
	if len(got.Stages) != 1 || got.Stages[0].Id != stages[0].Id {
		t.Fatalf("unexpected template stages: %#v", got.Stages)
	}
	listed, err := NewProjectManagementTemplateRepo(db).ListStages(t.Context(), template.Id)
	if err != nil || len(listed) != 1 || listed[0].Id != stages[0].Id {
		t.Fatalf("listed template stages: %#v, %v", listed, err)
	}
}

package repo

import (
	"context"
	"testing"

	"gorm.io/gorm"
)

func TestProjectManagementDependencyReachability(t *testing.T) {
	db := newRepoTestDB(t)
	if err := db.Exec("CREATE TABLE project_management_task_dependencies (id integer primary key, is_del integer not null default 0, task_id integer not null, dependency_task_id integer not null, dependency_type text not null default 'PreTask', dependency_rule text not null default 'FS')").Error; err != nil {
		t.Fatal(err)
	}
	for _, edge := range [][2]int64{{2, 1}, {3, 2}, {4, 3}} {
		if err := db.Exec("INSERT INTO project_management_task_dependencies (task_id, dependency_task_id) VALUES (?, ?)", edge[0], edge[1]).Error; err != nil {
			t.Fatal(err)
		}
	}
	repo := NewProjectManagementDependencyRepo(db)
	if err := db.Transaction(func(tx *gorm.DB) error {
		found, err := repo.HasPathWithTx(context.Background(), tx, 1, 4)
		if err != nil {
			return err
		}
		if !found {
			t.Fatal("recursive dependency path was not found")
		}
		found, err = repo.HasPathWithTx(context.Background(), tx, 4, 1)
		if err != nil {
			return err
		}
		if found {
			t.Fatal("reverse dependency path unexpectedly exists")
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
}

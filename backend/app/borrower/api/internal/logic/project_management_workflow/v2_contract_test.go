package project_management_workflow

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestGeneratedRoutesExposeV2WorkflowSurface(t *testing.T) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("failed to locate generated route source")
	}
	routes, err := os.ReadFile(filepath.Join(filepath.Dir(file), "../../handler/routes.go"))
	if err != nil {
		t.Fatalf("read generated routes: %v", err)
	}
	source := string(routes)
	for _, obsolete := range []string{
		"/project_management/workflow/:id/swap_stage_order",
		"StageSwapOrderHandler(serverCtx)",
		"/project_management/project/:id/workflow",
		"project_management_project.WorkflowDetailHandler(serverCtx)",
		"project_management_project.WorkflowUpdateHandler(serverCtx)",
	} {
		if strings.Contains(source, obsolete) {
			t.Fatalf("obsolete workflow route remains registered: %s", obsolete)
		}
	}
	for _, required := range []string{
		"/project_management/project/workflows/page",
		"ProjectManagementProjectWorkflowPageHandler(serverCtx)",
		"/project_management/project/:project_id/workflows/:workflow_id/stages",
	} {
		if !strings.Contains(source, required) {
			t.Fatalf("v2 workflow route is not registered: %s", required)
		}
	}
}

package routex

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/pathvar"
)

func TestGetGoZeroRoutePatternUsesRegisteredTemplate(t *testing.T) {
	registerRoutesForTest(t, []rest.Route{
		{Method: http.MethodGet, Path: "/api/v1/project_management/workflow/:workflow_id/stage/:id"},
		{Method: http.MethodGet, Path: "/api/v1/project_management/workflow/:workflow_id/stage/all"},
		{Method: http.MethodGet, Path: "/api/v1/project_management/task/:id"},
		{Method: http.MethodPost, Path: "/api/v1/project_management/task/:task_id"},
	})

	tests := []struct {
		name     string
		method   string
		path     string
		vars     map[string]string
		expected string
	}{
		{
			name:     "same parameter values",
			method:   http.MethodGet,
			path:     "/api/v1/project_management/workflow/1/stage/1",
			vars:     map[string]string{"workflow_id": "1", "id": "1"},
			expected: "/api/v1/project_management/workflow/:workflow_id/stage/:id",
		},
		{
			name:     "different parameter values",
			method:   http.MethodGet,
			path:     "/api/v1/project_management/workflow/12/stage/34",
			vars:     map[string]string{"workflow_id": "12", "id": "34"},
			expected: "/api/v1/project_management/workflow/:workflow_id/stage/:id",
		},
		{
			name:     "static route wins over parameter route",
			method:   http.MethodGet,
			path:     "/api/v1/project_management/workflow/1/stage/all",
			vars:     map[string]string{"workflow_id": "1"},
			expected: "/api/v1/project_management/workflow/:workflow_id/stage/all",
		},
		{
			name:     "method selects matching template",
			method:   http.MethodPost,
			path:     "/api/v1/project_management/task/1",
			vars:     map[string]string{"task_id": "1"},
			expected: "/api/v1/project_management/task/:task_id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := pathvar.WithVars(httptest.NewRequest(tt.method, tt.path, nil), tt.vars)
			assert.Equal(t, tt.expected, GetGoZeroRoutePattern(req))
		})
	}
}

func TestGetGoZeroRoutePatternFallback(t *testing.T) {
	registerRoutesForTest(t, nil)

	t.Run("single path variable", func(t *testing.T) {
		req := pathvar.WithVars(httptest.NewRequest(http.MethodGet, "/api/v1/resource/42", nil), map[string]string{"id": "42"})
		assert.Equal(t, "/api/v1/resource/:id", GetGoZeroRoutePattern(req))
	})

	t.Run("no path variables", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/health//", nil)
		assert.Equal(t, "/api/v1/health", GetGoZeroRoutePattern(req))
	})
}

func registerRoutesForTest(t *testing.T, routes []rest.Route) {
	t.Helper()
	routePatterns.Lock()
	previous := routePatterns.byMethod
	routePatterns.Unlock()
	t.Cleanup(func() {
		routePatterns.Lock()
		routePatterns.byMethod = previous
		routePatterns.Unlock()
	})
	RegisterRoutes(routes)
}

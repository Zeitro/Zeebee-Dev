package routex

import (
	"net/http"
	"strings"

	"github.com/zeromicro/go-zero/rest/pathvar"
)

func GetGoZeroRoutePattern(r *http.Request) string {
	path := r.URL.Path
	vars := pathvar.Vars(r)

	// If no path variables, return the path as-is
	if len(vars) == 0 {
		return path
	}

	// Reconstruct the route pattern by replacing parameter values with parameter names
	pattern := path
	for key, value := range vars {
		// Replace the value in the path with :key pattern
		// e.g., /api/v1/resource/7 -> /api/v1/resource/:id
		pattern = strings.Replace(pattern, "/"+value, "/:"+key, 1)
	}

	return pattern
}

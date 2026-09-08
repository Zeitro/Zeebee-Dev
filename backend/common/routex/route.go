package routex

import (
	"net/http"
	"path"
	"strings"
	"sync"

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/pathvar"
)

var routePatterns = struct {
	sync.RWMutex
	byMethod map[string][]string
}{byMethod: make(map[string][]string)}

// RegisterRoutes stores the exact route templates registered by the API server.
func RegisterRoutes(routes []rest.Route) {
	byMethod := make(map[string][]string)
	for _, route := range routes {
		byMethod[route.Method] = append(byMethod[route.Method], route.Path)
	}
	routePatterns.Lock()
	routePatterns.byMethod = byMethod
	routePatterns.Unlock()
}

func GetGoZeroRoutePattern(r *http.Request) string {
	requestPath := path.Clean(r.URL.Path)
	if pattern := registeredRoutePattern(r.Method, requestPath); pattern != "" {
		return pattern
	}

	vars := pathvar.Vars(r)

	// If no path variables, return the path as-is
	if len(vars) == 0 {
		return requestPath
	}

	// Reconstruct the route pattern by replacing parameter values with parameter names
	pattern := requestPath
	for key, value := range vars {
		// Replace the value in the path with :key pattern
		// e.g., /api/v1/resource/7 -> /api/v1/resource/:id
		pattern = strings.Replace(pattern, "/"+value, "/:"+key, 1)
	}

	return pattern
}

func registeredRoutePattern(method, requestPath string) string {
	routePatterns.RLock()
	patterns := routePatterns.byMethod[method]
	routePatterns.RUnlock()

	bestPattern, bestScore := "", -1
	for _, pattern := range patterns {
		score, matched := matchRoutePattern(pattern, requestPath)
		if matched && score > bestScore {
			bestPattern, bestScore = pattern, score
		}
	}
	return bestPattern
}

func matchRoutePattern(pattern, requestPath string) (int, bool) {
	patternParts := strings.Split(strings.Trim(pattern, "/"), "/")
	pathParts := strings.Split(strings.Trim(requestPath, "/"), "/")
	if len(patternParts) != len(pathParts) {
		return 0, false
	}

	staticParts := 0
	for i, part := range patternParts {
		if strings.HasPrefix(part, ":") {
			continue
		}
		if part != pathParts[i] {
			return 0, false
		}
		staticParts++
	}
	return staticParts, true
}

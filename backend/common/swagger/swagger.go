package swagger

import (
	"encoding/json"
	"sort"
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/ast"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/parser"
)

type ApiGroup map[string][]Api

type Api struct {
	Path        string `json:"path"`
	Method      string `json:"method"`
	Name        string `json:"name"`
	Description string `json:"description"`
	// AuthRequired       bool   `json:"auth_required"`
	// PermissionRequired bool   `json:"permission_required"`
}

type swaggerOperation struct {
	Summary     string   `json:"summary"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

// ParseSwagger extracts the route manifest from generated Swagger JSON.
func ParseSwagger(content []byte) (ApiGroup, error) {
	var doc struct {
		Paths map[string]map[string]json.RawMessage `json:"paths"`
	}
	if err := json.Unmarshal(content, &doc); err != nil {
		return nil, err
	}

	paths := make([]string, 0, len(doc.Paths))
	for path := range doc.Paths {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	groups := make(ApiGroup)
	for _, path := range paths {
		methods := make([]string, 0, len(doc.Paths[path]))
		for method := range doc.Paths[path] {
			if isHTTPMethod(method) {
				methods = append(methods, method)
			}
		}
		sort.Strings(methods)
		for _, method := range methods {
			var route swaggerOperation
			if err := json.Unmarshal(doc.Paths[path][method], &route); err != nil {
				return nil, err
			}
			name := strings.TrimSpace(route.Summary)
			if name == "" {
				name = strings.TrimSpace(route.Description)
			}
			if name == "" {
				name = path
			}
			for _, tag := range route.Tags {
				if tag == "" {
					continue
				}
				groups[tag] = append(groups[tag], Api{
					Path:        normalizePath(path),
					Method:      strings.ToUpper(method),
					Name:        name,
					Description: name,
				})
			}
		}
	}

	return groups, nil
}

func isHTTPMethod(method string) bool {
	switch strings.ToLower(method) {
	case "get", "put", "post", "delete", "options", "head", "patch", "trace":
		return true
	default:
		return false
	}
}

func normalizePath(path string) string {
	var normalized strings.Builder
	for i := 0; i < len(path); i++ {
		if path[i] != '{' {
			normalized.WriteByte(path[i])
			continue
		}
		end := strings.IndexByte(path[i+1:], '}')
		if end < 0 {
			normalized.WriteByte(path[i])
			continue
		}
		normalized.WriteByte(':')
		normalized.WriteString(path[i+1 : i+1+end])
		i += end + 1
	}
	return normalized.String()
}

func ParseApi(content string) (ApiGroup, error) {
	if strings.TrimSpace(content) == "" {
		return make(ApiGroup), nil
	}

	p := parser.New("stdin", content)
	a := p.Parse()
	if a == nil || len(a.Stmts) == 0 {
		return make(ApiGroup), nil
	}

	groups := make(ApiGroup)
	for _, stmt := range a.Stmts {
		svc, ok := stmt.(*ast.ServiceStmt)
		if !ok {
			continue
		}
		tag := getAtServerValue(svc.AtServerStmt, "tags")
		prefix := getAtServerValue(svc.AtServerStmt, "prefix")
		if tag == "" {
			continue
		}
		for _, route := range svc.Routes {
			name := getRouteDoc(route)
			groups[tag] = append(groups[tag], Api{
				Path:        prefix + route.Route.Path.Format(""),
				Method:      strings.ToUpper(route.Route.Method.Token.Text),
				Name:        name,
				Description: name,
			})
		}
	}

	return groups, nil
}

func getAtServerValue(atServer *ast.AtServerStmt, key string) string {
	if atServer == nil {
		return ""
	}
	for _, v := range atServer.Values {
		if strings.TrimSuffix(v.Key.Token.Text, ":") == key {
			return v.Value.RawText()
		}
	}
	return ""
}

func getRouteDoc(route *ast.ServiceItemStmt) string {
	if route.AtDoc != nil {
		switch v := route.AtDoc.(type) {
		case *ast.AtDocLiteralStmt:
			return strings.Trim(v.Value.Token.Text, "\"")
		case *ast.AtDocGroupStmt:
			for _, kv := range v.Values {
				if strings.TrimSuffix(kv.Key.Token.Text, ":") == "summary" {
					return kv.Value.RawText()
				}
			}
		}
	}
	head, _ := route.CommentGroup()
	return strings.Join(head.List(), " ")
}

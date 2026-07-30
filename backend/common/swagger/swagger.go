package swagger

import (
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

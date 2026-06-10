package {{.PkgName}}

import (
	"net/http"

    "github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	{{.ImportPackages}}


	"zeebee-backend/common/response"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/validator"
)

{{if .HasDoc}}{{.Doc}}{{end}}
func {{.HandlerName}}(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
        defer func() {
            if err != nil {
                logx.WithContext(r.Context()).
                    WithFields(
                        logx.Field("stacks", errorx.GetErrorStack(err)),
                        logx.Field("service_name", svcCtx.Config.Name),
                    ).
                    Error(err)
            }
        }()
		{{if .HasRequest}}var req types.{{.RequestType}}
		if err = httpx.Parse(r, &req); err != nil {
			response.BadRequest(w, errorx.GetErrorMsg(err))
			return
		}

		v, err := validator.GetEnValidator()
        if err != nil {
            response.InternalServerError(w, errorx.GetErrorMsg(err))
            return
        }
        if err := v.Validate(&req); err != nil {
            response.BadRequest(w, errorx.GetErrorMsg(err))
            return
        }

		{{end}}l := {{.LogicName}}.New{{.LogicType}}(r.Context(), svcCtx)
		{{if .HasResp}}resp, {{end}}err := l.{{.Call}}({{if .HasRequest}}&req{{end}})
		if err != nil {
			response.BadRequest(w, errorx.GetErrorMsg(err))
		} else {
			{{if .HasResp}}response.Success(w, resp){{else}}response.Success(w, nil){{end}}
		}
	}
}

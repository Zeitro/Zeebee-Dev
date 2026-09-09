package {{.packageName}}

import (
	"context"

	{{.imports}}

	"zeebee-backend/common/errorx"
)

type {{.logicName}} struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

}

func New{{.logicName}}(ctx context.Context,svcCtx *svc.ServiceContext) *{{.logicName}} {
	return &{{.logicName}}{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}
{{.functions}}

package handler

import (
	"zeebee-backend/webserver/internal/svc"
)

type BaseHandler struct {
	svcCtx *svc.ServiceContext
}

func NewBaseHandler(svcCtx *svc.ServiceContext) *BaseHandler {
	return &BaseHandler{
		svcCtx: svcCtx,
	}
}

package project_management_template

import (
	"context"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/domains/zeebee/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProjectManagementTemplateDefaultLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// set default template
func NewProjectManagementTemplateDefaultLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectManagementTemplateDefaultLogic {
	return &ProjectManagementTemplateDefaultLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProjectManagementTemplateDefaultLogic) ProjectManagementTemplateDefault(req *types.PathIdRequest) (resp *types.EmptyResponse, err error) {
	if _, err = l.svcCtx.ProjectManagementTemplateRpcClient.Defaults(l.ctx, &proto.PathIdRequest{Data: &proto.PathIdRequestData{Id: req.Id}}); err != nil {
		return nil, err
	}
	return &types.EmptyResponse{}, nil
}

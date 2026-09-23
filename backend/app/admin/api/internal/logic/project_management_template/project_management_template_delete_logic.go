package project_management_template

import (
	"context"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/domains/zeebee/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProjectManagementTemplateDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// delete
func NewProjectManagementTemplateDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectManagementTemplateDeleteLogic {
	return &ProjectManagementTemplateDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProjectManagementTemplateDeleteLogic) ProjectManagementTemplateDelete(req *types.PathIdRequest) (resp *types.EmptyResponse, err error) {
	_, err = l.svcCtx.ProjectManagementTemplateRpcClient.Delete(l.ctx, &proto.PathIdRequest{Data: &proto.PathIdRequestData{Id: req.Id}})
	if err != nil {
		return nil, err
	}
	return &types.EmptyResponse{}, nil
}

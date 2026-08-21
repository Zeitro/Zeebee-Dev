package project_management_task

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type AttachmentPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// page project task attachments
func NewAttachmentPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AttachmentPageLogic {
	return &AttachmentPageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AttachmentPageLogic) AttachmentPage(req *types.ProjectManagementAttachmentPageRequest) (resp *types.ProjectManagementAttachmentPageResponse, err error) {
	data, err := copierx.Trans[types.ProjectManagementAttachmentPageRequest, borrower.ProjectManagementAttachmentPageRequestData](req)
	if err != nil {
		return nil, err
	}
	data.TaskId = req.Id
	rpcResp, err := l.svcCtx.ProjectManagementTaskRpcClient.AttachmentPage(l.ctx, &borrower.ProjectManagementAttachmentPageRequest{Data: data})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	resp = &types.ProjectManagementAttachmentPageResponse{Page: rpcResp.Data.Page, Limit: rpcResp.Data.Limit, Total: rpcResp.Data.Total, TotalPage: rpcResp.Data.TotalPage, List: make([]*types.ProjectManagementAttachmentEntity, 0, len(rpcResp.Data.List))}
	for _, file := range rpcResp.Data.List {
		entity, err := attachmentEntity(file)
		if err != nil {
			return nil, err
		}
		resp.List = append(resp.List, entity)
	}
	return resp, nil
}

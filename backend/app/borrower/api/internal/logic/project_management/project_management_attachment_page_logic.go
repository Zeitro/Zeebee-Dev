package project_management

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProjectManagementAttachmentPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// page project task attachments
func NewProjectManagementAttachmentPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectManagementAttachmentPageLogic {
	return &ProjectManagementAttachmentPageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProjectManagementAttachmentPageLogic) ProjectManagementAttachmentPage(req *types.ProjectManagementAttachmentPageRequest) (resp *types.ProjectManagementAttachmentPageResponse, err error) {
	data, err := copierx.Trans[types.ProjectManagementAttachmentPageRequest, borrower.ProjectManagementAttachmentPageRequestData](req)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	data.TaskId = req.Id
	rpcResp, err := l.svcCtx.ProjectManagementRpcClient.AttachmentPage(l.ctx, &borrower.ProjectManagementAttachmentPageRequest{Data: data})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	return copierx.Trans[borrower.ProjectManagementAttachmentPageResponseData, types.ProjectManagementAttachmentPageResponse](rpcResp.Data)
}

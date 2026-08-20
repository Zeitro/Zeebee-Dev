package projectmanagementtasklogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/proto"
)

type AttachmentPageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAttachmentPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AttachmentPageLogic {
	return &AttachmentPageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AttachmentPageLogic) AttachmentPage(in *borrower.ProjectManagementAttachmentPageRequest) (resp *borrower.ProjectManagementAttachmentPageResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	base := logic.GetBaseLogic(l.svcCtx)
	claimValue, err := base.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	claim := &claimValue
	if _, err = base.ProjectManagementTaskLogic.DetailForProject(l.ctx, claim, in.Data.ProjectId, in.Data.TaskId); err != nil {
		return nil, err
	}
	page, err := base.ProjectManagementAttachmentLogic.Page(l.ctx, claim, in.Data.TaskId, dto.ProjectManagementAttachmentPageRequestDto{Page: in.Data.Page, Limit: in.Data.Limit})
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	data := &borrower.ProjectManagementAttachmentPageResponseData{Page: page.Page, Limit: page.Limit, Total: page.Total, TotalPage: page.TotalPage, List: make([]*proto.FileEntity, 0, len(page.List))}
	for _, file := range page.List {
		entity, err := attachmentEntity(file)
		if err != nil {
			return nil, errorx.GetFailedTransRespErr(err)
		}
		data.List = append(data.List, entity)
	}
	return &borrower.ProjectManagementAttachmentPageResponse{Data: data}, nil
}

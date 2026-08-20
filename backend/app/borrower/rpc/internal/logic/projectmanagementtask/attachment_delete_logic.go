package projectmanagementtasklogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/proto"
)

type AttachmentDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAttachmentDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AttachmentDeleteLogic {
	return &AttachmentDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AttachmentDeleteLogic) AttachmentDelete(in *borrower.ProjectManagementAttachmentDeleteRequest) (resp *proto.EmptyResponse, err error) {
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
	if err = base.ProjectManagementAttachmentLogic.Delete(l.ctx, claim, in.Data.TaskId, in.Data.FileId); err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	return &proto.EmptyResponse{}, nil
}

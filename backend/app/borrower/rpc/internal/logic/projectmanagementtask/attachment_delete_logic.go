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
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewAttachmentDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AttachmentDeleteLogic {
	return &AttachmentDeleteLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *AttachmentDeleteLogic) AttachmentDelete(in *borrower.ProjectManagementAttachmentDeleteRequest) (resp *proto.EmptyResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	base := l.baseLogic
	claim, err := base.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	if err = base.ProjectManagementAttachmentLogic.Delete(l.ctx, claim, in.Data.TaskId, in.Data.FileId); err != nil {
		return nil, err
	}
	return &proto.EmptyResponse{}, nil
}

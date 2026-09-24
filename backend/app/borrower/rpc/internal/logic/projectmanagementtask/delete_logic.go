package projectmanagementtasklogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/proto"
)

type DeleteLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *DeleteLogic) Delete(in *borrower.ProjectManagementTaskDeleteRequest) (resp *proto.EmptyResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	base := l.baseLogic
	claim, err := base.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	if _, err = base.ProjectManagementTaskLogic.Detail(l.ctx, claim, in.Data.Id); err != nil {
		return nil, err
	}
	if err = base.ProjectManagementTaskLogic.Delete(l.ctx, claim, in.Data.Id); err != nil {
		return nil, err
	}
	return &proto.EmptyResponse{}, nil
}

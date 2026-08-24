package projectmanagementprojectmemberlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/errorx"
)

type DetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DetailLogic) Detail(in *borrower.ProjectManagementMemberDetailRequest) (resp *borrower.ProjectManagementMemberDetailResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	claim, err := logic.GetBaseLogic(l.svcCtx).GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	member, err := logic.GetBaseLogic(l.svcCtx).ProjectManagementMemberLogic.Detail(l.ctx, claim, in.Data.ProjectId, in.Data.Id)
	if err != nil {
		return nil, err
	}
	data, err := memberEntity(member)
	if err != nil {
		return nil, err
	}
	return &borrower.ProjectManagementMemberDetailResponse{Data: data}, nil
}

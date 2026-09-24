package projectmanagementprojectmemberlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/proto"
)

type DetailLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *DetailLogic) Detail(in *borrower.ProjectManagementMemberDetailRequest) (resp *borrower.ProjectManagementMemberDetailResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	claim, err := l.baseLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	member, err := l.baseLogic.ProjectManagementMemberLogic.Detail(l.ctx, claim, in.Data.ProjectId, in.Data.Id)
	if err != nil {
		return nil, err
	}
	data, err := copierx.Trans[dto.ProjectManagementMemberDetailDto, proto.ProjectManagementMemberEntity](member, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}
	return &borrower.ProjectManagementMemberDetailResponse{Data: data}, nil
}

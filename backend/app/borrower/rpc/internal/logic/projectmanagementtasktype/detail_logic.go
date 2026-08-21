package projectmanagementtasktypelogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"

	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/proto"
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

func (l *DetailLogic) Detail(in *proto.PathIdRequest) (resp *borrower.ProjectManagementTaskTypeDetailResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	base := logic.GetBaseLogic(l.svcCtx)
	claim, err := base.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	taskType, err := base.ProjectManagementTaskTypeLogic.Detail(l.ctx, claim.OrganizationId, in.Data.Id)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	data, err := copierx.Trans[model.ProjectManagementTaskType, proto.ProjectManagementTaskTypeEntity](taskType, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	return &borrower.ProjectManagementTaskTypeDetailResponse{Data: data}, nil
}

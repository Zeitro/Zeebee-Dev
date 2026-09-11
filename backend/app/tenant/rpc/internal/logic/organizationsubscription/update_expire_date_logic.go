package organizationsubscriptionlogic

import (
	"context"
	"database/sql"

	"zeebee-backend/app/tenant/rpc/internal/logic"
	"zeebee-backend/app/tenant/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/timex"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/proto"
)

type UpdateExpireDateLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewUpdateExpireDateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateExpireDateLogic {
	return &UpdateExpireDateLogic{ctx: ctx, svcCtx: svcCtx, baseLogic: logic.GetBaseLogic(svcCtx)}
}

// UpdateExpireDate overrides a subscription's required expire date (MM/DD/YYYY).
func (l *UpdateExpireDateLogic) UpdateExpireDate(in *tenant.OrganizationSubscriptionUpdateExpireDateRequest) (resp *proto.EmptyResponse, err error) {
	if in.Data == nil {
		return nil, errorx.ErrorEmptyRequestData
	}
	if _, err = logic.ParseUserType(in.Data.UserType); err != nil {
		return
	}

	if in.Data.ExpireDate == "" {
		return nil, errorx.New("expire_date is required")
	}
	t, err := timex.FormateTimeToStrWithUsStyle(in.Data.ExpireDate)
	if err != nil {
		return nil, errorx.New("invalid expire_date, expect MM/DD/YYYY")
	}
	expireDate := sql.NullTime{Time: t, Valid: true}
	input := dto.OrganizationSubscriptionUpdateExpireDateRequestDto{
		Id:         in.Data.Id,
		ExpireDate: &expireDate,
	}
	if err = l.baseLogic.OrganizationSubscriptionLogic.UpdateExpireDate(l.ctx, input); err != nil {
		return
	}
	return &proto.EmptyResponse{}, nil
}

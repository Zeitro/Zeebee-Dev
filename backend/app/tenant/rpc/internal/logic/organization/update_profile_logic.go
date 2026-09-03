package organizationlogic

import (
	"context"

	"zeebee-backend/app/tenant/rpc/internal/logic"
	"zeebee-backend/app/tenant/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/proto"
)

type UpdateProfileLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewUpdateProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateProfileLogic {
	return &UpdateProfileLogic{ctx: ctx, svcCtx: svcCtx, baseLogic: logic.GetBaseLogic(svcCtx)}
}

func (l *UpdateProfileLogic) UpdateProfile(in *tenant.OrganizationUpdateProfileRequest) (*tenant.OrganizationDetailResponse, error) {
	if in == nil || in.Data == nil {
		return nil, errorx.ErrorEmptyRequestData
	}
	userType, err := logic.ParseUserType(in.Data.UserType)
	if err != nil {
		return nil, err
	}
	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, errorx.Wrap(err, "fail to get user claim from context")
	}
	organization, err := l.svcCtx.OrganizationRepo.GetOrganization(l.ctx, claim.OrganizationId, userType, false)
	if err != nil {
		return nil, errorx.Wrap(err, "fail to get organization")
	}
	if err = copierx.Assign(in.Data, organization); err != nil {
		return nil, err
	}
	if _, err = l.svcCtx.OrganizationRepo.UpdateOneById(l.ctx, claim.OrganizationId, organization); err != nil {
		return nil, errorx.Wrap(err, "fail to update organization")
	}
	result, err := l.baseLogic.OrganizationLogic.DetailWithType(l.ctx, claim.OrganizationId, userType)
	if err != nil {
		return nil, err
	}
	data, err := copierx.Trans[dto.OrganizationEntityDto, proto.OrganizationEntity](result, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}
	return &tenant.OrganizationDetailResponse{Data: data}, nil
}

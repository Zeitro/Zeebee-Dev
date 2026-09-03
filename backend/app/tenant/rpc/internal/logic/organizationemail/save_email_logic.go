package organizationemaillogic

import (
	"context"

	"zeebee-backend/app/tenant/rpc/internal/logic"
	"zeebee-backend/app/tenant/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/proto"
)

type SaveEmailLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewSaveEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveEmailLogic {
	return &SaveEmailLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

// SaveEmail persists an organization's sender email config (SMTP / strategy).
func (l *SaveEmailLogic) SaveEmail(in *tenant.OrganizationEmailSaveEmailRequest) (resp *proto.EmptyResponse, err error) {
	if in.Data == nil {
		return nil, errorx.ErrorEmptyRequestData
	}
	userType, err := logic.ParseUserType(in.Data.UserType)
	if err != nil {
		return
	}

	input := dto.OrganizationEmailSaveEmailRequestDto{
		OrganizationId:   in.Data.OrganizationId,
		UserId:           in.Data.UserId,
		Email:            in.Data.Email,
		Purpose:          constants.EmailPurpose(in.Data.Purpose),
		AuthStrategy:     constants.EmailAuthStrategy(in.Data.AuthStrategy),
		OrganizationType: userType,
	}
	if err = l.baseLogic.OrganizationEmailLogic.SaveEmail(l.ctx, input); err != nil {
		return
	}
	return &proto.EmptyResponse{}, nil
}

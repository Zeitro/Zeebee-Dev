package organizationuserlogic

import (
	"context"

	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type CheckInvitationLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewCheckInvitationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckInvitationLogic {
	return &CheckInvitationLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *CheckInvitationLogic) CheckInvitation(in *lender.OrganizationUserCheckInvitationRequest) (resp *lender.OrganizationUserCheckInvitationResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
				logx.Field("service_name", l.svcCtx.Config.Name),
			).Error(err)
		}
	}()
	input, err := copierx.Trans[lender.OrganizationUserCheckInvitationRequestData, dto.OrganizationUserCheckInvitationCodeRequestDto](in.Data)
	if err != nil {
		return
	}
	input.UserType = constants.UserTypeLender
	isExpired, err := l.baseLogic.OrganizationUserLogic.CheckInvitationCode(l.ctx, *input)
	if err != nil {
		return nil, err
	}

	return &lender.OrganizationUserCheckInvitationResponse{
		Data: &lender.OrganizationUserCheckInvitationResponseData{
			IsExpired: isExpired,
		},
	}, nil
}

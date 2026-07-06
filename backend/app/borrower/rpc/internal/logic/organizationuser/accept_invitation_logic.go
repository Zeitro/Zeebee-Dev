package organizationuserlogic

import (
	"context"

	"github.com/jinzhu/copier"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/proto"

	"zeebee-backend/common/errorx"
)

type AcceptInvitationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewAcceptInvitationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AcceptInvitationLogic {
	return &AcceptInvitationLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *AcceptInvitationLogic) AcceptInvitation(in *borrower.OrganizationUserAcceptInvitationRequest) (resp *proto.UserLoginResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
				logx.Field("service_name", l.svcCtx.Config.Name),
			).Error(err)
		}
	}()

	input, err := copierx.Trans[borrower.OrganizationUserAcceptInvitationRequestData, dto.OrganizationUserAcceptInvitationRequestDto](in.Data)
	if err != nil {
		return nil, err
	}

	input.UserType = constants.UserTypeBorrower
	userBo, err := l.baseLogic.OrganizationUserLogic.AcceptInvitation(l.ctx, *input)
	if err != nil {
		return nil, err
	}
	token, err := l.baseLogic.UserLogic.GenerateJwtTokenById(l.ctx, userBo.Id, userBo.DefaultOrganizationId, constants.UserTypeBorrower)
	if err != nil {
		return nil, err
	}
	data := &proto.UserLoginResponseData{}
	err = copier.CopyWithOption(data, token, copier.Option{DeepCopy: true})
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	return &proto.UserLoginResponse{Data: data}, nil
}

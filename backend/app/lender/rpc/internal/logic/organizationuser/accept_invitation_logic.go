package organizationuserlogic

import (
	"context"

	"github.com/jinzhu/copier"

	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/bo"
	"zeebee-backend/domains/zeebee/data/dto"

	"github.com/zeromicro/go-zero/core/logx"

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
		baseLogic: logic.NewBaseLogic(ctx, svcCtx),
	}
}

func (l *AcceptInvitationLogic) AcceptInvitation(in *lender.OrganizationUserAcceptInvitationRequest) (resp *lender.UserLoginResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	input := dto.OrganizationUserAcceptInvitationRequestDto{}
	err = copier.Copy(&input, &in.Data)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	input.UserType = constants.UserTypeLender
	userBo, err := l.baseLogic.OrganizationUserLogic.AcceptInvitation(l.ctx, input)
	if err != nil {
		return nil, err
	}
	token, err := l.baseLogic.UserLogic.GenerateJwtTokenById(l.ctx, userBo.Id, userBo.DefaultOrganizationId, constants.UserTypeLender)
	if err != nil {
		return nil, err
	}
	data, err := copierx.Trans[bo.UserJwtTokenBo, lender.UserLoginResponseData](token)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	return &lender.UserLoginResponse{Data: data}, nil
}

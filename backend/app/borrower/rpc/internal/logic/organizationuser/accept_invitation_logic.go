package organizationuserlogic

import (
	"context"

	"github.com/jinzhu/copier"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"

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

func (l *AcceptInvitationLogic) AcceptInvitation(in *borrower.OrganizationUserAcceptInvitationRequest) (resp *borrower.UserLoginResponse, err error) {
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
	input.UserType = constants.UserTypeBorrower
	userBo, err := l.baseLogic.OrganizationUserLogic.AcceptInvitation(l.ctx, input)
	if err != nil {
		return nil, err
	}
	token, err := l.baseLogic.UserLogic.GenerateJwtTokenById(l.ctx, userBo.Id, userBo.DefaultOrganizationId, constants.UserTypeBorrower)
	if err != nil {
		return nil, err
	}
	data := &borrower.UserLoginResponseData{}
	err = copier.CopyWithOption(data, token, copier.Option{DeepCopy: true})
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	return &borrower.UserLoginResponse{Data: data}, nil
}

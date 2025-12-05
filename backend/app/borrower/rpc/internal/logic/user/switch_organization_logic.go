package userlogic

import (
	"context"

	"github.com/jinzhu/copier"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/constants"

	"zeebee-backend/common/errorx"
)

type SwitchOrganizationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewSwitchOrganizationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SwitchOrganizationLogic {
	return &SwitchOrganizationLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.NewBaseLogic(ctx, svcCtx),
	}
}

func (l *SwitchOrganizationLogic) SwitchOrganization(in *borrower.SwitchOrganizationRequest) (resp *borrower.UserLoginResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	user, organizationUser, err := l.baseLogic.UserLogic.SwitchOrganization(l.ctx, in.Data.UserId, in.Data.OrganizationId, constants.UserTypeBorrower)
	if err != nil {
		return nil, err
	}

	token, err := l.baseLogic.UserLogic.GenerateJwtToken(user, organizationUser, constants.UserTypeBorrower)
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

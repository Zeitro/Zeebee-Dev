package userlogic

import (
	"context"

	"github.com/jinzhu/copier"

	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/constants"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/stringx"
	"zeebee-backend/common/utils"
	logic2 "zeebee-backend/domains/zeebee/logic"
	"zeebee-backend/domains/zeebee/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	BaseLogic *logic.BaseLogic
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		BaseLogic: logic.NewBaseLogic(ctx, svcCtx),
	}
}

func (l *RegisterLogic) Register(in *lender.UserRegisterRequest) (resp *lender.UserRegisterResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	// check if user already exists
	var existedUser *model.User
	existedUser, _ = l.svcCtx.UserRepo.GetLenderUserByEmail(l.ctx, in.Data.Email)

	if existedUser != nil && existedUser.Id > 0 {
		return nil, errorx.New("User Already Exists")
	}

	// check company
	isExisted, err := l.svcCtx.OrganizationRepo.IsExistedByColumns(l.ctx, map[string]interface{}{
		"id": in.Data.OrganizationId,
	})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	if !isExisted {
		return nil, errorx.New("Organization Is Not Exists")
	}

	var pwd string
	pwd, err = utils.HashPassword(in.Data.Password)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	userModel := &model.User{
		Email:                 in.Data.Email,
		PasswordHash:          pwd,
		Status:                constants.UserStatusActive,
		Type:                  constants.UserTypeLender,
		FirstName:             in.Data.FirstName,
		LastName:              in.Data.LastName,
		FullName:              stringx.FormatFullName(in.Data.FirstName, "", in.Data.LastName),
		DefaultOrganizationId: in.Data.OrganizationId,
	}

	var organization *model.Organization
	organization, err = l.svcCtx.OrganizationRepo.GetOneById(l.ctx, in.Data.OrganizationId, constants.UserTypeLender)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	userModel, err = l.svcCtx.UserRepo.RegisterWithOrganization(l.ctx, userModel, organization, constants.OrganizationAndUserStatusActive)
	if err != nil {
		return nil, err
	}
	jwtConfig := logic2.JwtConfig{
		AccessSecret:  l.svcCtx.Config.JwtAuth.AccessSecret,
		AccessExpire:  l.svcCtx.Config.JwtAuth.AccessExpire,
		RefreshSecret: l.svcCtx.Config.JwtAuth.RefreshSecret,
		RefreshExpire: l.svcCtx.Config.JwtAuth.RefreshExpire,
	}
	token, err := l.BaseLogic.UserLogic.GenerateJwtToken(userModel, jwtConfig, constants.UserTypeLender)
	if err != nil {
		return nil, err
	}
	data := &lender.UserRegisterResponseData{}
	err = copier.CopyWithOption(data, token, copier.Option{CaseSensitive: true})
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	return &lender.UserRegisterResponse{
		Data: data,
	}, err
}

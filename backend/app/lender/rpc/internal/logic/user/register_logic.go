package userlogic

import (
	"context"

	"github.com/jinzhu/copier"

	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/stringx"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.NewBaseLogic(ctx, svcCtx),
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

	var pwd string
	pwd, err = utils.HashPassword(in.Data.Password)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	userInput := dto.UserRegisterRequestDto{
		Email:                 in.Data.Email,
		PasswordHash:          pwd,
		Status:                constants.UserStatusActive,
		Type:                  constants.UserTypeLender,
		FirstName:             in.Data.FirstName,
		LastName:              in.Data.LastName,
		FullName:              stringx.FormatFullName(in.Data.FirstName, "", in.Data.LastName),
		DefaultOrganizationId: in.Data.OrganizationId,
	}

	var organizationInput *dto.OrganizationRegisterRequestDto
	switch in.Data.Type {
	case "Join":
		organization, err := l.svcCtx.OrganizationRepo.GetOrganization(l.ctx, in.Data.OrganizationId, constants.UserTypeLender, false)
		if err != nil {
			return nil, errorx.WithStack(err)
		}
		organizationInput, err = copierx.Trans[model.Organization, dto.OrganizationRegisterRequestDto](organization)
	case "Create":

		organizationInput = &dto.OrganizationRegisterRequestDto{
			Name:         in.Data.OrganizationName,
			Status:       constants.OrganizationStatusActive,
			OwnerType:    constants.OrganizationOwnerTypeThirdParty,
			Type:         constants.UserTypeLender,
			ContactEmail: in.Data.ContactEmail,
			ContactPhone: in.Data.ContactPhone,
			ContactName:  in.Data.ContactName,
			Logo:         in.Data.Logo,
			Website:      in.Data.Website,
		}
	default:
		err = errorx.New("Invalid Organization Type")
	}

	if err != nil {
		return nil, errorx.WithStack(err)
	}

	userModel, organizationUser, err := l.baseLogic.UserLogic.RegisterWithOrganization(l.ctx, userInput, *organizationInput, constants.UserTypeLender)
	if err != nil {
		return nil, err
	}

	token, err := l.baseLogic.UserLogic.GenerateJwtToken(userModel, organizationUser, constants.UserTypeLender)
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

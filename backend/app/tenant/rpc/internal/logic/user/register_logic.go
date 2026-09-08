package userlogic

import (
	"context"
	"strings"

	"github.com/zeromicro/go-zero/core/threading"

	"zeebee-backend/app/tenant/rpc/internal/logic"
	"zeebee-backend/app/tenant/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/stringx"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

// Register builds the identity (user + organization) only. Borrower pre-approval
// side effects stay at the borrower backend (it consumes user_id / organization_id
// from the response, and keeps the pre-approval JSON in its own request). Email
// is normalized to lowercase; an existing user gets the verify or login-guidance
// email resent instead of a duplicate-creation error.
func (l *RegisterLogic) Register(in *tenant.RegisterRequest) (resp *tenant.RegisterResponse, err error) {
	userType, err := logic.ParseUserType(in.Data.UserType)
	if err != nil {
		return nil, err
	}
	emailAddr := strings.ToLower(in.Data.Email)

	existedUser, _ := l.svcCtx.UserRepo.GetUserByEmail(l.ctx, emailAddr, userType)
	if existedUser != nil && existedUser.Id > 0 {
		threading.GoSafe(func() {
			ctx := context.WithoutCancel(l.ctx)
			if existedUser.Status == constants.UserStatusUnverified {
				if err := sendVerifyEmail(ctx, l.baseLogic, l.svcCtx, userType, emailAddr, existedUser.FirstName); err != nil {
					logSendError(l.svcCtx, ctx, "fail to resend register verify email", err)
				}
			} else {
				if err := sendLoginUrlEmail(ctx, l.baseLogic, l.svcCtx, userType, emailAddr, existedUser.FirstName); err != nil {
					logSendError(l.svcCtx, ctx, "fail to send login url email", err)
				}
			}
		})
		return l.registerResponse(existedUser, nil)
	}

	pwdHash, err := utils.HashPassword(in.Data.Password)
	if err != nil {
		return nil, err
	}

	userInput, err := copierx.Trans[tenant.RegisterRequestData, dto.UserRegisterRequestDto](in.Data)
	if err != nil {
		return nil, err
	}
	userInput.PasswordHash = pwdHash
	userInput.Status = constants.UserStatusUnverified
	userInput.Type = userType
	userInput.Email = emailAddr
	userInput.FullName = stringx.FormatFullName(in.Data.FirstName, "", in.Data.LastName)

	ownerType := constants.OrganizationOwnerType(in.Data.OwnerType)
	if userType == constants.UserTypeBorrower {
		ownerType = constants.OrganizationOwnerTypeNone
	} else if ownerType == "" {
		ownerType = constants.OrganizationOwnerTypeThirdParty
	}
	organizationInput := dto.OrganizationRegisterRequestDto{
		Name:      in.Data.OrganizationName,
		Type:      userType,
		OwnerType: ownerType,
	}

	userModel, organizationUser, err := l.baseLogic.UserLogic.RegisterWithOrganization(l.ctx, *userInput, organizationInput, userType, false)
	if err != nil {
		return nil, err
	}

	// Send the verify email asynchronously using the platform's SMTP client + brand config.
	threading.GoSafe(func() {
		ctx := context.WithoutCancel(l.ctx)
		if err := sendVerifyEmail(ctx, l.baseLogic, l.svcCtx, userType, emailAddr, userModel.FirstName); err != nil {
			logSendError(l.svcCtx, ctx, "fail to send register verify email", err)
		}
	})

	return l.registerResponse(userModel, organizationUser)
}

func (l *RegisterLogic) registerResponse(user *model.User, organizationUser *model.OrganizationUser) (*tenant.RegisterResponse, error) {
	if organizationUser == nil && user.DefaultOrganizationId > 0 {
		organizationUser, _ = l.svcCtx.OrganizationUserRepo.GetOrganizationUserByUserId(l.ctx, user.DefaultOrganizationId, user.Id)
	}
	accessToken := ""
	if organizationUser != nil {
		if token, err := l.baseLogic.UserLogic.GenerateJwtToken(user, organizationUser); err == nil {
			accessToken = token.AccessToken
		}
	}
	return &tenant.RegisterResponse{Data: &tenant.RegisterResponseData{
		UserId: user.Id, OrganizationId: user.DefaultOrganizationId, Status: string(user.Status), AccessToken: accessToken,
	}}, nil
}

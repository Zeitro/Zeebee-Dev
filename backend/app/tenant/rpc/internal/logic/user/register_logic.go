package userlogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"

	"zeebee-backend/app/tenant/rpc/internal/logic"
	"zeebee-backend/app/tenant/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/stringx"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
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
// from the response). The verify-email send is handled by SendVerifyEmail / the
// caller; tenant does not yet own per-platform email delivery.
func (l *RegisterLogic) Register(in *tenant.RegisterRequest) (resp *tenant.RegisterResponse, err error) {
	userType, err := logic.ParseUserType(in.Data.UserType)
	if err != nil {
		return nil, err
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
	userInput.FullName = stringx.FormatFullName(in.Data.FirstName, "", in.Data.LastName)

	organizationInput := dto.OrganizationRegisterRequestDto{
		Name:      in.Data.OrganizationName,
		Type:      userType,
		OwnerType: constants.OrganizationOwnerType(in.Data.OwnerType),
	}

	userModel, _, err := l.baseLogic.UserLogic.RegisterWithOrganization(l.ctx, *userInput, organizationInput, userType, false)
	if err != nil {
		return nil, err
	}

	// Send the verify email asynchronously using the platform's SMTP client + brand config.
	threading.GoSafe(func() {
		ctx := context.WithoutCancel(l.ctx)
		if err := sendVerifyEmail(ctx, l.baseLogic, l.svcCtx, userType, in.Data.Email, userModel.FirstName); err != nil {
			l.svcCtx.Logger.WithContext(ctx).
				WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).
				Error("fail to send register verify email")
		}
	})

	return &tenant.RegisterResponse{
		Data: &tenant.RegisterResponseData{
			UserId:         userModel.Id,
			OrganizationId: userModel.DefaultOrganizationId,
			Status:         string(userModel.Status),
		},
	}, nil
}

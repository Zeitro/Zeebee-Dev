package organizationlogic

import (
	"context"

	"gorm.io/gorm"

	"zeebee-backend/app/tenant/rpc/internal/logic"
	"zeebee-backend/app/tenant/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/stringx"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/proto"
)

type CreateLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{ctx: ctx, svcCtx: svcCtx, baseLogic: logic.GetBaseLogic(svcCtx)}
}

// Create handles both lender self-service creation and admin provisioning.
// The presence of owner credentials selects the admin-provisioning path.
func (l *CreateLogic) Create(in *tenant.OrganizationCreateRequest) (*tenant.OrganizationDetailResponse, error) {
	if in == nil || in.Data == nil {
		return nil, errorx.ErrorEmptyRequestData
	}
	userType, err := logic.ParseUserType(in.Data.UserType)
	if err != nil {
		return nil, err
	}
	if userType != constants.UserTypeLender {
		return nil, errorx.New("organization creation is only supported for lenders")
	}

	if in.Data.OwnerEmail != "" {
		return l.createWithOwner(in)
	}

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, errorx.Wrap(err, "fail to get user claim from context")
	}
	organization, err := copierx.Trans[tenant.OrganizationCreateRequestData, model.Organization](in.Data)
	if err != nil {
		return nil, err
	}
	organization.Status = constants.OrganizationStatusNotActivated
	organization.OwnerType = constants.OrganizationOwnerTypeThirdParty
	organization.Type = constants.UserTypeLender
	organization.OwnerId = claim.UserId

	exists, err := l.svcCtx.OrganizationRepo.IsExistedByColumns(l.ctx, map[string]interface{}{
		"name": organization.Name,
		"type": constants.UserTypeLender,
	})
	if err != nil {
		return nil, errorx.Wrap(err, "fail to check organization name uniqueness")
	}
	if exists {
		return nil, errorx.New("organization with this name already exists")
	}

	err = l.svcCtx.OrganizationRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		organization, err = l.svcCtx.OrganizationRepo.CreateWithTx(l.ctx, tx, organization, false)
		if err != nil {
			return errorx.Wrap(err, "fail to create organization")
		}
		firstName, midName, lastName := stringx.ParseFullName(claim.Username)
		_, err = l.svcCtx.OrganizationUserRepo.CreateWithTx(l.ctx, tx, &model.OrganizationUser{
			OrganizationId: organization.Id,
			UserId:         claim.UserId,
			Status:         constants.OrganizationUserStatusActive,
			FullName:       claim.Username,
			FirstName:      firstName,
			MidName:        midName,
			LastName:       lastName,
			Roles: []*model.Role{{
				Name:           constants.AdminRoleName,
				Type:           constants.UserTypeLender,
				OrganizationId: organization.Id,
			}},
		}, true)
		if err != nil {
			return err
		}
		return l.baseLogic.OrganizationTemplateLogic.SyncByOrganizationIdWithTx(l.ctx, tx, organization)
	})
	if err != nil {
		return nil, err
	}
	return l.detail(organization.Id)
}

func (l *CreateLogic) createWithOwner(in *tenant.OrganizationCreateRequest) (*tenant.OrganizationDetailResponse, error) {
	if in.Data.OwnerPassword == "" || in.Data.OwnerFirstName == "" || in.Data.OwnerLastName == "" {
		return nil, errorx.New("owner email, password, first name and last name are required")
	}
	pwdHash, err := utils.HashPassword(in.Data.OwnerPassword)
	if err != nil {
		return nil, err
	}
	userInput := dto.UserRegisterRequestDto{
		Email:        in.Data.OwnerEmail,
		PasswordHash: pwdHash,
		Status:       constants.UserStatusActive,
		Type:         constants.UserTypeLender,
		FirstName:    in.Data.OwnerFirstName,
		LastName:     in.Data.OwnerLastName,
		FullName:     stringx.FormatFullName(in.Data.OwnerFirstName, "", in.Data.OwnerLastName),
	}
	orgInput, err := copierx.Trans[tenant.OrganizationCreateRequestData, dto.OrganizationRegisterRequestDto](in.Data)
	if err != nil {
		return nil, err
	}
	orgInput.Status = constants.OrganizationStatusActive
	orgInput.Type = constants.UserTypeLender
	orgInput.OwnerType = constants.OrganizationOwnerTypeThirdParty
	if orgInput.Source == "" {
		orgInput.Source = constants.OrganizationSourceManual
	}
	if orgInput.Order == 0 {
		orgInput.Order, err = l.svcCtx.OrganizationRepo.MaxOrderByType(l.ctx, constants.UserTypeLender)
		if err != nil {
			return nil, err
		}
		orgInput.Order++
	}
	_, organizationUser, err := l.baseLogic.UserLogic.RegisterWithOrganization(l.ctx, userInput, *orgInput, constants.UserTypeLender, true)
	if err != nil {
		return nil, err
	}
	_ = l.baseLogic.OrganizationSubscriptionLogic.SubscribeDefaultPlan(l.ctx, organizationUser.OrganizationId, constants.UserTypeLender)
	return l.detail(organizationUser.OrganizationId)
}

func (l *CreateLogic) detail(id int64) (*tenant.OrganizationDetailResponse, error) {
	organization, err := l.baseLogic.OrganizationLogic.DetailWithType(l.ctx, id, constants.UserTypeLender)
	if err != nil {
		return nil, err
	}
	data, err := copierx.Trans[dto.OrganizationEntityDto, proto.OrganizationEntity](organization, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}
	return &tenant.OrganizationDetailResponse{Data: data}, nil
}

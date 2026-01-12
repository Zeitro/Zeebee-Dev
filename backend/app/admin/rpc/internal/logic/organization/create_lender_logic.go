package organizationlogic

import (
	"context"

	"gorm.io/gorm"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/stringx"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLenderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	*logic.BaseLogic
}

func NewCreateLenderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLenderLogic {
	return &CreateLenderLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		BaseLogic: logic.NewBaseLogic(ctx, svcCtx),
	}
}

// CreateLender creates a new Lender organization with owner user.
// Key differences from RegisterWithOrganization:
// - All statuses set to Active (no activation flow)
// - Validates organization name uniqueness (Lender type)
// - Validates owner email uniqueness (Lender type)
func (l *CreateLenderLogic) CreateLender(in *admin.OrganizationCreateLenderRequest) (resp *admin.OrganizationDetailResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	if in.Data == nil || in.Data.User == nil || in.Data.Organization == nil {
		return nil, errorx.New("Request Data Is Required")
	}

	// Validate organization name uniqueness (Lender type)
	isOrgExisted, err := l.svcCtx.AllRepos.OrganizationRepo.IsExistedByColumns(l.ctx, map[string]interface{}{
		"name": in.Data.Organization.Name,
		"type": constants.UserTypeLender,
	})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	if isOrgExisted {
		return nil, errorx.New("Organization Name Already Exists")
	}

	// Validate owner email uniqueness (Lender type)
	isUserExisted, err := l.svcCtx.AllRepos.UserRepo.IsExistedByColumns(l.ctx, map[string]interface{}{
		"email": in.Data.User.Email,
		"type":  constants.UserTypeLender,
	})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	if isUserExisted {
		return nil, errorx.New("Owner Email Already Exists")
	}

	// Hash password
	pwdHash, err := utils.HashPassword(in.Data.User.Password)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	// Prepare user input
	userInput := dto.UserRegisterRequestDto{
		Email:        in.Data.User.Email,
		PasswordHash: pwdHash,
		Status:       constants.UserStatusActive, // Admin created: Active directly
		Type:         constants.UserTypeLender,
		FirstName:    in.Data.User.FirstName,
		LastName:     in.Data.User.LastName,
		FullName:     stringx.FormatFullName(in.Data.User.FirstName, "", in.Data.User.LastName),
	}

	// Prepare organization input
	orgInput := dto.OrganizationRegisterRequestDto{
		Name:         in.Data.Organization.Name,
		Type:         constants.UserTypeLender,
		Status:       constants.OrganizationStatusActive, // Admin created: Active directly
		ContactEmail: in.Data.Organization.ContactEmail,
		ContactPhone: in.Data.Organization.ContactPhone,
		ContactName:  in.Data.Organization.ContactName,
		Logo:         in.Data.Organization.Logo,
		Website:      in.Data.Organization.Website,
		OwnerType:    constants.OrganizationOwnerTypeSelf,
	}

	var organization *model.Organization
	var user *model.User
	var organizationUser *model.OrganizationUser

	err = l.svcCtx.AllRepos.UserRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		var innerErr error

		// Create organization
		organization, innerErr = copierx.Trans[dto.OrganizationRegisterRequestDto, model.Organization](&orgInput)
		if innerErr != nil {
			return innerErr
		}
		// Set additional fields not in DTO
		organization.Tagline = in.Data.Organization.Tagline
		if in.Data.Organization.LenderMetadata != nil {
			organization.Metadata = &model.OrganizationMetadata{
				Lender: &model.OrganizationLenderMetadata{
					LoanAmountRange: model.RangeFloat64{
						Min: in.Data.Organization.LenderMetadata.LoanAmountRange.Min,
						Max: in.Data.Organization.LenderMetadata.LoanAmountRange.Max,
					},
					LoanTermRange: model.RangeInt{
						Min: int(in.Data.Organization.LenderMetadata.LoanTermRange.Min),
						Max: int(in.Data.Organization.LenderMetadata.LoanTermRange.Max),
					},
					LoanRateRange: model.RangeFloat64{
						Min: in.Data.Organization.LenderMetadata.LoanRateRange.Min,
						Max: in.Data.Organization.LenderMetadata.LoanRateRange.Max,
					},
					LoanPointRange: model.RangeFloat64{
						Min: in.Data.Organization.LenderMetadata.LoanPointRange.Min,
						Max: in.Data.Organization.LenderMetadata.LoanPointRange.Max,
					},
				},
			}
		}
		organization, innerErr = l.svcCtx.AllRepos.OrganizationRepo.CreateWithTx(l.ctx, tx, organization, false)
		if innerErr != nil {
			return innerErr
		}

		// Create user
		user, innerErr = copierx.Trans[dto.UserRegisterRequestDto, model.User](&userInput)
		if innerErr != nil {
			return innerErr
		}
		user.DefaultOrganizationId = organization.Id
		user.RemoveRelationships()
		user, innerErr = l.svcCtx.AllRepos.UserRepo.CreateWithTx(l.ctx, tx, user, false)
		if innerErr != nil {
			return innerErr
		}

		// Update organization's owner
		innerErr = l.svcCtx.AllRepos.OrganizationRepo.UpdateOwnerIdWithTx(l.ctx, tx, organization.Id, user.Id)
		if innerErr != nil {
			return innerErr
		}
		organization.OwnerId = user.Id

		// Create Admin role for the organization
		roles := []*model.Role{
			{
				Name:           constants.AdminRoleName,
				Type:           constants.UserTypeLender,
				OrganizationId: organization.Id,
			},
		}

		// Create organization user relation
		organizationUser = &model.OrganizationUser{
			OrganizationId: organization.Id,
			UserId:         user.Id,
			FirstName:      user.FirstName,
			MidName:        user.MidName,
			LastName:       user.LastName,
			FullName:       user.FullName,
			Status:         constants.OrganizationUserStatusActive, // Admin created: Active directly
			Roles:          roles,
		}
		organizationUser, innerErr = l.svcCtx.AllRepos.OrganizationUserRepo.CreateWithTx(l.ctx, tx, organizationUser, true)
		if innerErr != nil {
			return innerErr
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Load organization with relations for response
	organization, err = l.svcCtx.AllRepos.OrganizationRepo.GetOrganization(l.ctx, organization.Id, constants.UserTypeLender, true)
	if err != nil {
		return nil, err
	}

	// Convert to response
	data, err := copierx.Trans[*model.Organization, admin.OrganizationEntity](&organization, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	return &admin.OrganizationDetailResponse{
		Data: data,
	}, nil
}

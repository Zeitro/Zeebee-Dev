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
	"zeebee-backend/domains/zeebee/proto"
)

type UpdateLenderLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewUpdateLenderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLenderLogic {
	return &UpdateLenderLogic{ctx: ctx, svcCtx: svcCtx, baseLogic: logic.GetBaseLogic(svcCtx)}
}

func (l *UpdateLenderLogic) UpdateLender(in *tenant.OrganizationUpdateLenderRequest) (*tenant.OrganizationDetailResponse, error) {
	if in == nil || in.Data == nil {
		return nil, errorx.ErrorEmptyRequestData
	}
	if _, err := logic.ParseUserType(in.Data.UserType); err != nil {
		return nil, err
	}
	organization, err := l.svcCtx.OrganizationRepo.GetOrganization(l.ctx, in.Data.Id, constants.UserTypeLender, false)
	if err != nil {
		return nil, errorx.Wrap(err, "fail to get organization")
	}
	needOwner := organization.Status == constants.OrganizationStatusNotActivated && organization.OwnerId == 0 && organization.Source == constants.OrganizationSourceSearch
	if needOwner {
		if in.Data.OwnerEmail == "" || in.Data.OwnerPassword == "" || in.Data.OwnerFirstName == "" || in.Data.OwnerLastName == "" {
			return nil, errorx.New("owner account (email, password, name) is required to activate this lender")
		}
		if !l.baseLogic.UserLogic.ValidatePortalPassword(in.Data.OwnerPassword) {
			return nil, errorx.New("password must be at least 8 characters and include upper, lower, and numeric characters")
		}
	}

	attachedOwner := false
	err = l.svcCtx.OrganizationRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		org, innerErr := l.svcCtx.OrganizationRepo.GetOrganizationWithTx(l.ctx, tx, in.Data.Id, constants.UserTypeLender, false)
		if innerErr != nil {
			return innerErr
		}
		wasNA := org.Status == constants.OrganizationStatusNotActivated
		needOwner = wasNA && org.OwnerId == 0 && org.Source == constants.OrganizationSourceSearch
		if needOwner {
			pwdHash, hashErr := utils.HashPassword(in.Data.OwnerPassword)
			if hashErr != nil {
				return hashErr
			}
			innerErr = l.baseLogic.UserLogic.CreatePrimaryOwnerForLenderOrgWithTx(l.ctx, tx, org, dto.UserRegisterRequestDto{
				Email: in.Data.OwnerEmail, PasswordHash: pwdHash, Status: constants.UserStatusActive,
				Type: constants.UserTypeLender, FirstName: in.Data.OwnerFirstName, LastName: in.Data.OwnerLastName,
				FullName: stringx.FormatFullName(in.Data.OwnerFirstName, "", in.Data.OwnerLastName),
			})
			if innerErr != nil {
				return innerErr
			}
			attachedOwner = true
		}
		preservedSource, preservedName := org.Source, org.Name
		if innerErr = copierx.Assign(in.Data, org); innerErr != nil {
			return innerErr
		}
		org.Source, org.Name = preservedSource, preservedName
		if wasNA {
			org.Status = constants.OrganizationStatusActive
			org.Order, innerErr = l.svcCtx.OrganizationRepo.MaxOrderByTypeWithTx(l.ctx, tx, constants.UserTypeLender)
			if innerErr != nil {
				return innerErr
			}
			org.Order++
		}
		_, innerErr = l.svcCtx.OrganizationRepo.UpdateOneByIdWithTx(l.ctx, tx, in.Data.Id, org)
		return innerErr
	})
	if err != nil {
		return nil, err
	}
	if attachedOwner {
		l.baseLogic.UserLogic.PublishRoleReloadForOrganization(l.ctx, in.Data.Id, constants.UserTypeLender)
	}
	result, err := l.baseLogic.OrganizationLogic.DetailWithType(l.ctx, in.Data.Id, constants.UserTypeLender)
	if err != nil {
		return nil, err
	}
	data, err := copierx.Trans[dto.OrganizationEntityDto, proto.OrganizationEntity](result, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}
	return &tenant.OrganizationDetailResponse{Data: data}, nil
}

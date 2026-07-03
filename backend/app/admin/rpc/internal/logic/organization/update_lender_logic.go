package organizationlogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/protobuf/proto"
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
	zeebeeproto "zeebee-backend/domains/zeebee/proto"
)

type UpdateLenderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewUpdateLenderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLenderLogic {
	return &UpdateLenderLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

// UpdateLender updates editable fields of a Lender organization.
// When activating a search-sourced org without an owner, request must include owner credentials.
func (l *UpdateLenderLogic) UpdateLender(in *admin.OrganizationUpdateLenderRequest) (resp *admin.OrganizationDetailResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
				logx.Field("service_name", l.svcCtx.Config.Name),
			).Error(err)
		}
	}()

	organization, err := l.svcCtx.OrganizationRepo.GetOrganization(l.ctx, in.Data.Id, constants.UserTypeLender, false)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	wasNotActivated := organization.Status == constants.OrganizationStatusNotActivated
	needOwnerOnboard := wasNotActivated &&
		organization.OwnerId == 0 &&
		organization.Source == constants.OrganizationSourceSearch

	if needOwnerOnboard {
		o := in.Data.GetOwner()
		if o == nil {
			return nil, errorx.New("Owner account (email, password, name) is required to activate this lender")
		}
		if !l.baseLogic.UserLogic.ValidatePortalPassword(o.GetPassword()) {
			return nil, errorx.New("Password must be at least 8 characters and include upper, lower, and numeric characters")
		}
	}

	var attachedOwner bool
	err = l.svcCtx.OrganizationRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		org, innerErr := l.svcCtx.OrganizationRepo.GetOrganizationWithTx(l.ctx, tx, in.Data.Id, constants.UserTypeLender, false)
		if innerErr != nil {
			return innerErr
		}
		wasNA := org.Status == constants.OrganizationStatusNotActivated
		needOwner := wasNA &&
			org.OwnerId == 0 &&
			org.Source == constants.OrganizationSourceSearch

		if needOwner {
			o := in.Data.GetOwner()
			if o == nil {
				return errorx.New("Owner account is required to activate this lender")
			}
			var pwdHash string
			pwdHash, innerErr = utils.HashPassword(o.GetPassword())
			if innerErr != nil {
				return errorx.WithStack(innerErr)
			}
			userDto := dto.UserRegisterRequestDto{
				Email:        o.GetEmail(),
				PasswordHash: pwdHash,
				Status:       constants.UserStatusActive,
				Type:         constants.UserTypeLender,
				FirstName:    o.GetFirstName(),
				LastName:     o.GetLastName(),
				FullName:     stringx.FormatFullName(o.GetFirstName(), "", o.GetLastName()),
			}
			innerErr = l.baseLogic.UserLogic.CreatePrimaryOwnerForLenderOrgWithTx(l.ctx, tx, org, userDto)
			if innerErr != nil {
				return innerErr
			}
			attachedOwner = true
		}

		assignSrc := proto.Clone(in.Data).(*admin.OrganizationUpdateLenderRequestData)
		assignSrc.Owner = nil
		preservedSource, preservedName := org.Source, org.Name
		innerErr = copierx.Assign(assignSrc, org)
		if innerErr != nil {
			return innerErr
		}
		org.Source = preservedSource
		org.Name = preservedName
		if wasNA {
			org.Status = constants.OrganizationStatusActive
			// Append the newly activated lender at the end of the lender list.
			var maxOrder int
			maxOrder, innerErr = l.svcCtx.OrganizationRepo.MaxOrderByTypeWithTx(l.ctx, tx, constants.UserTypeLender)
			if innerErr != nil {
				return innerErr
			}
			org.Order = maxOrder + 1
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
	// Update the model
	_, err = l.svcCtx.OrganizationRepo.UpdateOneById(l.ctx, in.Data.Id, organization)
	if err != nil {
		return nil, errorx.Wrap(err, "fail to update organization")
	}

	organizationBo, err := l.baseLogic.OrganizationLogic.Detail(l.ctx, in.Data.Id)
	if err != nil {
		return nil, err
	}
	data, err := copierx.Trans[dto.OrganizationEntityDto, zeebeeproto.OrganizationEntity](organizationBo, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}

	return &admin.OrganizationDetailResponse{
		Data: data,
	}, nil
}

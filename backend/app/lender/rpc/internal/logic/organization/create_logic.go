package organizationlogic

import (
	"context"

	"gorm.io/gorm"

	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/stringx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type CreateLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
	logx.Logger
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.NewBaseLogic(ctx, svcCtx),
		Logger:    logx.WithContext(ctx),
	}
}

func (l *CreateLogic) Create(in *lender.OrganizationCreateRequest) (resp *lender.OrganizationDetailResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	// Get current user claim (creator will become owner)
	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, errorx.Wrap(err, "Failed To Get User Claim From Context")
	}

	// Convert request data to organization model
	organization := &model.Organization{}
	organization, err = copierx.Trans[lender.OrganizationCreateRequestData, model.Organization](
		in.Data,
	)
	if err != nil {
		return nil, err
	}

	// Set required default values for new organization
	organization.Status = constants.OrganizationStatusNotActivated
	organization.OwnerType = constants.OrganizationOwnerTypeThirdParty
	organization.Type = constants.UserTypeLender
	organization.OwnerId = claim.UserId

	// Check if organization with the same name already exists
	exists, err := l.svcCtx.OrganizationRepo.IsExistedByColumns(l.ctx, map[string]interface{}{
		"name": organization.Name,
		"type": constants.UserTypeLender,
	})
	if err != nil {
		return nil, errorx.Wrap(err, "Failed To Check Organization Name Uniqueness")
	}
	if exists {
		return nil, errorx.New("Organization With This Name Already Exists")
	}

	// Use transaction to ensure data consistency across multiple operations
	err = l.svcCtx.OrganizationRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		// Create organization record
		organization, err = l.svcCtx.OrganizationRepo.CreateWithTx(l.ctx, tx, organization, false)
		if err != nil {
			return errorx.Wrap(err, "Failed To Create Organization")
		}

		// Create admin role for the organization
		adminRole := &model.Role{
			Name:           constants.AdminRoleName,
			Type:           constants.UserTypeLender,
			OrganizationId: organization.Id,
		}

		// Create organization user association with admin role
		firstName, midName, lastName := stringx.ParseFullName(claim.Username)
		organizationUser := &model.OrganizationUser{
			OrganizationId: organization.Id,
			UserId:         claim.UserId,
			Status:         constants.OrganizationUserStatusActive,
			FullName:       claim.Username,
			FirstName:      firstName,
			MidName:        midName,
			LastName:       lastName,
			Roles:          []*model.Role{adminRole},
		}

		_, err = l.svcCtx.OrganizationUserRepo.CreateWithTx(l.ctx, tx, organizationUser, true)
		if err != nil {
			return errorx.Wrap(err, "Failed To Create Organization User")
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Convert created organization to proto entity for response
	resultOrganization, err := copierx.Trans[model.Organization, lender.OrganizationEntity](
		organization,
		copierx.TimeToRFC3339Str(),
	)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	return &lender.OrganizationDetailResponse{
		Data: resultOrganization,
	}, nil
}

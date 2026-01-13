package organizationlogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/stringx"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/bo"
	"zeebee-backend/domains/zeebee/data/dto"
)

type CreateLenderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewCreateLenderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLenderLogic {
	return &CreateLenderLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.NewBaseLogic(ctx, svcCtx),
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

	// Hash password
	pwdHash, err := utils.HashPassword(in.Data.User.Password)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	// Prepare user input
	userInput, err := copierx.Trans[admin.OrganizationCreateLenderUserEntity, dto.UserRegisterRequestDto](in.Data.User)
	if err != nil {
		return nil, err
	}
	userInput.Status = constants.UserStatusActive
	userInput.Type = constants.UserTypeLender
	userInput.PasswordHash = pwdHash
	userInput.FullName = stringx.FormatFullName(in.Data.User.FirstName, "", in.Data.User.LastName)

	// Prepare organization input
	lenderMetadata, err := copierx.Trans[admin.OrganizationLenderMetadataEntity, dto.OrganizationLenderMetadata](in.Data.Organization.LenderMetadata)
	if err != nil {
		return nil, err
	}

	orgInput, err := copierx.Trans[admin.OrganizationCreateLenderOrganizationEntity, dto.OrganizationRegisterRequestDto](in.Data.Organization)
	if err != nil {
		return nil, err
	}
	orgInput.Status = constants.OrganizationStatusActive
	orgInput.Type = constants.UserTypeLender
	orgInput.OwnerType = constants.OrganizationOwnerTypeThirdParty
	orgInput.Metadata = &dto.OrganizationMetadata{
		Lender: lenderMetadata,
	}
	// Register with organization
	_, organizationUser, err := l.baseLogic.UserLogic.RegisterWithOrganization(l.ctx, *userInput, *orgInput, constants.UserTypeLender, true)
	if err != nil {
		return nil, err
	}

	// Load organization with relations for response
	organization, err := l.baseLogic.OrganizationLogic.DetailWithType(l.ctx, organizationUser.OrganizationId, constants.UserTypeLender)
	if err != nil {
		return nil, err
	}

	// Convert to response
	data, err := copierx.Trans[bo.OrganizationBo, admin.OrganizationEntity](organization, copierx.TimeToRFC3339Str())
	return &admin.OrganizationDetailResponse{
		Data: data,
	}, err
}

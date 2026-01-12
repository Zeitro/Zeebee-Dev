package organizationlogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLenderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	*logic.BaseLogic
}

func NewUpdateLenderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLenderLogic {
	return &UpdateLenderLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		BaseLogic: logic.NewBaseLogic(ctx, svcCtx),
	}
}

// UpdateLender updates editable fields of a Lender organization.
// Organization name cannot be modified (not included in the request).
// Only updates: ContactEmail, ContactPhone, ContactName, Logo, Website, Tagline, LenderMetadata
func (l *UpdateLenderLogic) UpdateLender(in *admin.OrganizationUpdateLenderRequest) (resp *admin.OrganizationDetailResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	// Get existing organization
	organization, err := l.svcCtx.AllRepos.OrganizationRepo.GetOrganization(l.ctx, in.Data.Id, constants.UserTypeLender, false)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	// Update editable fields
	organization.ContactEmail = in.Data.ContactEmail
	organization.ContactPhone = in.Data.ContactPhone
	organization.ContactName = in.Data.ContactName
	organization.Logo = in.Data.Logo
	organization.Website = in.Data.Website
	organization.Tagline = in.Data.Tagline

	// Update LenderMetadata if provided
	if in.Data.LenderMetadata != nil {
		if organization.Metadata == nil {
			organization.Metadata = &model.OrganizationMetadata{}
		}
		if organization.Metadata.Lender == nil {
			organization.Metadata.Lender = &model.OrganizationLenderMetadata{}
		}
		organization.Metadata.Lender.LoanAmountRange = model.RangeFloat64{
			Min: in.Data.LenderMetadata.LoanAmountRange.Min,
			Max: in.Data.LenderMetadata.LoanAmountRange.Max,
		}
		organization.Metadata.Lender.LoanTermRange = model.RangeInt{
			Min: int(in.Data.LenderMetadata.LoanTermRange.Min),
			Max: int(in.Data.LenderMetadata.LoanTermRange.Max),
		}
		organization.Metadata.Lender.LoanRateRange = model.RangeFloat64{
			Min: in.Data.LenderMetadata.LoanRateRange.Min,
			Max: in.Data.LenderMetadata.LoanRateRange.Max,
		}
		organization.Metadata.Lender.LoanPointRange = model.RangeFloat64{
			Min: in.Data.LenderMetadata.LoanPointRange.Min,
			Max: in.Data.LenderMetadata.LoanPointRange.Max,
		}
	}

	// Save update using BaseRepo UpdateOneById
	_, err = l.svcCtx.AllRepos.OrganizationRepo.UpdateOneById(l.ctx, organization.Id, organization)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	// Reload with relations for response
	organization, err = l.svcCtx.AllRepos.OrganizationRepo.GetOrganization(l.ctx, in.Data.Id, constants.UserTypeLender, true)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	// Update Owner info from OrganizationUser
	if organization.Owner != nil {
		organizationUser, err := l.svcCtx.AllRepos.OrganizationUserRepo.GetOrganizationUserByUserId(l.ctx, organization.Id, organization.Owner.Id)
		if err == nil && organizationUser != nil {
			if len(organizationUser.FullName) > 0 {
				organization.Owner.FullName = organizationUser.FullName
				organization.Owner.FirstName = organizationUser.FirstName
				organization.Owner.LastName = organizationUser.LastName
			}
		}
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

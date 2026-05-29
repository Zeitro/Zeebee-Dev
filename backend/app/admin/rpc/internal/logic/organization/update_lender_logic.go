package organizationlogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
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
// Organization name cannot be modified (not included in the request).
// Only updates: ContactEmail, ContactPhone, ContactName, Logo, Website, Tagline, LenderMetadata
func (l *UpdateLenderLogic) UpdateLender(in *admin.OrganizationUpdateLenderRequest) (resp *admin.OrganizationDetailResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	// Get existing organization
	organization, err := l.svcCtx.OrganizationRepo.GetOrganization(l.ctx, in.Data.Id, constants.UserTypeLender, false)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	err = copierx.Assign[admin.OrganizationUpdateLenderRequestData, model.Organization](in.Data, organization)
	if err != nil {
		return nil, err
	}

	// Update the model
	_, err = l.svcCtx.OrganizationRepo.UpdateOneById(l.ctx, in.Data.Id, organization)
	if err != nil {
		return nil, errorx.Wrap(err, "Failed to update organization")
	}
	// Get updated organization
	organizationBo, err := l.baseLogic.OrganizationLogic.Detail(l.ctx, in.Data.Id)
	if err != nil {
		return nil, err
	}
	// Convert to response
	data, err := copierx.Trans[dto.OrganizationEntityDto, admin.OrganizationEntity](organizationBo, copierx.TimeToRFC3339Str())

	return &admin.OrganizationDetailResponse{
		Data: data,
	}, nil
}

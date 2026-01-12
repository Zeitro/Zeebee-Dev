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

type DetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	*logic.BaseLogic
}

func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		BaseLogic: logic.NewBaseLogic(ctx, svcCtx),
	}
}

// Detail retrieves organization details by ID.
// It fetches the organization with preloaded relations (Owner, OrganizationEmails).
func (l *DetailLogic) Detail(in *admin.IdRequest) (resp *admin.OrganizationDetailResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	// Get organization with preloaded relations
	organization, err := l.svcCtx.AllRepos.OrganizationRepo.GetOrganization(l.ctx, in.Data, constants.UserTypeLender, true)
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

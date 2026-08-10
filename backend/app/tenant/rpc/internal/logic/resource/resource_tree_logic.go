package resourcelogic

import (
	"context"

	"zeebee-backend/app/tenant/rpc/internal/logic"
	"zeebee-backend/app/tenant/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/proto"
)

type ResourceTreeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewResourceTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResourceTreeLogic {
	return &ResourceTreeLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

// ResourceTree returns the resource tree of a platform, either global (all) or
// scoped to one organization.
func (l *ResourceTreeLogic) ResourceTree(in *tenant.ResourceTreeRequest) (resp *tenant.ResourceTreeResponse, err error) {
	if in.Data == nil {
		return nil, errorx.ErrorEmptyRequestData
	}
	userType, err := logic.ParseUserType(in.Data.UserType)
	if err != nil {
		return nil, err
	}
	targetType, err := logic.ParseUserType(in.Data.Type)
	if err != nil {
		return nil, err
	}
	if userType != constants.UserTypeAdmin && targetType != userType {
		return nil, errorx.New("resource type must match caller platform")
	}

	var trees []*dto.ResourceTreeDto
	switch constants.ResourceTreeScope(in.Data.Scope) {
	case constants.ResourceTreeScopeAll:
		trees, err = l.baseLogic.ResourceLogic.ResourceTree(l.ctx, targetType)
	case constants.ResourceTreeScopeOrganization:
		trees, err = l.baseLogic.ResourceLogic.GetOrganizationResources(l.ctx, in.Data.OrganizationId, targetType)
	default:
		return nil, errorx.New("invalid resource tree scope")
	}
	if err != nil {
		return nil, err
	}

	data, err := copierx.Trans[[]*dto.ResourceTreeDto, []*proto.ResourceTreeEntity](&trees, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}
	return &tenant.ResourceTreeResponse{Data: *data}, nil
}

package resourcelogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
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
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *ResourceTreeLogic) ResourceTree(in *admin.ResourceTreeRequest) (resp *admin.ResourceTreeResponse, err error) {

	var trees []*dto.ResourceTreeDto
	switch constants.ResourceTreeScope(in.Data.Scope) {
	case constants.ResourceTreeScopeAll:
		trees, err = l.baseLogic.ResourceLogic.ResourceTree(l.ctx, constants.UserType(in.Data.Type))
	case constants.ResourceTreeScopeOrganization:
		claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
		if err != nil {
			return nil, err
		}
		trees, err = l.baseLogic.ResourceLogic.GetOrganizationResources(l.ctx, claim.OrganizationId, constants.UserType(in.Data.Type))
	}

	if err != nil {
		return nil, err
	}

	data, err := copierx.Trans[[]*dto.ResourceTreeDto, []*proto.ResourceTreeEntity](&trees, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}

	return &admin.ResourceTreeResponse{
		Data: *data,
	}, nil
}

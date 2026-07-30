package resourcelogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
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

func (l *ResourceTreeLogic) ResourceTree(in *borrower.ResourceTreeRequest) (resp *borrower.ResourceTreeResponse, err error) {

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	trees, err := l.baseLogic.ResourceLogic.GetOrganizationResources(l.ctx, claim.OrganizationId, constants.UserType(in.Data.Type))
	if err != nil {
		return nil, err
	}

	data, err := copierx.Trans[[]*dto.ResourceTreeDto, []*proto.ResourceTreeEntity](&trees, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}

	return &borrower.ResourceTreeResponse{
		Data: *data,
	}, nil
}

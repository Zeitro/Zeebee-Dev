package organizationlogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/proto"
)

type DetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

// Detail retrieves organization details by ID.
// It fetches the organization with preloaded relations (Owner, OrganizationEmails).
func (l *DetailLogic) Detail(in *proto.IdRequest) (resp *admin.OrganizationDetailResponse, err error) {

	organization, err := l.baseLogic.OrganizationLogic.Detail(l.ctx, in.Data)
	if err != nil {
		return nil, err
	}
	// Convert to response
	data, err := copierx.Trans[dto.OrganizationEntityDto, proto.OrganizationEntity](organization, copierx.TimeToRFC3339Str())

	return &admin.OrganizationDetailResponse{
		Data: data,
	}, err
}

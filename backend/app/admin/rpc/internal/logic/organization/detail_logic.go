package organizationlogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/data/dto"
)

type DetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.NewBaseLogic(ctx, svcCtx),
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

	organization, err := l.baseLogic.OrganizationLogic.Detail(l.ctx, in.Data)
	if err != nil {
		return nil, err
	}
	// Convert to response
	data, err := copierx.Trans[dto.OrganizationEntityDto, admin.OrganizationEntity](organization, copierx.TimeToRFC3339Str())

	return &admin.OrganizationDetailResponse{
		Data: data,
	}, err
}

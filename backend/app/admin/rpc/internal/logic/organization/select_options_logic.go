package organizationlogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/bo"

	"zeebee-backend/common/errorx"
)

type SelectOptionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewSelectOptionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SelectOptionsLogic {
	return &SelectOptionsLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.NewBaseLogic(ctx, svcCtx),
	}
}

func (l *SelectOptionsLogic) SelectOptions(in *admin.OrganizationSelectOptionsRequest) (resp *admin.OrganizationSelectOptionsResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	if in.Data == nil {
		return nil, errorx.New("Data is required")
	}

	organizationType := constants.UserType(in.Data.Type)

	organization, err := l.baseLogic.OrganizationLogic.AllForOption(l.ctx, organizationType, constants.OrganizationStatusActive)
	if err != nil {
		return nil, err
	}
	data, err := copierx.Trans[[]*bo.OrganizationOptionBo, []*admin.OrganizationOptionsEntity](&organization)
	if err != nil {
		return nil, err
	}

	return &admin.OrganizationSelectOptionsResponse{
		Data: *data,
	}, nil
}

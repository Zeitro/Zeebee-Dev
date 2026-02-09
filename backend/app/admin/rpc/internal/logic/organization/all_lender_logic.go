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
	"zeebee-backend/domains/zeebee/data/bo"
)

type AllLenderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewAllLenderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AllLenderLogic {
	return &AllLenderLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.NewBaseLogic(ctx, svcCtx),
	}
}

func (l *AllLenderLogic) AllLender(in *admin.EmptyRequest) (resp *admin.OrganizationSelectOptionsResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	allLender, err := l.baseLogic.OrganizationLogic.AllForOption(l.ctx, constants.UserTypeLender, constants.OrganizationStatusNone)
	if err != nil {
		return nil, err
	}
	data, err := copierx.Trans[[]*bo.OrganizationOptionBo, []*admin.OrganizationOptionsEntity](&allLender)
	if err != nil {
		return nil, err
	}

	return &admin.OrganizationSelectOptionsResponse{
		Data: *data,
	}, nil
}

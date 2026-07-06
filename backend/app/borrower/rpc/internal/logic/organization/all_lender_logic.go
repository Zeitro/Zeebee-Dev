package organizationlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/proto"

	"github.com/zeromicro/go-zero/core/logx"
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
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *AllLenderLogic) AllLender(in *proto.EmptyRequest) (*borrower.OrganizationAllLenderResponse, error) {
	var err error
	var allLender []*dto.OrganizationOptionDto
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
				logx.Field("service_name", l.svcCtx.Config.Name),
			).Error(err)
		}
	}()
	allLender, err = l.baseLogic.OrganizationLogic.AllForOption(l.ctx, constants.UserTypeLender, constants.OrganizationStatusActive)
	if err != nil {
		return nil, err
	}
	data, err := copierx.Trans[[]*dto.OrganizationOptionDto, []*proto.OrganizationOptionsEntity](&allLender)
	if err != nil {
		return nil, err
	}

	return &borrower.OrganizationAllLenderResponse{
		Data: *data,
	}, nil
}

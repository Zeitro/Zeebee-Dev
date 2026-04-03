package resourcelogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/data/dto"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type UserMenuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewUserMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserMenuLogic {
	return &UserMenuLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *UserMenuLogic) UserMenu(in *borrower.EmptyRequest) (resp *borrower.ResourceTreeResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	input := dto.ResourceUserMenuRequestDto{
		OrganizationId:     claim.OrganizationId,
		OrganizationUserId: claim.OrganizationUserId,
		Type:               claim.UserType,
	}

	result, err := l.baseLogic.ResourceLogic.GetUserResources(l.ctx, input)
	if err != nil {
		return nil, err
	}
	// trans response
	resp = &borrower.ResourceTreeResponse{}
	menus, err := copierx.Trans[[]*dto.ResourceTreeDto, []*borrower.ResourceTreeEntity](&result)

	return &borrower.ResourceTreeResponse{
		Data: *menus,
	}, nil
}

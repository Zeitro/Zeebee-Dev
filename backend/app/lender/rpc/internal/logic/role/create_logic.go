package rolelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/bo"
	"zeebee-backend/domains/zeebee/data/dto"
)

type CreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.NewBaseLogic(ctx, svcCtx),
	}
}

func (l *CreateLogic) Create(in *lender.RoleCreateRequest) (resp *lender.RoleCreateResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	input, err := copierx.Trans[lender.RoleCreateRequestData, dto.RoleCreateRequestDto](in.Data)
	if err != nil {
		return nil, err
	}
	input.Type = constants.UserTypeLender
	input.OrganizationId = claim.OrganizationId

	role, err := l.baseLogic.RoleLogic.Create(l.ctx, *input)
	if err != nil {
		return nil, err
	}

	data, err := copierx.Trans[bo.RoleBo, lender.RoleEntity](role, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}

	return &lender.RoleCreateResponse{
		Data: data,
	}, nil
}

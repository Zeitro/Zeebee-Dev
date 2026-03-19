package rolelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
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
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *CreateLogic) Create(in *borrower.RoleCreateRequest) (resp *borrower.RoleCreateResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	input := &dto.RoleCreateRequestDto{}
	input, err = copierx.Trans[borrower.RoleCreateRequestData, dto.RoleCreateRequestDto](in.Data)
	if err != nil {
		return nil, err
	}

	input.Type = constants.UserTypeBorrower
	input.OrganizationId = claim.OrganizationId

	role, err := l.baseLogic.RoleLogic.Create(l.ctx, *input)
	if err != nil {
		return nil, err
	}

	respData := &borrower.RoleEntity{}
	respData, err = copierx.Trans[dto.RoleEntityDto, borrower.RoleEntity](role, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}

	return &borrower.RoleCreateResponse{
		Data: respData,
	}, nil
}

package organizationuserlogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/data/dto"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
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

func (l *CreateLogic) Create(in *admin.OrganizationUserCreateRequest) (resp *admin.OrganizationUserDetailResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	input, err := copierx.Trans[admin.OrganizationUserCreateRequestData, dto.OrganizationUserCreateAdminRequestDto](in.Data)
	if err != nil {
		return nil, err
	}

	input.UserType = claim.UserType
	input.OrganizationId = claim.OrganizationId
	user, err := l.baseLogic.OrganizationUserLogic.CreateAdminOrganizationUser(l.ctx, *input)
	if err != nil {
		return nil, err
	}

	data, err := copierx.Trans[dto.OrganizationUserEntityDto, admin.OrganizationUserEntity](user, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}
	return &admin.OrganizationUserDetailResponse{
		Data: data,
	}, nil
}

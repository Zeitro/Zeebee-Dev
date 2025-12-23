package organizationuserlogic

import (
	"context"

	"github.com/jinzhu/copier"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
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
		baseLogic: logic.NewBaseLogic(ctx, svcCtx),
	}
}

func (l *CreateLogic) Create(in *admin.OrganizationUserCreateRequest) (resp *admin.OrganizationUserCreateResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	input := dto.OrganizationUserCreateAdminRequestDto{}
	err = copier.Copy(&input, in.Data)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	input.UserType = claim.UserType
	input.OrganizationId = claim.OrganizationId
	user, err := l.baseLogic.OrganizationUserLogic.CreateAdminOrganizationUser(l.ctx, input)
	if err != nil {
		return nil, err
	}
	data := &admin.OrganizationUserEntity{}
	err = copier.CopyWithOption(data, user, copier.Option{
		DeepCopy:      true,
		CaseSensitive: true,
	})
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	return &admin.OrganizationUserCreateResponse{
		Data: data,
	}, nil
}

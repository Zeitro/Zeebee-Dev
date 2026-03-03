package organizationuserlogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/data/dto"

	"zeebee-backend/common/errorx"
)

type UpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.NewBaseLogic(ctx, svcCtx),
	}
}

func (l *UpdateLogic) Update(in *admin.OrganizationUserUpdateRequest) (resp *admin.OrganizationUserDetailResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	input, err := copierx.Trans[admin.OrganizationUserUpdateRequest, dto.OrganizationUserUpdateRequestDto](in)
	if err != nil {
		return nil, err
	}

	// input := dto.OrganizationUserUpdateRequestDto{}
	// err = copier.CopyWithOption(&input, &in.Data, copier.Option{CaseSensitive: true, DeepCopy: true})
	// if err != nil {
	// 	return nil, err
	// }
	input.UserType = claim.UserType
	input.OrganizationId = claim.OrganizationId

	result, err := l.baseLogic.OrganizationUserLogic.UpdateOrganizationUser(l.ctx, *input)
	if err != nil {
		return nil, err
	}
	data, err := copierx.Trans[dto.OrganizationUserEntityDto, admin.OrganizationUserEntity](result, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}
	// data := &admin.OrganizationUserEntity{}
	// err = copier.CopyWithOption(data, result, copier.Option{
	// 	CaseSensitive: true,
	// 	DeepCopy:      true,
	// 	Converters: []copier.TypeConverter{
	// 		copierx.TimeToRFC3339Str(),
	// 	},
	// })
	// if err != nil {
	// 	return nil, err
	// }

	return &admin.OrganizationUserDetailResponse{
		Data: data,
	}, nil
}

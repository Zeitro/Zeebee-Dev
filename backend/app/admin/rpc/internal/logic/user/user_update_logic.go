package userlogic

import (
	"context"

	"github.com/jinzhu/copier"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/constants"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/stringx"
	"zeebee-backend/common/utils"

	"zeebee-backend/common/errorx"
)

type UserUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewUserUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserUpdateLogic {
	return &UserUpdateLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.NewBaseLogic(ctx, svcCtx),
	}
}

func (l *UserUpdateLogic) UserUpdate(in *admin.UserUpdateRequest) (resp *admin.UserUpdateResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	userModel, err := l.svcCtx.UserRepo.GetOrganizationUser(l.ctx, in.Data.Id, in.Data.OrganizationId, constants.UserTypeAdmin)
	if err != nil {
		return nil, err
	}
	if len(in.Data.Password) > 0 {
		userModel.PasswordHash, err = utils.HashPassword(in.Data.Password)
		if err != nil {
			return nil, errorx.WithStack(err)
		}
	}

	userModel.FirstName, userModel.MidName, userModel.LastName = stringx.ParseFullName(in.Data.FullName)

	user, err := l.baseLogic.UserLogic.UpdateOrganizationUser(l.ctx, in.Data.Id, userModel, in.Data.OrganizationId, in.Data.RoleIds)
	if err != nil {
		return nil, err
	}
	userEntity := &admin.UserEntity{}
	err = copier.CopyWithOption(userEntity, user, copier.Option{
		DeepCopy:      true,
		CaseSensitive: true,
		Converters: []copier.TypeConverter{
			copierx.TimeToRFC3339Str(),
		},
	})
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	return &admin.UserUpdateResponse{
		Data: userEntity,
	}, nil
}

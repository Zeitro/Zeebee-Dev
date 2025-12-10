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
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/stringx"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/model"
)

type UserCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewUserCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserCreateLogic {
	return &UserCreateLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.NewBaseLogic(ctx, svcCtx),
	}
}

func (l *UserCreateLogic) UserCreate(in *admin.UserCreateRequest) (resp *admin.UserCreateResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	// Password is required for user creation
	if in.Data.Password == "" {
		return nil, errorx.New("Password Is Required")
	}

	userModel := &model.User{}
	err = copier.CopyWithOption(userModel, in.Data, copier.Option{
		DeepCopy:      true,
		CaseSensitive: true,
		Converters: []copier.TypeConverter{
			copierx.TimeToRFC3339Str(),
		},
	})
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	userModel.Status = constants.UserStatusActive
	userModel.Type = constants.UserTypeAdmin
	userModel.PasswordHash, err = utils.HashPassword(in.Data.Password)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	userModel.FirstName, userModel.MidName, userModel.LastName = stringx.ParseFullName(in.Data.FullName)
	userModel.OrganizationUsers = []*model.OrganizationUser{
		{
			OrganizationId: claim.OrganizationId,
			Status:         constants.OrganizationUserStatusActive,
			FirstName:      userModel.FirstName,
			MidName:        userModel.MidName,
			LastName:       userModel.LastName,
			Phone:          in.Data.Phone,
			FullName:       userModel.FullName,
		},
	}
	// todo remove to OrganizationUser repo
	user, err := l.baseLogic.UserLogic.AddOrganizationUser(l.ctx, userModel, claim.OrganizationId, in.Data.RoleIds)
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
	return &admin.UserCreateResponse{
		Data: userEntity,
	}, nil
}

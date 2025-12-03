package userlogic

import (
	"context"

	"github.com/jinzhu/copier"

	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/constants"
	"zeebee-backend/common/copierx"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type UserDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserDetailLogic {
	return &UserDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserDetailLogic) UserDetail(in *admin.UserDetailRequest) (resp *admin.UserDetailResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	user, err := l.svcCtx.UserRepo.GetOrganizationUser(l.ctx, in.Data.Id, in.Data.OrganizationId, constants.UserTypeAdmin)
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

	relation, err := l.svcCtx.OrganizationAndUserRepo.GetRelation(l.ctx, in.Data.OrganizationId, in.Data.Id)
	if err != nil {
		return nil, err
	}
	userEntity.PivotOrganizationStatus = string(relation.Status)

	return &admin.UserDetailResponse{
		Data: userEntity,
	}, nil
}

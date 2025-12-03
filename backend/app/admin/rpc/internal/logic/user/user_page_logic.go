package userlogic

import (
	"context"

	"github.com/jinzhu/copier"

	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/constants"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/timex"
	"zeebee-backend/domains/zeebee/repo"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type UserPageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserPageLogic {
	return &UserPageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// admin user
func (l *UserPageLogic) UserPage(in *admin.UserPageRequest) (resp *admin.UserPageResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	resultData := &admin.UserPageResponseData{}
	pageReq := &repo.UserPageRequest{
		Filter: &repo.UserPageRequestFilter{},
		Sorter: &repo.UserPageRequestSorter{},
	}
	err = copier.CopyWithOption(pageReq, in, copier.Option{
		CaseSensitive: true,
		DeepCopy:      true,
	})
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	pageReq.Filter.UserType = constants.UserTypeAdmin
	page, err := l.svcCtx.UserRepo.Page(l.ctx, pageReq)
	if err != nil {
		return nil, err
	}
	err = copier.CopyWithOption(resultData, page, copier.Option{
		CaseSensitive: true,
		DeepCopy:      true,
		Converters: []copier.TypeConverter{
			copierx.StrToSqlTimePtr(timex.LayoutUsStyle),
			copierx.TimeToRFC3339Str(),
		},
	})
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	return &admin.UserPageResponse{
		Data: resultData,
	}, nil
}

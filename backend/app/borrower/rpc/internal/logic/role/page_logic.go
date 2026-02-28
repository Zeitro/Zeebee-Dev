package rolelogic

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
)

type PageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PageLogic {
	return &PageLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.NewBaseLogic(ctx, svcCtx),
	}
}

func (l *PageLogic) Page(in *borrower.RolePageRequest) (resp *borrower.RolePageResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	pageReq := dto.RolePageRequestDto{}
	_ = copier.CopyWithOption(&pageReq, in.Data, copier.Option{
		CaseSensitive: true,
		DeepCopy:      true,
	})

	pageReq.Filter.RoleType = constants.UserTypeBorrower
	pageReq.Filter.OrganizationId = claim.OrganizationId
	// pageReq.Filter.UserId = claim.OrganizationUserId

	result, err := l.baseLogic.RoleLogic.Page(l.ctx, pageReq)
	if err != nil {
		return nil, err
	}

	data, err := copierx.Trans[gormx.PageResponse[dto.RoleEntityDto], borrower.RolePageResponseData](&result, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}

	// data := &borrower.RolePageResponseData{}
	// err = copier.CopyWithOption(data, &result, copier.Option{
	// 	DeepCopy: true,
	// 	Converters: []copier.TypeConverter{
	// 		copierx.TimeToRFC3339Str(),
	// 	},
	// })
	// if err != nil {
	// 	return nil, errorx.GetFailedTransRespErr(err)
	// }

	return &borrower.RolePageResponse{
		Data: data,
	}, nil
}

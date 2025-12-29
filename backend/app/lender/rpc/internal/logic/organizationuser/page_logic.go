package organizationuserlogic

import (
	"context"

	"github.com/jinzhu/copier"

	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/bo"
	"zeebee-backend/domains/zeebee/data/dto"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
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

func (l *PageLogic) Page(in *lender.OrganizationUserPageRequest) (resp *lender.OrganizationUserPageResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	input := dto.OrganizationUserPageDto{
		OrganizationId: claim.OrganizationId,
		Page:           in.Data.Page,
		PageSize:       in.Data.PageSize,
		Keyword:        in.Data.Keyword,
		Status:         in.Data.Status,
		UserType:       constants.UserTypeLender,
	}
	var result *gormx.PageResponse[*bo.OrganizationUserBo]
	if claim.UserType == constants.UserTypeSuperAdmin {
		result, err = l.baseLogic.OrganizationUserLogic.Page(l.ctx, input, 0)
	} else {
		result, err = l.baseLogic.OrganizationUserLogic.Page(l.ctx, input, claim.OrganizationId)
	}
	if err != nil {
		return nil, err
	}
	var list []*lender.OrganizationUserEntity
	for _, v := range result.List {
		item, err := copier.CopyWithOption(&lender.OrganizationUserEntity{}, v, copier.Option{
			CaseSensitive: true,
			Converters: []copier.TypeConverter{
				copierx.TimeToRFC3339Str(),
			},
		})
		if err != nil {
			return nil, errorx.GetFailedTransRespErr(err)
		}
		list = append(list, item.(*lender.OrganizationUserEntity))
	}

	return &lender.OrganizationUserPageResponse{
		Data: &lender.PageResponseOrganizationUserEntity{
			Total: result.Total,
			List:  list,
		},
	}, nil
}

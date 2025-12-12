package organizationuserlogic

import (
	"context"

	"github.com/jinzhu/copier"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
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

func (l *PageLogic) Page(in *admin.OrganizationUserPageRequest) (resp *admin.OrganizationUserPageResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	pageReq := dto.OrganizationUserPageRequestDto{}
	err = copier.Copy(&pageReq, &in.Data)
	if err != nil {
		return nil, err
	}
	pageReq.Filter.UserType = claim.UserType
	pageReq.Filter.OrganizationId = claim.OrganizationId
	pageReq.Filter.UserId = claim.UserId

	results, err := l.baseLogic.OrganizationUserLogic.Page(l.ctx, pageReq)
	if err != nil {
		return nil, err
	}
	data := &admin.OrganizationUserPageResponseData{}
	err = copier.CopyWithOption(data, &results, copier.Option{
		DeepCopy:      true,
		CaseSensitive: true,
		Converters: []copier.TypeConverter{
			copierx.TimeToRFC3339Str(),
		},
	})
	if err != nil {
		return nil, err
	}

	return &admin.OrganizationUserPageResponse{
		Data: data,
	}, nil
}

package organizationuserlogic

import (
	"context"

	"github.com/jinzhu/copier"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/proto"

	"zeebee-backend/common/errorx"
)

type AllLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	*logic.BaseLogic
}

func NewAllLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AllLogic {
	return &AllLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		BaseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *AllLogic) All(in *proto.EmptyRequest) (resp *borrower.OrganizationUserAllResponse, err error) {

	claimUser, err := l.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	results, err := l.svcCtx.OrganizationUserRepo.AllOrganizationUsers(l.ctx, claimUser.OrganizationId)
	if err != nil {
		return nil, err
	}
	data := []*proto.OrganizationUserEntity{}
	err = copier.CopyWithOption(&data, &results, copier.Option{
		DeepCopy:      true,
		CaseSensitive: true,
		Converters: []copier.TypeConverter{
			copierx.TimeToRFC3339Str(),
		},
	})
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	return &borrower.OrganizationUserAllResponse{
		Data: data,
	}, nil
}

package organization_user

import (
	"context"

	"github.com/jinzhu/copier"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrganizationUserAllLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// all
func NewOrganizationUserAllLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationUserAllLogic {
	return &OrganizationUserAllLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationUserAllLogic) OrganizationUserAll(req *types.EmptyRequest) (resp *types.OrganizationUserAllResponse, err error) {
	rpcReq := &borrower.EmptyRequest{}
	rpcResp, err := l.svcCtx.OrganizationUserRpcClient.All(l.ctx, rpcReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	resp = &types.OrganizationUserAllResponse{
		List: []*types.OrganizationUserEntity{},
	}
	err = copier.CopyWithOption(&resp.List, rpcResp.Data, copier.Option{
		DeepCopy:      true,
		CaseSensitive: true,
		Converters: []copier.TypeConverter{
			copierx.TimeToRFC3339Str(),
		},
	})
	if err != nil {
		err = errorx.GetFailedTransRespErr(err)
		return nil, err
	}
	return resp, nil
}

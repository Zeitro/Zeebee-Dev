package organization_user

import (
	"context"

	"github.com/jinzhu/copier"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrganizationUserUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// update
func NewOrganizationUserUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationUserUpdateLogic {
	return &OrganizationUserUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationUserUpdateLogic) OrganizationUserUpdate(req *types.OrganizationUserUpdateRequest) (resp *types.OrganizationUserUpdateResponse, err error) {
	rpcReq := &admin.OrganizationUserUpdateRequest{
		Data: &admin.OrganizationUserUpdateRequestData{},
	}
	err = copier.CopyWithOption(rpcReq.Data, req, copier.Option{
		DeepCopy:      true,
		CaseSensitive: true,
	})
	if err != nil {
		err = errorx.GetFailedTransReqErr(err)
		return nil, err
	}

	rpcResp, err := l.svcCtx.OrganizationUserRpcClient.Update(l.ctx, rpcReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	resp = &types.OrganizationUserUpdateResponse{}
	err = copier.CopyWithOption(resp, rpcResp.Data, copier.Option{
		DeepCopy:      true,
		CaseSensitive: true,
		Converters: []copier.TypeConverter{
			copierx.TimeToRFC3339Str(),
		},
	})
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	return resp, nil
}

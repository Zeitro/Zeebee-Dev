// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package organization_email

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrganizationEmailGenerateAuthCodeUrlLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// generate auth code url
func NewOrganizationEmailGenerateAuthCodeUrlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationEmailGenerateAuthCodeUrlLogic {
	return &OrganizationEmailGenerateAuthCodeUrlLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationEmailGenerateAuthCodeUrlLogic) OrganizationEmailGenerateAuthCodeUrl(req *types.OrganizationEmailGenerateAuthCodeUrlRequest) (resp *types.OrganizationEmailGenerateAuthCodeUrlResponse, err error) {
	rpcReq := &borrower.OrganizationEmailGenerateAuthCodeUrlRequest{}

	rpcReq.Data, err = copierx.Trans[types.OrganizationEmailGenerateAuthCodeUrlRequest, borrower.OrganizationEmailGenerateAuthCodeUrlRequestData](req)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	rpcResp, err := l.svcCtx.OrganizationEmailRpcClient.GenerateAuthCodeUrl(l.ctx, rpcReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	resp, err = copierx.Trans[borrower.OrganizationEmailGenerateAuthCodeUrlResponseData, types.OrganizationEmailGenerateAuthCodeUrlResponse](rpcResp.Data)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	return resp, nil
}

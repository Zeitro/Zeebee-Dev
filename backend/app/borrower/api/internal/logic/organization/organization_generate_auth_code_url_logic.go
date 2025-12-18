// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package organization

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrganizationGenerateAuthCodeUrlLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// generate auth code url
func NewOrganizationGenerateAuthCodeUrlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationGenerateAuthCodeUrlLogic {
	return &OrganizationGenerateAuthCodeUrlLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationGenerateAuthCodeUrlLogic) OrganizationGenerateAuthCodeUrl(req *types.OrganizationGenerateAuthCodeUrlRequest) (resp *types.OrganizationGenerateAuthCodeUrlResponse, err error) {
	rpcReq := &borrower.OrganizationGenerateAuthCodeUrlRequest{}

	rpcReq.Data, err = copierx.Trans[types.OrganizationGenerateAuthCodeUrlRequest, borrower.OrganizationGenerateAuthCodeUrlRequestData](req)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	rpcResp, err := l.svcCtx.OrganizationRpcClient.GenerateAuthCodeUrl(l.ctx, rpcReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	resp, err = copierx.Trans[borrower.OrganizationGenerateAuthCodeUrlResponseData, types.OrganizationGenerateAuthCodeUrlResponse](rpcResp.Data)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	return resp, nil
}

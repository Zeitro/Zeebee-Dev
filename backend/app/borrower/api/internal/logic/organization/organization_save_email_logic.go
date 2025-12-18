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

type OrganizationSaveEmailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// save email
func NewOrganizationSaveEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationSaveEmailLogic {
	return &OrganizationSaveEmailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationSaveEmailLogic) OrganizationSaveEmail(req *types.OrganizationSaveEmailRequest) (resp *types.EmptyResponse, err error) {
	rpcReq := &borrower.OrganizationSaveEmailRequest{}

	rpcReq.Data, err = copierx.Trans[types.OrganizationSaveEmailRequest, borrower.OrganizationSaveEmailRequestData](req)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	_, err = l.svcCtx.OrganizationRpcClient.SaveEmail(l.ctx, rpcReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	return &types.EmptyResponse{}, nil
}

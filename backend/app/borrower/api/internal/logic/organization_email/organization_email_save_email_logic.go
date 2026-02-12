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

type OrganizationEmailSaveEmailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// save email
func NewOrganizationEmailSaveEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationEmailSaveEmailLogic {
	return &OrganizationEmailSaveEmailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationEmailSaveEmailLogic) OrganizationEmailSaveEmail(req *types.OrganizationEmailSaveEmailRequest) (resp *types.EmptyResponse, err error) {
	rpcReq := &borrower.OrganizationEmailSaveEmailRequest{}

	rpcReq.Data, err = copierx.Trans[types.OrganizationEmailSaveEmailRequest, borrower.OrganizationEmailSaveEmailRequestData](req)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	_, err = l.svcCtx.OrganizationEmailRpcClient.SaveEmail(l.ctx, rpcReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	return &types.EmptyResponse{}, nil
}

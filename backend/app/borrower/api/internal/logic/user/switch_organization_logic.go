package user

import (
	"context"

	"github.com/jinzhu/copier"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type SwitchOrganizationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// switch organization
func NewSwitchOrganizationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SwitchOrganizationLogic {
	return &SwitchOrganizationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SwitchOrganizationLogic) SwitchOrganization(req *types.UserSwitchOrganizationRequest) (resp *types.UserLoginResponse, err error) {
	resp = &types.UserLoginResponse{}
	rpcReq := &borrower.SwitchOrganizationRequest{
		Data: &borrower.SwitchOrganizationRequestData{},
	}
	claim, err := utils.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	userId := claim.UserId
	if userId == 0 {
		return nil, errorx.WithStack(errorx.New("User Id Not Found"))
	}
	rpcReq.Data.UserId = userId
	rpcReq.Data.OrganizationId = req.OrganizationId

	rpcResp, err := l.svcCtx.UserRpcClient.SwitchOrganization(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	err = copier.Copy(resp, &rpcResp.Data)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	return resp, err
}

// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"

	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/utils"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type SwitchLenderCompanyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// switch lender company
func NewSwitchLenderCompanyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SwitchLenderCompanyLogic {
	return &SwitchLenderCompanyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SwitchLenderCompanyLogic) SwitchLenderCompany(req *types.UserSwitchLenderCompanyRequest) (resp *types.UserRefreshTokenResponse, err error) {
	_, employeeId, _ := utils.GetLenderCompanyAndEmployeeId(l.ctx)
	if employeeId == 0 {
		return nil, errors.WithStack(errors.New("employee id not found"))
	}
	rpcResp, err := l.svcCtx.UserRpcClient.SwitchLenderCompany(l.ctx, &lender.SwitchLenderCompanyRequest{
		Data: &lender.SwitchLenderCompanyRequestData{
			EmployeeId:      employeeId,
			LenderCompanyId: req.LenderCompanyId,
		},
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &types.UserRefreshTokenResponse{
		AccessToken:        rpcResp.Data.AccessToken,
		AccessTokenExpire:  rpcResp.Data.AccessTokenExpire,
		RefreshToken:       rpcResp.Data.RefreshToken,
		RefreshTokenExpire: rpcResp.Data.RefreshTokenExpire,
	}, nil
}

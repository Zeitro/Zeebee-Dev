package user

import (
	"context"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// user menu
func NewUserMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserMenuLogic {
	return &UserMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserMenuLogic) UserMenu(req *types.EmptyRequest) (resp *types.ResourceTreeResponse, err error) {
	resp = &types.ResourceTreeResponse{
		List: []*types.ResourceTreeEntity{},
	}

	rpcResp, err := l.svcCtx.ResourceRpcClient.UserMenu(l.ctx, &lender.EmptyRequest{})
	if err != nil {
		return nil, err
	}

	// trans response
	result, err := copierx.Trans[[]*lender.ResourceTreeEntity, []*types.ResourceTreeEntity](&rpcResp.Data)
	if err != nil {
		return nil, err
	}

	return &types.ResourceTreeResponse{
		List: *result,
	}, nil
}

package common

import (
	"context"
	"encoding/json"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/swagger"

	"github.com/zeromicro/go-zero/core/logx"
)

type SyncApiLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// sync api
func NewSyncApiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncApiLogic {
	return &SyncApiLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SyncApiLogic) SyncApi(req *types.EmptyRequest) (resp *types.EmptyResponse, err error) {
	apiGroup, err := swagger.ParseApi(l.svcCtx.BorrowerApiContent)
	if err != nil {
		return nil, errorx.Wrap(err, "failed to parse lender.api")
	}

	apiJson, err := json.Marshal(apiGroup)
	if err != nil {
		return nil, errorx.Wrap(err, "failed to marshal api group")
	}

	_, err = l.svcCtx.ApiRpcClient.SyncApi(l.ctx, &lender.SyncApiRequest{
		Data: &lender.SyncApiRequestData{
			ApiJson: string(apiJson),
		},
	})
	if err != nil {
		return nil, err
	}

	return &types.EmptyResponse{}, nil
}

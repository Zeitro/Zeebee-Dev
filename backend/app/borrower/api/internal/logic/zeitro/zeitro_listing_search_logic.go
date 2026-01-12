package zeitro

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"

	"github.com/zeromicro/go-zero/core/logx"
)

type ZeitroListingSearchLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Search property listings from Zeitro API
func NewZeitroListingSearchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ZeitroListingSearchLogic {
	return &ZeitroListingSearchLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ZeitroListingSearchLogic) ZeitroListingSearch(req *types.ZeitroListingSearchRequest) (resp *types.ZeitroListingSearchResponse, err error) {
	// Convert API request to RPC request using copierx.Trans
	rpcReq, err := copierx.Trans[types.ZeitroListingSearchRequest, borrower.ZeitroListingSearchRequest](req)
	if err != nil {
		return nil, err
	}

	// Call RPC to search listings
	rpcResp, err := l.svcCtx.ZeitroRpcClient.ZeitroListingSearch(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// Convert RPC response to API response using copierx.Trans
	entities, err := copierx.Trans[[]*borrower.ZeitroPropertyEntity, []types.ZeitroPropertyEntity](&rpcResp.Data)
	if err != nil {
		return nil, err
	}

	return &types.ZeitroListingSearchResponse{
		List: *entities,
	}, nil
}

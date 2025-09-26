package loan

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/utils"
)

type DetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// detail
func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DetailLogic) Detail(req *types.LoanDetailRequest) (resp *types.LoanDetailResponse, err error) {
	resp = &types.LoanDetailResponse{}
	// get lender id
	lenderCompanyId, err := utils.GetLenderCompanyIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	rpcResp, err := l.svcCtx.LoanRpcClient.LoanDetail(l.ctx, &lender.LoanDetailRequest{
		Data: &lender.LoanDetailRequestData{
			Id:              req.Id,
			LenderCompanyId: lenderCompanyId,
		},
	})
	if err != nil {
		return nil, err
	}
	err = copier.CopyWithOption(resp, &rpcResp.Data, copier.Option{
		CaseSensitive: true,
		DeepCopy:      true,
	})
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	return resp, nil
}

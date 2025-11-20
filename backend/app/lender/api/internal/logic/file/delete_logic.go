package file

import (
	"context"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// delete
func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteLogic) Delete(req *types.PathIdRequest) (resp *types.EmptyResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	// Get lender company and employee id from context
	lenderCompanyId, lenderEmployeeId, err := utils.GetLenderCompanyAndEmployeeId(l.ctx)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	// Call RPC to delete file
	_, err = l.svcCtx.FileRpcClient.DeleteFile(l.ctx, &lender.FileDeleteRequest{
		Data: &lender.FileDeleteRequestData{
			Id:               req.Id,
			LenderCompanyId:  lenderCompanyId,
			LenderEmployeeId: lenderEmployeeId,
		},
	})
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	return &types.EmptyResponse{}, nil
}

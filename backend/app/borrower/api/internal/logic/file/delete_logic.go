package file

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/errorx"
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
	// Call RPC to delete file
	_, err = l.svcCtx.FileRpcClient.DeleteFile(l.ctx, &borrower.FileDeleteRequest{
		Data: &borrower.FileDeleteRequestData{
			Id: req.Id,
		},
	})
	if err != nil {
		return nil, errorx.GetFailedTransReqErr(err)
	}

	return &types.EmptyResponse{}, nil
}

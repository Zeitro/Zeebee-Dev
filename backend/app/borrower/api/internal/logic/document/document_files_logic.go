package document

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/proto"
)

type DocumentFilesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDocumentFilesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DocumentFilesLogic {
	return &DocumentFilesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DocumentFilesLogic) DocumentFiles(req *types.DocumentFilesRequest) (resp *types.DocumentFilesResponse, err error) {
	rpcResp, err := l.svcCtx.DocumentRpcClient.Files(l.ctx, &borrower.DocumentFilesRequest{
		Data: &borrower.DocumentFilesRequestData{
			DocumentId: req.Id,
		},
	})
	if err != nil {
		return nil, err
	}

	files, err := copierx.Trans[[]*proto.FileEntity, []types.DocumentFile](&rpcResp.Data)
	if err != nil {
		return nil, err
	}

	return &types.DocumentFilesResponse{
		Files: *files,
	}, nil
}

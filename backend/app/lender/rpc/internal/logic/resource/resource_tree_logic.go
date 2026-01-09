package resourcelogic

import (
	"context"

	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type ResourceTreeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewResourceTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResourceTreeLogic {
	return &ResourceTreeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ResourceTreeLogic) ResourceTree(in *lender.ResourceTreeRequest) (resp *lender.ResourceTreeResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	// todo: add your logic here and delete this line

	return &lender.ResourceTreeResponse{}, nil
}

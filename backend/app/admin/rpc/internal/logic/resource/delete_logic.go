package resourcelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
)

type DeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.NewBaseLogic(ctx, svcCtx),
	}
}

func (l *DeleteLogic) Delete(in *admin.ResourceDeleteRequest) (resp *admin.EmptyResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	err = l.svcCtx.ResourceRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		var innerErr error
		innerErr = l.svcCtx.ResourceRepo.DelByIdWithTx(l.ctx, tx, in.Data.Id, constants.UserType(in.Data.Type))
		if innerErr != nil {
			return innerErr
		}
		// todo update user permissions

		return innerErr
	})
	if err != nil {
		return nil, err
	}

	return &admin.EmptyResponse{}, nil
}

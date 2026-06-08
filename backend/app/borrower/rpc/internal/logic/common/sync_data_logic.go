package commonlogic

import (
	"context"

	"gorm.io/gorm"

	"zeebee-backend/app/borrower/rpc/internal/biz_helper"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type SyncDataLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSyncDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncDataLogic {
	return &SyncDataLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SyncDataLogic) SyncData(in *borrower.EmptyRequest) (resp *borrower.EmptyResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	if err = l.syncBorrowerPropertyStatuses(); err != nil {
		return nil, err
	}

	return &borrower.EmptyResponse{}, nil
}

func (l *SyncDataLogic) syncBorrowerPropertyStatuses() error {
	propertyHelper := biz_helper.NewBorrowerProperty(l.svcCtx.AllRepos)
	err := l.svcCtx.BorrowerPropertyRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		var innerErr error
		var propertyIds []int64
		innerErr = l.svcCtx.BorrowerPropertyRepo.PluckByColumnsWithTx(l.ctx, tx, "id", &propertyIds, nil)
		if innerErr != nil {
			innerErr = errorx.Wrap(innerErr, "failed to pluck property ids")
			return innerErr
		}

		for _, propertyId := range propertyIds {
			innerErr = propertyHelper.CalculatePropertyStatuses(l.ctx, tx, propertyId)
			if innerErr != nil {
				innerErr = errorx.Wrapf(innerErr, "failed to calculate property[%d] statuses", propertyId)
				return innerErr
			}
		}
		return innerErr
	})

	return err
}

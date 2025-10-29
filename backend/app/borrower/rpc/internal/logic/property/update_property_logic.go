package propertylogic

import (
	"context"

	"github.com/jinzhu/copier"

	"zeebee-backend/app/borrower/rpc/internal/repo/model"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type UpdatePropertyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdatePropertyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePropertyLogic {
	return &UpdatePropertyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdatePropertyLogic) UpdateProperty(in *borrower.PropertyUpdateRequest) (resp *borrower.EmptyResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	var result *model.PropertyBase
	result, err = l.svcCtx.PropertyBaseRepo.FindOneByColumns(l.ctx, map[string]interface{}{
		"id":          in.Data.Id,
		"borrower_id": in.Data.BorrowerId,
	})
	if err != nil {
		return nil, err
	}

	err = copier.CopyWithOption(result, in.Data, copier.Option{
		CaseSensitive: true,
		DeepCopy:      true,
		Converters:    []copier.TypeConverter{},
	})
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	_, err = l.svcCtx.PropertyBaseRepo.UpdateOne(l.ctx, in.Data.Id, result)

	return &borrower.EmptyResponse{}, err
}

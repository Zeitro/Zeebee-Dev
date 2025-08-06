package commonlogic

import (
	"context"

	"github.com/jinzhu/copier"

	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type GetUsCitiesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUsCitiesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUsCitiesLogic {
	return &GetUsCitiesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUsCitiesLogic) GetUsCities(in *borrower.UsCitiesRequest) (resp *borrower.UsCitiesResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	usCityModels, err := l.svcCtx.UsCityRepo.FindAll(l.ctx, map[string]interface{}{"state_id": in.Data.StateId})
	if err != nil {
		return nil, err
	}
	cities := make([]*borrower.UsCityEntry, 0, len(usCityModels))
	err = copier.Copy(&cities, &usCityModels)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	return &borrower.UsCitiesResponse{
		Data: cities,
	}, nil
}

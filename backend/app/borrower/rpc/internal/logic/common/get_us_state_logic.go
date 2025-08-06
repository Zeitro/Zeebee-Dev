package commonlogic

import (
	"context"

	"github.com/jinzhu/copier"

	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type GetUsStateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUsStateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUsStateLogic {
	return &GetUsStateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUsStateLogic) GetUsState(in *borrower.EmptyRequest) (resp *borrower.UsStateResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	usStateModels, err := l.svcCtx.UsStateRepo.FindAll(l.ctx, nil)
	if err != nil {
		return nil, err
	}

	states := make([]*borrower.UsStateEntry, 0, len(usStateModels))

	err = copier.Copy(&states, &usStateModels)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	return &borrower.UsStateResponse{Data: states}, nil
}

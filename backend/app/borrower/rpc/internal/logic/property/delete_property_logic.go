package propertylogic

import (
	"context"

	"github.com/pkg/errors"

	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type DeletePropertyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeletePropertyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletePropertyLogic {
	return &DeletePropertyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeletePropertyLogic) DeleteProperty(in *borrower.PropertyDeleteRequest) (resp *borrower.EmptyResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	result, err := l.svcCtx.PropertyBaseRepo.FindOneById(l.ctx, in.Data.Id)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	if result.BorrowerId != in.Data.BorrowerId {
		return nil, errorx.GetFailedTransRespErr(errors.New("You Have No Access To This Property"))
	}

	if result.LoanCount > 0 {
		return nil, errorx.GetFailedTransRespErr(errors.New("You Cannot Delete A Property That Has Loans"))
	}
	_, err = l.svcCtx.PropertyBaseRepo.DeleteByID(l.ctx, in.Data.Id)

	return &borrower.EmptyResponse{}, err
}

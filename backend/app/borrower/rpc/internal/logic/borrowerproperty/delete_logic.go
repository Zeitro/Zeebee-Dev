package borrowerpropertylogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type DeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteLogic) Delete(in *borrower.PropertyDeleteRequest) (resp *borrower.EmptyResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	result, err := l.svcCtx.BorrowerPropertyRepo.FindOneById(l.ctx, in.Data.Id)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	// todo check organization
	if result.UserId != in.Data.UserId {
		return nil, errorx.GetFailedTransRespErr(errorx.New("You Have No Access To This Property"))
	}

	if result.LoanCount > 0 {
		return nil, errorx.GetFailedTransRespErr(errorx.New("You Cannot Delete A Property That Has Loans"))
	}
	_, err = l.svcCtx.BorrowerPropertyRepo.DeleteByID(l.ctx, in.Data.Id)

	return &borrower.EmptyResponse{}, err
}

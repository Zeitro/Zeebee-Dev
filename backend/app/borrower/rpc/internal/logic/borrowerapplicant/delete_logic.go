package borrowerapplicantlogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
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

func (l *DeleteLogic) Delete(in *borrower.BorrowerApplicantDeleteRequest) (resp *borrower.EmptyResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	result, err := l.svcCtx.BorrowerApplicantRepo.FindOneById(l.ctx, in.Data.Id)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	// todo check organization
	if result.UserId != in.Data.UserId {
		return nil, errorx.GetFailedTransRespErr(errorx.New("You Have No Access To This Applicant"))
	}

	_, err = l.svcCtx.BorrowerApplicantRepo.DeleteByID(l.ctx, in.Data.Id)

	return &borrower.EmptyResponse{}, err
}

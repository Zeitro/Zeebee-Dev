package borrowerentitylogic

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

func (l *DeleteLogic) Delete(in *borrower.BorrowerEntityDeleteRequest) (resp *borrower.EmptyResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	result, err := l.svcCtx.BorrowerEntityRepo.GetOneByIdForOrganization(l.ctx, in.Data.Id, in.Data.OrganizationId)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	if result.OrganizationId != in.Data.OrganizationId {
		return nil, errorx.GetFailedTransRespErr(errorx.New("You Have No Access To This Entity"))
	}

	_, err = l.svcCtx.BorrowerEntityRepo.DeleteByID(l.ctx, in.Data.Id)

	return &borrower.EmptyResponse{}, err
}

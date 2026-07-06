package borrowerentitylogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/domains/zeebee/proto"

	"zeebee-backend/common/errorx"
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
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *DeleteLogic) Delete(in *proto.PathIdRequest) (resp *proto.EmptyResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
				logx.Field("service_name", l.svcCtx.Config.Name),
			).Error(err)
		}
	}()

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	result, err := l.svcCtx.BorrowerEntityRepo.GetOneByIdForOrganization(l.ctx, in.Data.Id, claim.OrganizationId)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	if result.OrganizationId != claim.OrganizationId {
		return nil, errorx.GetFailedTransRespErr(errorx.New("You have no access to this entity"))
	}

	_, err = l.svcCtx.BorrowerEntityRepo.DeleteByID(l.ctx, in.Data.Id)

	return &proto.EmptyResponse{}, err
}

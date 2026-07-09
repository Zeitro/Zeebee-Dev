package borrowerentitylogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/proto"

	"zeebee-backend/common/errorx"
)

type UpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *UpdateLogic) Update(in *borrower.BorrowerEntityUpdateRequest) (resp *proto.EmptyResponse, err error) {
	defer func() {
		if err != nil {
			l.svcCtx.Logger.WithContext(l.ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
			).Error(err)
		}
	}()

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	var result *model.BorrowerEntity
	result, err = l.svcCtx.BorrowerEntityRepo.GetOneByIdForOrganization(l.ctx, in.Data.Id, claim.OrganizationId)
	if err != nil {
		return nil, err
	}
	result, err = copierx.Trans[borrower.BorrowerEntityEditEntity, model.BorrowerEntity](in.Data)
	if err != nil {
		return nil, err
	}

	_, err = l.svcCtx.BorrowerEntityRepo.UpdateOneById(l.ctx, in.Data.Id, result)

	return &proto.EmptyResponse{}, err
}

package borrowerapplicantlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/proto"
)

type DeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *DeleteLogic) Delete(in *proto.PathIdRequest) (resp *proto.EmptyResponse, err error) {

	claim, err := l.baseLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	result, err := l.svcCtx.BorrowerApplicantRepo.FindOneById(l.ctx, in.Data.Id)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	if result.OrganizationId != claim.OrganizationId {
		return nil, errorx.New("you have no access to this applicant")
	}

	_, err = l.svcCtx.BorrowerApplicantRepo.DeleteByID(l.ctx, in.Data.Id)

	return &proto.EmptyResponse{}, err
}

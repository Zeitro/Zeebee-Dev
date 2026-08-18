package borrowerpropertylogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/domains/zeebee/proto"

	"zeebee-backend/common/errorx"
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

func (l *DeleteLogic) Delete(in *borrower.PropertyDeleteRequest) (resp *proto.EmptyResponse, err error) {

	claim, err := l.baseLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	result, err := l.svcCtx.BorrowerPropertyRepo.GetOneByIdForOrganization(l.ctx, in.Data.Id, claim.OrganizationId)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	// todo check organization
	if result.OrganizationId != claim.OrganizationId {
		return nil, errorx.GetFailedTransRespErr(errorx.New("You have no access to this property"))
	}

	if result.LoanCount > 0 {
		return nil, errorx.GetFailedTransRespErr(errorx.New("You cannot delete a property that has loans"))
	}
	if err = l.baseLogic.ProjectManagementLogic.RequirePropertyDeletable(l.ctx, result.Id); err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	_, err = l.svcCtx.BorrowerPropertyRepo.DeleteByID(l.ctx, in.Data.Id)

	return &proto.EmptyResponse{}, err
}

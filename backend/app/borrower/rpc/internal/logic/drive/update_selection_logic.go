package drivelogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/proto"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type UpdateSelectionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewUpdateSelectionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSelectionLogic {
	return &UpdateSelectionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *UpdateSelectionLogic) UpdateSelection(in *borrower.DriveUpdateSelectionRequest) (resp *proto.EmptyResponse, err error) {
	defer func() {
		if err != nil {
			l.svcCtx.Logger.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	claim, err := l.baseLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	if in.Data == nil {
		return nil, errorx.New("Request data is empty")
	}

	input := dto.DriveUpdateSelectionRequestDto{
		OrganizationId:   claim.OrganizationId,
		OrganizationType: constants.UserTypeBorrower,
		DriveId:          in.Data.DriveId,
		DriveName:        in.Data.DriveName,
	}

	err = l.baseLogic.DriveLogic.UpdateSelection(l.ctx, input)
	if err != nil {
		return nil, err
	}

	return &proto.EmptyResponse{}, nil
}

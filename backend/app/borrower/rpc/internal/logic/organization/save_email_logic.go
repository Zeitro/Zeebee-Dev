package organizationlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type SaveEmailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewSaveEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveEmailLogic {
	return &SaveEmailLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.NewBaseLogic(ctx, svcCtx),
	}
}

func (l *SaveEmailLogic) SaveEmail(in *borrower.OrganizationSaveEmailRequest) (resp *borrower.EmptyResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	input, err := copierx.Trans[borrower.OrganizationSaveEmailRequestData, dto.OrganizationSaveEmailRequestDto](in.Data)
	if err != nil {
		return nil, err
	}
	input.OrganizationId = claim.OrganizationId
	input.OrganizationType = constants.UserTypeBorrower
	if input.EmailAuthStrategy == constants.EmailAuthStrategyUser {
		input.UserId = claim.UserId
	}

	err = l.baseLogic.OrganizationLogic.SaveEmail(l.ctx, *input)
	if err != nil {
		return nil, err
	}

	return &borrower.EmptyResponse{}, nil
}

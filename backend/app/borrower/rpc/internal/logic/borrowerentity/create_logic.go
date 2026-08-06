package borrowerentitylogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/timex"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/proto"
)

type CreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *CreateLogic) Create(in *borrower.BorrowerEntityCreateRequest) (resp *borrower.BorrowerEntityDetailResponse, err error) {

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	saveData := &model.BorrowerEntity{}
	responseData := &proto.BorrowerEntityEntity{}
	saveData, err = copierx.Trans[borrower.BorrowerEntityEditEntity, model.BorrowerEntity](in.Data)
	if err != nil {
		return nil, err
	}

	saveData.CreatorId = claim.OrganizationUserId
	saveData.OrganizationId = claim.OrganizationId
	saveData.CreatorName = claim.Username

	result, err := l.svcCtx.BorrowerEntityRepo.Create(l.ctx, saveData)
	if err != nil {
		return nil, err
	}
	responseData, err = copierx.Trans[model.BorrowerEntity, proto.BorrowerEntityEntity](result,
		copierx.TimeToRFC3339Str(),
		copierx.SqlTimePtrToStr(timex.LayoutUsStyle),
	)
	if err != nil {
		return nil, err
	}

	return &borrower.BorrowerEntityDetailResponse{
		Data: responseData,
	}, nil
}

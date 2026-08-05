package borrowerentitylogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/timex"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/proto"

	"zeebee-backend/common/errorx"
)

type AllLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewAllLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AllLogic {
	return &AllLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *AllLogic) All(in *proto.EmptyRequest) (resp *borrower.BorrowerEntityAllResponse, err error) {

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx, constants.UserTypeBorrower)
	if err != nil {
		return nil, err
	}

	results, err := l.svcCtx.BorrowerEntityRepo.AllForOrganization(l.ctx, claim.OrganizationId)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	respData := make([]*borrower.BorrowerEntityPageEntity, 0, len(results))
	respDataPtr, err := copierx.Trans[[]*model.BorrowerEntity, []*borrower.BorrowerEntityPageEntity](&results,
		copierx.TimeToRFC3339Str(),
		copierx.SqlTimePtrToStr(timex.LayoutUsStyle),
	)
	if err != nil {
		return nil, err
	}
	respData = *respDataPtr

	return &borrower.BorrowerEntityAllResponse{
		Data: respData,
	}, nil
}

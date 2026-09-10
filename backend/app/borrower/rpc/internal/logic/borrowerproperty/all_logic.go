package borrowerpropertylogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/timex"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/proto"
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

func (l *AllLogic) All(in *borrower.PropertyAllRequest) (resp *borrower.PropertyAllResponse, err error) {

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	results, err := l.svcCtx.BorrowerPropertyRepo.AllForOrganization(l.ctx, claim.OrganizationId)
	if err != nil {
		return nil, errorx.Wrap(err, "fail to get properties")
	}

	// respData := make([]*proto.BorrowerPropertyEntity, 0, len(results))

	data, err := copierx.Trans[[]*model.BorrowerProperty, []*proto.BorrowerPropertyEntity](&results, copierx.TimeToRFC3339Str(), copierx.SqlTimePtrToStr(timex.LayoutUsStyle))
	if err != nil {
		return nil, err
	}

	return &borrower.PropertyAllResponse{
		Data: *data,
	}, nil
}

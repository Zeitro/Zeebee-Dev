package borrowerapplicantlogic

import (
	"context"

	"github.com/jinzhu/copier"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/stringx"
	"zeebee-backend/common/timex"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/proto"
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

func (l *UpdateLogic) Update(in *borrower.BorrowerApplicantUpdateRequest) (resp *proto.EmptyResponse, err error) {

	var result *model.BorrowerApplicant

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx, constants.UserTypeBorrower)
	if err != nil {
		return nil, err
	}

	result, err = l.svcCtx.BorrowerApplicantRepo.FindOneByColumns(l.ctx, map[string]interface{}{
		"id":              in.Data.Id,
		"organization_id": claim.OrganizationId,
	})
	if err != nil {
		return nil, err
	}

	err = copier.CopyWithOption(result, in.Data, copier.Option{
		CaseSensitive: true,
		DeepCopy:      true,
		Converters: []copier.TypeConverter{
			copierx.StrToSqlTimePtr(timex.LayoutUsStyle),
		},
	})
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	result.FullName = stringx.FormatFullName(in.Data.FirstName, in.Data.MidName, in.Data.LastName)

	_, err = l.svcCtx.BorrowerApplicantRepo.UpdateOneById(l.ctx, in.Data.Id, result)

	return &proto.EmptyResponse{}, err
}

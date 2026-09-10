package borrowerpropertylogic

import (
	"context"

	"github.com/jinzhu/copier"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
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

func (l *CreateLogic) Create(in *borrower.PropertyCreateRequest) (resp *borrower.PropertyCreateResponse, err error) {

	claim, err := l.baseLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	saveData := &model.BorrowerProperty{}
	responseData := &proto.BorrowerPropertyEntity{}
	err = copier.CopyWithOption(saveData, in.Data, copier.Option{
		CaseSensitive: true,
		DeepCopy:      true,
		Converters:    []copier.TypeConverter{},
	})
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	saveData.CreatorId = claim.OrganizationUserId
	saveData.OrganizationId = claim.OrganizationId
	saveData.CreatorName = claim.Username

	// todo check organization
	isExisted, err := l.svcCtx.BorrowerPropertyRepo.IsExisted(l.ctx, 0, claim.OrganizationId, &saveData.Address)
	if err != nil {
		return nil, err
	}
	if isExisted {
		err = errorx.WithStack(errorx.New("Address already exists"))
		return nil, err
	}

	result, err := l.svcCtx.BorrowerPropertyRepo.Create(l.ctx, saveData)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	err = copier.CopyWithOption(responseData, result, copier.Option{
		CaseSensitive: true,
		DeepCopy:      true,
		Converters: []copier.TypeConverter{
			copierx.TimeToRFC3339Str(),
			copierx.SqlTimePtrToStr(timex.LayoutUsStyle),
		},
	})
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	return &borrower.PropertyCreateResponse{
		Data: responseData,
	}, nil
}

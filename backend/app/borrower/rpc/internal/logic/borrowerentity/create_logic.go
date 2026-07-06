package borrowerentitylogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

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
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *CreateLogic) Create(in *borrower.BorrowerEntityCreateRequest) (resp *borrower.BorrowerEntityDetailResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
				logx.Field("service_name", l.svcCtx.Config.Name),
			).Error(err)
		}
	}()
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

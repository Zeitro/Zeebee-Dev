package borrowerapplicantlogic

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/stringx"
	"zeebee-backend/common/timex"
	"zeebee-backend/domains/zeebee/model"
)

type UpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateLogic) Update(in *borrower.BorrowerApplicantUpdateRequest) (resp *borrower.EmptyResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	var result *model.BorrowerApplicant
	// todo check organization
	result, err = l.svcCtx.BorrowerApplicantRepo.FindOneByColumns(l.ctx, map[string]interface{}{
		"id":      in.Data.Id,
		"user_id": in.Data.UserId,
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

	return &borrower.EmptyResponse{}, err
}

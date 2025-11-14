package borrowerpropertylogic

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/timex"
	"zeebee-backend/domains/zeebee/model"

	"zeebee-backend/common/errorx"
)

type CreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateLogic) Create(in *borrower.PropertyCreateRequest) (resp *borrower.PropertyCreateResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	saveData := &model.BorrowerProperty{}
	responseData := &borrower.BorrowerPropertyEntry{}
	err = copier.CopyWithOption(saveData, in.Data, copier.Option{
		CaseSensitive: true,
		DeepCopy:      true,
		Converters:    []copier.TypeConverter{},
	})
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	isExisted, err := l.svcCtx.BorrowerPropertyRepo.IsExisted(l.ctx, 0, saveData.BorrowerId, &saveData.Address)
	if err != nil {
		return nil, err
	}
	if isExisted {
		err = errors.WithStack(errors.New("Address Is Existed"))
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

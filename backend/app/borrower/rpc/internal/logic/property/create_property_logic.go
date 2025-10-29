package propertylogic

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/rpc/internal/repo/model"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/timex"

	"zeebee-backend/common/errorx"
)

type CreatePropertyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreatePropertyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePropertyLogic {
	return &CreatePropertyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreatePropertyLogic) CreateProperty(in *borrower.PropertyCreateRequest) (resp *borrower.PropertyCreateResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	saveData := &model.PropertyBase{}
	responseData := &borrower.PropertyBaseEntry{}
	err = copier.CopyWithOption(saveData, in.Data, copier.Option{
		CaseSensitive: true,
		DeepCopy:      true,
		Converters:    []copier.TypeConverter{},
	})
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	isExisted, err := l.svcCtx.PropertyBaseRepo.IsExisted(l.ctx, 0, saveData.BorrowerId, &saveData.Address)
	if err != nil {
		return nil, err
	}
	if isExisted {
		err = errors.WithStack(errors.New("Address Is Existed"))
		return nil, err
	}

	result, err := l.svcCtx.PropertyBaseRepo.Create(l.ctx, saveData)
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

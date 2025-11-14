package borrowerentitylogic

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/timex"
	"zeebee-backend/domains/zeebee/model"
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

func (l *CreateLogic) Create(in *borrower.BorrowerEntityCreateRequest) (resp *borrower.BorrowerEntityCreateResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	saveData := &model.BorrowerEntity{}
	responseData := &borrower.BorrowerEntityEntry{}
	err = copier.CopyWithOption(saveData, in.Data, copier.Option{
		CaseSensitive: true,
		DeepCopy:      true,
		Converters:    []copier.TypeConverter{},
	})
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	result, err := l.svcCtx.BorrowerEntityRepo.Create(l.ctx, saveData)
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

	return &borrower.BorrowerEntityCreateResponse{
		Data: responseData,
	}, nil
}

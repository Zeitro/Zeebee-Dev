package documentlogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/model"
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

func (l *CreateLogic) Create(in *borrower.DocumentCreateRequest) (resp *borrower.DocumentCreateResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	// validate parent exists and has leaf children (can only add siblings to existing leaves)
	parent, err := l.svcCtx.DocumentRepo.FindOneById(l.ctx, in.Data.ParentId)
	if err != nil {
		return nil, errorx.Wrap(err, "Parent document not found")
	}

	// check name uniqueness among siblings
	siblings, err := l.svcCtx.DocumentRepo.GetChildren(l.ctx, parent.Id)
	if err != nil {
		return nil, err
	}
	for _, s := range siblings {
		if s.Name == in.Data.Name {
			return nil, errorx.New("Document name already exists at this level")
		}
	}

	// create leaf document
	doc := &model.Document{
		ParentId:  parent.Id,
		Type:      "leaf",
		Name:      in.Data.Name,
		Level:     parent.Level + 1,
		SortOrder: int32(len(siblings)) + 1,
	}

	result, err := l.svcCtx.DocumentRepo.Create(l.ctx, doc)
	if err != nil {
		return nil, errorx.Wrap(err, "Failed to create document")
	}

	protoDoc, err := copierx.Trans[model.Document, borrower.DocumentEntity](result, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	return &borrower.DocumentCreateResponse{
		Data: protoDoc,
	}, nil
}

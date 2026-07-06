package documentlogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/timex"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/proto"
)

type DetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *DetailLogic) Detail(in *borrower.DocumentDetailRequest) (resp *borrower.DocumentDetailResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
				logx.Field("service_name", l.svcCtx.Config.Name),
			).Error(err)
		}
	}()

	entityDoc, err := l.svcCtx.EntityDocumentRepo.GetByEntity(l.ctx, constants.EntityDocumentType(in.Data.EntityType), in.Data.EntityId)
	if err != nil {
		return nil, errorx.Wrap(err, "Entity document not found")
	}

	rootDoc, err := l.buildDocumentTree(entityDoc.Id)
	if err != nil {
		return nil, err
	}

	protoDoc, err := l.convertDocumentEntity(rootDoc)
	if err != nil {
		return nil, err
	}

	return &borrower.DocumentDetailResponse{
		Data: protoDoc,
	}, nil
}

func (l *DetailLogic) buildDocumentTree(entityDocId int64) (*model.Document, error) {
	allDocs, err := l.svcCtx.DocumentRepo.GetAllByEntityDocumentId(l.ctx, entityDocId)
	if err != nil {
		return nil, errorx.Wrap(err, "Failed to load documents")
	}

	docMap := make(map[int64]*model.Document)
	var root *model.Document
	for _, d := range allDocs {
		docMap[d.Id] = d
		if d.ParentId == 0 {
			root = d
		}
	}
	if root == nil {
		return nil, errorx.New("Root document not found")
	}

	// load files for all documents
	fileEntityType := string(constants.FileEntityTypeDocument)
	for _, d := range allDocs {
		var files []*model.File
		err := l.svcCtx.FileRepo.Client(l.ctx).
			Where("entity_type = ? AND entity_id = ?", fileEntityType, d.Id).
			Find(&files).Error
		if err == nil {
			d.Files = files
		}
	}

	// build parent-child tree
	for _, d := range allDocs {
		if d.ParentId == 0 {
			continue
		}
		if parent, ok := docMap[d.ParentId]; ok {
			parent.Children = append(parent.Children, d)
		}
	}

	return root, nil
}

func (l *DetailLogic) convertDocumentEntity(doc *model.Document) (*proto.DocumentEntity, error) {
	entity, err := copierx.Trans[model.Document, proto.DocumentEntity](
		doc,
		copierx.TimeToRFC3339Str(),
		copierx.SqlTimePtrToStr(timex.LayoutUsStyle),
	)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	// metadata is []byte in proto, copierx can't auto-convert jsonb
	if doc.Metadata != nil {
		entity.Metadata = string(doc.Metadata)
	}

	// convert files
	for _, f := range doc.Files {
		protoFile, err := copierx.Trans[model.File, proto.FileEntity](
			f,
			copierx.TimeToRFC3339Str(),
		)
		if err != nil {
			return nil, errorx.GetFailedTransRespErr(err)
		}
		entity.Files = append(entity.Files, protoFile)
	}

	// convert children recursively
	for _, child := range doc.Children {
		childEntity, err := l.convertDocumentEntity(child)
		if err != nil {
			return nil, err
		}
		entity.Children = append(entity.Children, childEntity)
	}

	return entity, nil
}

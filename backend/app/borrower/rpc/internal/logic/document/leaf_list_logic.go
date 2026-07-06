package documentlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/domains/zeebee/proto"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
)

type LeafListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewLeafListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LeafListLogic {
	return &LeafListLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *LeafListLogic) LeafList(in *borrower.DocumentLeafListRequest) (resp *borrower.DocumentLeafListResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
				logx.Field("service_name", l.svcCtx.Config.Name),
			).Error(err)
		}
	}()

	// find entity_document to get root
	entityDoc, err := l.svcCtx.EntityDocumentRepo.GetByEntity(l.ctx, constants.EntityDocumentType(in.Data.EntityType), in.Data.EntityId)
	if err != nil {
		return nil, errorx.Wrap(err, "Entity document not found")
	}

	// get all descendants
	descendants, err := l.svcCtx.DocumentRepo.GetDescendants(l.ctx, []int64{entityDoc.RootDocumentId})
	if err != nil {
		return nil, errorx.Wrap(err, "Failed to load documents")
	}

	// filter leaf nodes
	leaves := make([]*proto.DocumentEntity, 0)
	fileEntityType := string(constants.FileEntityTypeDocument)
	for _, doc := range descendants {
		if doc.Type != "leaf" {
			continue
		}

		// load files count
		var fileCount int64
		l.svcCtx.FileRepo.Client(l.ctx).
			Where("entity_type = ? AND entity_id = ?", fileEntityType, doc.Id).
			Count(&fileCount)

		protoDoc, innerErr := copierx.Trans[model.Document, proto.DocumentEntity](
			doc,
			copierx.TimeToRFC3339Str(),
		)
		if innerErr != nil {
			continue
		}

		leaves = append(leaves, protoDoc)
	}

	return &borrower.DocumentLeafListResponse{
		Data: leaves,
	}, nil
}

package documentlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/proto"
)

type FilesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewFilesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FilesLogic {
	return &FilesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *FilesLogic) Files(in *borrower.DocumentFilesRequest) (resp *borrower.DocumentFilesResponse, err error) {

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx, constants.UserTypeBorrower)
	if err != nil {
		return nil, err
	}

	// get document with files preloaded
	doc, err := l.svcCtx.DocumentRepo.FindOneWithFiles(l.ctx, in.Data.DocumentId)
	if err != nil {
		err = errorx.Wrap(err, "failed to find document")
		return
	}

	// verify org ownership: walk up to root, find entity_document, check transaction org
	if err = l.verifyOrgOwnership(doc, claim.OrganizationId); err != nil {
		return nil, err
	}

	// convert files to proto
	result := make([]*proto.FileEntity, 0, len(doc.Files))
	for _, f := range doc.Files {
		protoFile, innerErr := copierx.Trans[model.File, proto.FileEntity](f, copierx.TimeToRFC3339Str())
		if innerErr != nil {
			continue
		}
		result = append(result, protoFile)
	}

	return &borrower.DocumentFilesResponse{
		Data: result,
	}, nil
}

func (l *FilesLogic) verifyOrgOwnership(doc *model.Document, orgId int64) (err error) {
	var entityDoc *model.EntityDocument

	entityDoc, err = l.svcCtx.EntityDocumentRepo.FindOneById(l.ctx, doc.EntityDocumentId)
	if err != nil {
		return errorx.Wrap(err, "document entity binding not found")
	}

	if entityDoc.OrganizationId != orgId {
		return errorx.New("document not belongs to your organization")
	}

	return nil
}

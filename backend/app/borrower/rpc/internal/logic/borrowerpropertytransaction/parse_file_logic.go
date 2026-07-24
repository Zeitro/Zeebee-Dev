package borrowerpropertytransactionlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/biz_helper/parse_file"
	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/domains/zeebee/constants"

	"zeebee-backend/common/errorx"
)

type ParseFileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewParseFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ParseFileLogic {
	return &ParseFileLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

// ParseFile triggers a manual parse (manual_tree / manual_select). IDOR
// + rule resolution live in parse_file.Dispatch; here we only resolve the
// caller's org from the claim. A duplicate trigger surfaces as gRPC
// AlreadyExists (mapped to HTTP 409 at the API layer).
func (l *ParseFileLogic) ParseFile(in *borrower.ParseFileRequest) (resp *borrower.ParseFileResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	// check organization,transaction and document

	document, err := l.svcCtx.DocumentRepo.FindOneById(l.ctx, in.Data.DocumentId)
	if err != nil {
		return nil, err
	}
	entityDocument, err := l.svcCtx.EntityDocumentRepo.FindOneById(l.ctx, document.EntityDocumentId)
	if err != nil {
		return nil, err
	}

	if entityDocument.OrganizationId != claim.OrganizationId {
		err = errorx.New("document does not belong to your organization")
		return
	}
	if entityDocument.EntityId != in.Data.Id || entityDocument.EntityType != constants.EntityDocumentTypeBorrowerPropertyTransaction {
		err = errorx.New("document dose not belong to the transaction")
		return
	}

	requestId, err := l.svcCtx.FileParser.Dispatch(l.ctx, parse_file.DispatchParam{
		TransactionId:  in.Data.Id,
		DocumentId:     in.Data.DocumentId,
		OrganizationId: claim.OrganizationId,
		Parser:         in.Data.Parser,
		Trigger:        parse_file.TriggerManual,
	})
	if err != nil {
		return nil, err
	}

	return &borrower.ParseFileResponse{
		Data: &borrower.ParseFileResponseData{
			RequestId: requestId,
		},
	}, nil
}

package borrowerpropertytransactionlogic

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/timex"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
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

func (l *CreateLogic) Create(in *borrower.BorrowerPropertyTransactionCreateRequest) (resp *borrower.BorrowerPropertyTransactionCreateResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	// validate property exists and belongs to organization
	_, err = l.svcCtx.BorrowerPropertyRepo.GetOneByIdForOrganization(l.ctx, in.Data.BorrowerPropertyId, claim.OrganizationId)
	if err != nil {
		return nil, errorx.Wrap(err, "property not found")
	}

	// validate transaction type
	if !constants.BorrowerPropertyTransactionType(in.Data.Type).IsValid() {
		return nil, errorx.New("invalid transaction type")
	}

	// generate code: tc + year + base62(ms timestamp)
	code := generateTransactionCode()

	// build transaction model
	transactionData := &model.BorrowerPropertyTransaction{
		BorrowerPropertyId: in.Data.BorrowerPropertyId,
		OrganizationId:     claim.OrganizationId,
		Code:               code,
		Type:               constants.BorrowerPropertyTransactionType(in.Data.Type),
		Status:             constants.BorrowerPropertyTransactionStatusPending,
	}

	// load document template
	template, err := l.svcCtx.DocumentTemplateRepo.GetByType(l.ctx, constants.DocumentTemplateTypeBorrowerPropertyTransaction)
	if err != nil {
		return nil, errorx.Wrap(err, "document template not found")
	}

	var result *model.BorrowerPropertyTransaction

	// execute in transaction
	err = l.svcCtx.BorrowerPropertyTransactionRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		var innerErr error

		// create transaction
		result, innerErr = l.svcCtx.BorrowerPropertyTransactionRepo.CreateWithTx(l.ctx, tx, transactionData, false)
		if innerErr != nil {
			return innerErr
		}

		// create entity_document binding first (root_document_id=0, backfilled later)
		entityDoc := &model.EntityDocument{
			EntityType:     constants.EntityDocumentTypeBorrowerPropertyTransaction,
			OrganizationId: claim.OrganizationId,
			EntityId:       result.Id,
			RootDocumentId: 0,
		}
		entityDoc, innerErr = l.svcCtx.EntityDocumentRepo.CreateWithTx(l.ctx, tx, entityDoc, false)
		if innerErr != nil {
			return innerErr
		}

		// build document tree with entity_document_id
		rootDoc, innerErr := l.createDocumentTree(tx, entityDoc.Id, template.Structure, string(result.Type))
		if innerErr != nil {
			return innerErr
		}

		// backfill root_document_id
		innerErr = l.svcCtx.EntityDocumentRepo.UpdateColumnsByIdWithTx(l.ctx, tx, entityDoc.Id, map[string]interface{}{
			"root_document_id": rootDoc.Id,
		})
		if innerErr != nil {
			return innerErr
		}

		// initialize checklist
		checklist, innerErr := l.baseLogic.RecalculateBorrowerPropertyTransactionChecklist(l.ctx, tx, result.Id)
		if innerErr != nil {
			return innerErr
		}

		result.Checklist = checklist

		metadata, _ := json.Marshal(map[string]interface{}{
			"after": struct {
				*model.BorrowerPropertyTransaction
				Id        int64     `json:"-"`
				CreatedAt time.Time `json:"-"`
				UpdatedAt time.Time `json:"-"`
			}{
				BorrowerPropertyTransaction: result,
			},
		})

		// create operation log
		innerErr = l.baseLogic.OperationLogic.CreateOperationLog(l.ctx, tx, dto.OperationCreateDto{
			EntityId:       result.Id,
			EntityType:     constants.OperationEntityTypeBorrowerPropertyTransaction,
			OperatorId:     claim.OrganizationUserId,
			OperatorName:   claim.Username,
			OperatorType:   claim.UserType,
			Operation:      constants.OperationActionBorrowerPropertyTransactionCreate,
			OrganizationId: claim.OrganizationId,
			Metadata:       metadata,
		})
		return innerErr
	})
	if err != nil {
		return nil, errorx.Wrap(err, "failed to create transaction")
	}

	// convert response
	responseEntity, err := copierx.Trans[model.BorrowerPropertyTransaction, borrower.BorrowerPropertyTransactionEntity](
		result,
		copierx.TimeToRFC3339Str(),
		copierx.SqlTimePtrToStr(timex.LayoutUsStyle),
	)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	return &borrower.BorrowerPropertyTransactionCreateResponse{
		Data: responseEntity,
	}, nil
}

func (l *CreateLogic) createDocumentTree(tx *gorm.DB, entityDocumentId int64, structureMap model.DocumentTemplateStructureMap, transactionType string) (*model.Document, error) {
	// merge common + per-type nodes
	nodes := structureMap["common"]
	typeKey := strings.ToLower(transactionType)
	if typeNodes, ok := structureMap[typeKey]; ok {
		nodes = append(nodes, typeNodes...)
	}

	root := &model.Document{
		ParentId:         0,
		EntityDocumentId: entityDocumentId,
		Type:             constants.DocumentTypeRoot,
		Name:             "Root",
		Level:            0,
		SortOrder:        0,
	}
	root, err := l.svcCtx.DocumentRepo.CreateWithTx(l.ctx, tx, root, false)
	if err != nil {
		return nil, err
	}

	for _, node := range nodes {
		if err := l.createDocumentNode(tx, entityDocumentId, root.Id, node, 1); err != nil {
			return nil, err
		}
	}

	return root, nil
}

func (l *CreateLogic) createDocumentNode(tx *gorm.DB, entityDocumentId int64, parentId int64, node model.DocumentTemplateStructure, level int32) error {
	docType := constants.DocumentTypeBranch
	if len(node.Children) == 0 {
		docType = constants.DocumentTypeLeaf
	}

	metadata, _ := json.Marshal(map[string]interface{}{
		"is_required": node.IsRequired,
		"status":      constants.BorrowerPropertyTransactionStatusPending,
		"create_type": "System",
	})

	doc := &model.Document{
		ParentId:         parentId,
		EntityDocumentId: entityDocumentId,
		Type:             docType,
		Name:             node.Name,
		Level:            level,
		SortOrder:        node.SortOrder,
		Metadata:         metadata,
	}

	doc, err := l.svcCtx.DocumentRepo.CreateWithTx(l.ctx, tx, doc, false)
	if err != nil {
		return err
	}

	for _, child := range node.Children {
		if err := l.createDocumentNode(tx, entityDocumentId, doc.Id, child, level+1); err != nil {
			return err
		}
	}

	return nil
}

func generateTransactionCode() string {
	now := time.Now()
	year := now.Format("06")
	ms := now.UnixMilli()
	base62 := encodeBase62(ms)
	return fmt.Sprintf("TC%s%s", year, base62)
}

func encodeBase62(n int64) string {
	const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	if n == 0 {
		return "0"
	}
	result := make([]byte, 0)
	for n > 0 {
		result = append(result, charset[n%62])
		n /= 62
	}
	// reverse
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}
	return string(result)
}

package borrowerpropertytransactionlogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/repo"

	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/timex"
	"zeebee-backend/domains/zeebee/model"
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

func (l *DetailLogic) Detail(in *borrower.BorrowerPropertyTransactionDetailRequest) (resp *borrower.BorrowerPropertyTransactionDetailResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	result, err := l.svcCtx.BorrowerPropertyTransactionRepo.GetOrganizationTransaction(l.ctx, in.Data.Id, claim.OrganizationId)
	if err != nil {
		return nil, errorx.Wrap(err, "transaction not found")
	}

	fillChecklistNodeFiles(l.ctx, l.svcCtx.FileRepo, &result.Checklist)

	entity, err := copierx.Trans[model.BorrowerPropertyTransaction, borrower.BorrowerPropertyTransactionEntity](
		result,
		copierx.TimeToRFC3339Str(),
		copierx.SqlTimePtrToStr(timex.LayoutUsStyle),
	)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	// fill files into checklist leaf nodes

	return &borrower.BorrowerPropertyTransactionDetailResponse{
		Data: entity,
	}, nil
}

// fillChecklistNodeFiles recursively fills files for leaf checklist nodes
func fillChecklistNodeFiles(ctx context.Context, fileRepo *repo.FileRepo, checklist *dto.ChecklistRoot) {

	fillNodeFiles(ctx, fileRepo, checklist.RequiredDocuments)
	fillNodeFiles(ctx, fileRepo, checklist.OptionalDocuments)
}

func fillNodeFiles(ctx context.Context, fileRepo *repo.FileRepo, nodes []*dto.ChecklistNode) {
	for i, node := range nodes {
		if len(node.Children) == 0 {
			// leaf node: query files for this document
			files, _ := fileRepo.GetFiles(ctx, node.DocumentId, constants.FileEntityTypeDocument)

			// 整段切片一次性转换：天然支持同一文档下的多个文件；
			// 额外带上 SqlTimePtrToStr 以填充 drive_synced_at（Google Drive 标签需要）。
			dtoFiles, _ := copierx.Trans[[]*model.File, []*dto.FileEntityDto](&files, copierx.TimeToRFC3339Str(), copierx.SqlTimePtrToStr(timex.LayoutUsStyle))
			if dtoFiles != nil && len(*dtoFiles) > 0 {
				nodes[i].FileCount = int32(len(*dtoFiles))
				nodes[i].Files = *dtoFiles
			}
		} else {
			fillNodeFiles(ctx, fileRepo, node.Children)
		}
	}
}

package borrowerpropertytransactionlogic

import (
	"context"

	"github.com/duke-git/lancet/v2/slice"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/proto"
	"zeebee-backend/domains/zeebee/repo"

	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/timex"
	"zeebee-backend/domains/zeebee/model"
)

type DetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *DetailLogic) Detail(in *proto.PathIdRequest) (resp *borrower.BorrowerPropertyTransactionDetailResponse, err error) {
	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	result, err := l.svcCtx.BorrowerPropertyTransactionRepo.GetOrganizationTransaction(l.ctx, in.Data.Id, claim.OrganizationId)
	if err != nil {
		return nil, errorx.Wrap(err, "transaction not found")
	}

	// fill files into checklist leaf nodes
	fillChecklistNodeFiles(l.ctx, l.svcCtx.FileRepo, &result.Checklist)

	entity, err := copierx.Trans[model.BorrowerPropertyTransaction, proto.BorrowerPropertyTransactionEntity](
		result,
		copierx.TimeToRFC3339Str(),
		copierx.SqlTimePtrToStr(timex.LayoutUsStyle),
	)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	// load property transaction escrows
	escrows, err := l.svcCtx.PropertyTransactionEscrowRepo.FindAll(l.ctx, map[string]any{
		"borrower_property_transaction_id": in.Data.Id,
	})
	if err != nil {
		return nil, errorx.Wrap(err, "fail to load escrows")
	}
	if len(escrows) > 0 {
		protoEscrows, terr := copierx.Trans[[]*model.PropertyTransactionEscrow, []*proto.PropertyTransactionEscrowEntity](&escrows, copierx.TimeToRFC3339Str())
		if terr != nil {
			return nil, errorx.GetFailedTransRespErr(terr)
		}
		entity.Escrows = *protoEscrows
	}
	slice.SortBy(entity.Escrows, func(a, b *proto.PropertyTransactionEscrowEntity) bool {
		return a.Id < b.Id
	})

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
			files, err := fileRepo.GetFiles(ctx, node.DocumentId, constants.FileEntityTypeDocument)
			if err != nil {
				continue
			}
			dtoFiles := toFileEntityDtos(files)
			if len(dtoFiles) > 0 {
				nodes[i].FileCount = int32(len(dtoFiles))
				nodes[i].Files = dtoFiles
			}
		} else {
			fillNodeFiles(ctx, fileRepo, node.Children)
		}
	}
}

// toFileEntityDtos converts file models into checklist DTOs. DriveSyncedAt is a
// *sql.NullTime, so SqlTimePtrToStr is required for the Google Drive sync marker
// to survive the conversion; without it drive_synced_at is silently dropped and
// the UI "Google Drive" label never shows.
func toFileEntityDtos(files []*model.File) []*dto.FileEntityDto {
	dtoFiles, _ := copierx.Trans[[]*model.File, []*dto.FileEntityDto](
		&files,
		copierx.TimeToRFC3339Str(),
		copierx.SqlTimePtrToStr(timex.LayoutUsStyle),
	)
	if dtoFiles == nil {
		return nil
	}
	return *dtoFiles
}

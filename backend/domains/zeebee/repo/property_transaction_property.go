package repo

import (
	"context"

	"gorm.io/gorm"

	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/model"
)

// PropertyTransactionPropertyRepo relies on gormx.BaseRepo for find/create/update/
// delete. Only the "delete by transaction, id NOT IN set" operation (full-replace
// sync) has no BaseRepo equivalent and is defined here.
type PropertyTransactionPropertyRepo struct {
	*gormx.BaseRepo[*model.PropertyTransactionProperty]
}

func NewPropertyTransactionPropertyRepo(db *gorm.DB) *PropertyTransactionPropertyRepo {
	return &PropertyTransactionPropertyRepo{
		BaseRepo: gormx.NewBaseRepo[*model.PropertyTransactionProperty](db),
	}
}

// DeleteByIdsNotInWithTx removes Propertys of a transaction whose id is not in
// keepIds (full-replace sync). Empty keepIds deletes all rows for the transaction.
func (repo *PropertyTransactionPropertyRepo) DeleteByIdsNotInWithTx(
	ctx context.Context, tx *gorm.DB, transactionId int64, keepIds []int64,
) error {

	q := tx.WithContext(ctx).
		Where("borrower_property_transaction_id = ?", transactionId)
	if len(keepIds) > 0 {
		q = q.Where("id NOT IN ?", keepIds)
	}
	return errorx.WithStack(q.Delete(&model.PropertyTransactionProperty{}).Error)
}

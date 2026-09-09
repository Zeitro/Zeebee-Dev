package repo

import (
	"context"

	"gorm.io/gorm"

	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/model"
)

// PropertyTransactionEscrowRepo relies on gormx.BaseRepo for find/create/update/
// delete. Only the "delete by transaction, id NOT IN set" operation (full-replace
// sync) has no BaseRepo equivalent and is defined here.
type PropertyTransactionEscrowRepo struct {
	*gormx.BaseRepo[*model.PropertyTransactionEscrow]
}

func NewPropertyTransactionEscrowRepo(db *gorm.DB) *PropertyTransactionEscrowRepo {
	return &PropertyTransactionEscrowRepo{
		BaseRepo: gormx.NewBaseRepo[*model.PropertyTransactionEscrow](db),
	}
}

// DeleteByIdsNotInWithTx removes escrows of a transaction whose id is not in
// keepIds (full-replace sync). Empty keepIds deletes all rows for the transaction.
func (repo *PropertyTransactionEscrowRepo) DeleteByIdsNotInWithTx(
	ctx context.Context, tx *gorm.DB, transactionId int64, keepIds []int64,
) error {

	q := tx.WithContext(ctx).
		Where("borrower_property_transaction_id = ?", transactionId)
	if len(keepIds) > 0 {
		q = q.Where("id NOT IN ?", keepIds)
	}
	return errorx.WithStack(q.Delete(&model.PropertyTransactionEscrow{}).Error)
}

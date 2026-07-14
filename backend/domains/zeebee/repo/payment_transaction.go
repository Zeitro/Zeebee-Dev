package repo

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/model"
)

type PaymentTransactionRepo struct {
	*gormx.BaseRepo[*model.PaymentTransaction]
}

func NewPaymentTransactionRepo(db *gorm.DB) *PaymentTransactionRepo {
	return &PaymentTransactionRepo{
		BaseRepo: gormx.NewBaseRepo[*model.PaymentTransaction](db),
	}
}

// UpsertByInvoiceIdWithTx inserts the transaction or, when the invoice already
// exists (replayed webhook), refreshes its mutable fields.
func (repo *PaymentTransactionRepo) UpsertByInvoiceIdWithTx(ctx context.Context, tx *gorm.DB, transaction *model.PaymentTransaction) error {
	err := tx.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "stripe_invoice_id"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"organization_subscription_id", "amount", "currency", "status",
				"period_start", "period_end", "paid_at", "updated_at",
			}),
		}).
		Create(transaction).Error
	return errorx.WithStack(err)
}

func (repo *PaymentTransactionRepo) FindByOrganizationId(ctx context.Context, organizationId int64) (result []*model.PaymentTransaction, err error) {
	err = repo.Client(ctx).
		Where("organization_id = ?", organizationId).
		Order("id desc").
		Find(&result).Error
	return result, errorx.WithStack(err)
}

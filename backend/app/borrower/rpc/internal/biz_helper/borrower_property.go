package biz_helper

import (
	"context"

	"github.com/lib/pq"
	"gorm.io/gorm"

	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/repo"
)

type BorrowerProperty struct {
	repos *repo.AllRepos
}

func NewBorrowerProperty(repos *repo.AllRepos) *BorrowerProperty {
	return &BorrowerProperty{
		repos: repos,
	}
}

// MapTransactionToPropertyStatus maps transaction type+status to property status.
// Returns empty string and false if the combination is not in the mapping table (e.g. Cancelled).
func MapTransactionToPropertyStatus(txType constants.BorrowerPropertyTransactionType, txStatus constants.BorrowerPropertyTransactionStatus) (constants.BorrowerPropertyStatus, bool) {
	type mappingKey struct {
		t constants.BorrowerPropertyTransactionType
		s constants.BorrowerPropertyTransactionStatus
	}
	mapping := map[mappingKey]constants.BorrowerPropertyStatus{
		{constants.BorrowerPropertyTransactionTransactionTypePurchase, constants.BorrowerPropertyTransactionStatusPending}:    constants.BorrowerPropertyStatusPurchasePending,
		{constants.BorrowerPropertyTransactionTransactionTypePurchase, constants.BorrowerPropertyTransactionStatusClosed}:     constants.BorrowerPropertyStatusActive,
		{constants.BorrowerPropertyTransactionTransactionTypeReFinance, constants.BorrowerPropertyTransactionStatusClosed}:    constants.BorrowerPropertyStatusActive,
		{constants.BorrowerPropertyTransactionTransactionTypeReFinance, constants.BorrowerPropertyTransactionStatusPending}:   constants.BorrowerPropertyStatusReFinancePending,
		{constants.BorrowerPropertyTransactionTransactionTypeSale, constants.BorrowerPropertyTransactionStatusPending}:        constants.BorrowerPropertyStatusSalePending,
		{constants.BorrowerPropertyTransactionTransactionTypeSale, constants.BorrowerPropertyTransactionStatusClosed}:         constants.BorrowerPropertyStatusSold,
		{constants.BorrowerPropertyTransactionTransactionTypeChangeTitle, constants.BorrowerPropertyTransactionStatusPending}: constants.BorrowerPropertyStatusChangeTitlePending,
		{constants.BorrowerPropertyTransactionTransactionTypeAssignment, constants.BorrowerPropertyTransactionStatusPending}:  constants.BorrowerPropertyStatusAssignmentPending,
		{constants.BorrowerPropertyTransactionTransactionTypeAssignment, constants.BorrowerPropertyTransactionStatusClosed}:   constants.BorrowerPropertyStatusSold,
	}

	if v, ok := mapping[mappingKey{txType, txStatus}]; ok {
		return v, true
	}
	return "", false
}

// CalculatePropertyStatuses queries all transactions for a property, maps each to a property status,
// deduplicates, and updates the property's status column.
func (h *BorrowerProperty) CalculatePropertyStatuses(ctx context.Context, tx *gorm.DB, propertyId int64) error {
	transactions, err := h.repos.BorrowerPropertyTransactionRepo.FindByBorrowerPropertyIdWithTx(ctx, tx, propertyId)
	if err != nil {
		return errorx.Wrap(err, "failed to query transactions for property status calculation")
	}

	seen := make(map[constants.BorrowerPropertyStatus]struct{})
	var statuses pq.StringArray
	for _, t := range transactions {
		if status, ok := MapTransactionToPropertyStatus(t.Type, t.Status); ok {
			if _, exists := seen[status]; !exists {
				seen[status] = struct{}{}
				statuses = append(statuses, string(status))
			}
		}
	}

	if statuses == nil {
		statuses = pq.StringArray{}
	}

	return h.repos.BorrowerPropertyRepo.UpdateColumnsByIdWithTx(ctx, tx, propertyId, map[string]interface{}{
		"status": statuses,
	})
}

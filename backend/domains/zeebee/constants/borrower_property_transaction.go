package constants

import (
	"zeebee-backend/common/enums"
)

type BorrowerPropertyTransactionType string

const BorrowerPropertyTransactionTypeName = "borrower_property_transaction_type"

const (
	BorrowerPropertyTransactionTransactionTypePurchase    BorrowerPropertyTransactionType = "Purchase"
	BorrowerPropertyTransactionTransactionTypeSale        BorrowerPropertyTransactionType = "Sale"
	BorrowerPropertyTransactionTransactionTypeReFinance   BorrowerPropertyTransactionType = "Re-finance"
	BorrowerPropertyTransactionTransactionTypeChangeTitle BorrowerPropertyTransactionType = "Change Title"
	BorrowerPropertyTransactionTransactionTypeAssignment  BorrowerPropertyTransactionType = "Assignment"
)

func (b BorrowerPropertyTransactionType) IsValid() bool {
	switch b {
	case BorrowerPropertyTransactionTransactionTypePurchase,
		BorrowerPropertyTransactionTransactionTypeSale, BorrowerPropertyTransactionTransactionTypeReFinance, BorrowerPropertyTransactionTransactionTypeChangeTitle, BorrowerPropertyTransactionTransactionTypeAssignment:
		return true
	default:
		return false

	}
}

type BorrowerPropertyTransactionStatus string

const BorrowerPropertyTransactionStatusName = "borrower_property_transaction_status"

const (
	BorrowerPropertyTransactionStatusPending   BorrowerPropertyTransactionStatus = "Pending"
	BorrowerPropertyTransactionStatusClosed    BorrowerPropertyTransactionStatus = "Closed"
	BorrowerPropertyTransactionStatusCancelled BorrowerPropertyTransactionStatus = "Cancelled"
)

func init() {
	enums.Register(BorrowerPropertyTransactionTypeName, []enums.EnumOpt{
		{Label: string(BorrowerPropertyTransactionTransactionTypePurchase), Value: string(BorrowerPropertyTransactionTransactionTypePurchase)},
		{Label: string(BorrowerPropertyTransactionTransactionTypeSale), Value: string(BorrowerPropertyTransactionTransactionTypeSale)},
		{Label: string(BorrowerPropertyTransactionTransactionTypeReFinance), Value: string(BorrowerPropertyTransactionTransactionTypeReFinance)},
		{Label: string(BorrowerPropertyTransactionTransactionTypeChangeTitle), Value: string(BorrowerPropertyTransactionTransactionTypeChangeTitle)},
		{Label: string(BorrowerPropertyTransactionTransactionTypeAssignment), Value: string(BorrowerPropertyTransactionTransactionTypeAssignment)},
	})

	enums.Register(BorrowerPropertyTransactionStatusName, []enums.EnumOpt{
		{Label: string(BorrowerPropertyTransactionStatusPending), Value: string(BorrowerPropertyTransactionStatusPending)},
		{Label: string(BorrowerPropertyTransactionStatusClosed), Value: string(BorrowerPropertyTransactionStatusClosed)},
		{Label: string(BorrowerPropertyTransactionStatusCancelled), Value: string(BorrowerPropertyTransactionStatusCancelled)},
	})
}

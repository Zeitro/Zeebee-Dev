package constants

import "zeebee-backend/common/enums"

// stage: Pending,Active,Archived
type LoanStage string

const LoanStageName = "loan_stage"
const (
	LoanStagePending  LoanStage = "Pending"
	LoanStageActive   LoanStage = "Active"
	LoanStageArchived LoanStage = "Archived"
)

// pending_status: Quoting,Processing,Closing
type LoanPendingStatus string

const LoanPendingStatusName = "loan_pending_status"
const (
	LoanPendingStatusQuoting    LoanPendingStatus = "Quoting"
	LoanPendingStatusProcessing LoanPendingStatus = "Processing"
	LoanPendingStatusClosing    LoanPendingStatus = "Closing"
)

// purpose: Purchase,Rate-term refi,Cash-out refi
type LoanPurpose string

const LoanPurposeName = "loan_purpose"
const (
	LoanPurposePurchase LoanPurpose = "Purchase"
	LoanPurposeRateTerm LoanPurpose = "Rate-term refi"
	LoanPurposeCashOut  LoanPurpose = "Cash-out refi"
)

// loan_type: Bridge,Rehab,Construction
type LoanType string

const LoanTypeName = "loan_type"
const (
	LoanTypeBridge       LoanType = "Bridge"
	LoanTypeRehab        LoanType = "Rehab"
	LoanTypeConstruction LoanType = "Construction"
)

// source: borrower,lender,system
type LoanSource string

const LoanSourceName = "loan_source"
const (
	LoanSourceBorrower LoanSource = "borrower"
	LoanSourceLender   LoanSource = "lender"
	LoanSourceSystem   LoanSource = "system"
)

type LienPosition string

const LienPositionName = "lien_position"
const (
	LienPosition1     LienPosition = "1st"
	LienPosition2     LienPosition = "2nd"
	LienPosition3     LienPosition = "3rd"
	LienPosition4     LienPosition = "4th"
	LienPositionOther LienPosition = "Other"
)

func init() {
	// register loan stage
	enums.Register(LoanStageName, []enums.EnumOpt{
		{
			Label: string(LoanStagePending),
			Value: LoanStagePending,
		},
		{
			Label: string(LoanStageActive),
			Value: LoanStageActive,
		},
		{
			Label: string(LoanStageArchived),
			Value: LoanStageArchived,
		},
	})

	// register loan pending status
	enums.Register(LoanPendingStatusName, []enums.EnumOpt{
		{
			Label: string(LoanPendingStatusQuoting),
			Value: LoanPendingStatusQuoting,
		},
		{
			Label: string(LoanPendingStatusProcessing),
			Value: LoanPendingStatusProcessing,
		},
		{
			Label: string(LoanPendingStatusClosing),
			Value: LoanPendingStatusClosing,
		},
	})
	// register loan purpose
	enums.Register(LoanPurposeName, []enums.EnumOpt{
		{
			Label: string(LoanPurposePurchase),
			Value: LoanPurposePurchase,
		},
		{
			Label: string(LoanPurposeRateTerm),
			Value: LoanPurposeRateTerm,
		},
		{
			Label: string(LoanPurposeCashOut),
			Value: LoanPurposeCashOut,
		},
	})
	// register loan type
	enums.Register(LoanTypeName, []enums.EnumOpt{
		{
			Label: string(LoanTypeBridge),
			Value: LoanTypeBridge,
		},
		{
			Label: string(LoanTypeRehab),
			Value: LoanTypeRehab,
		},
		{
			Label: string(LoanTypeConstruction),
			Value: LoanTypeConstruction,
		},
	})
	// register loan source
	enums.Register(LoanSourceName, []enums.EnumOpt{
		{
			Label: string(LoanSourceBorrower),
			Value: LoanSourceBorrower,
		},
		{
			Label: string(LoanSourceLender),
			Value: LoanSourceLender,
		},
		{
			Label: string(LoanSourceSystem),
			Value: LoanSourceSystem,
		},
	})
	// register lien position
	enums.Register(LienPositionName, []enums.EnumOpt{
		{
			Label: string(LienPosition1),
			Value: LienPosition1,
		},
		{
			Label: string(LienPosition2),
			Value: LienPosition2,
		},
		{
			Label: string(LienPosition3),
			Value: LienPosition3,
		},
		{
			Label: string(LienPosition4),
			Value: LienPosition4,
		},
		{
			Label: string(LienPositionOther),
			Value: LienPositionOther,
		},
	})

}

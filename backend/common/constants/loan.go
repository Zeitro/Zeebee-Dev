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
	LoanTypeConstruction LoanType = "Ground up construction"
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

// business_strategy: Sale,Rental,Flip,ADU/Addition,Single Lot SFR,SB9 lot split,Conforming lot split,New Condo/TH,Community development (>3 units),Commercial
type BusinessStrategy string

const BusinessStrategyName = "business_strategy"
const (
	BusinessStrategySale                 BusinessStrategy = "Sale"
	BusinessStrategyRental               BusinessStrategy = "Rental"
	BusinessStrategyFlip                 BusinessStrategy = "Flip"
	BusinessStrategyADUAddition          BusinessStrategy = "ADU/Addition"
	BusinessStrategySingleLotSFR         BusinessStrategy = "Single Lot SFR"
	BusinessStrategySB9LotSplit          BusinessStrategy = "SB9 lot split"
	BusinessStrategyConformingLotSplit   BusinessStrategy = "Conforming lot split"
	BusinessStrategyNewCondoTH           BusinessStrategy = "New Condo/TH"
	BusinessStrategyCommunityDevelopment BusinessStrategy = "Community development (>3 units)"
	BusinessStrategyCommercial           BusinessStrategy = "Commercial"
)

func init() {
	// register loan stage
	enums.Register(LoanStageName, []enums.EnumOpt{
		{
			Label: string(LoanStagePending),
			Value: string(LoanStagePending),
		},
		{
			Label: string(LoanStageActive),
			Value: string(LoanStageActive),
		},
		{
			Label: string(LoanStageArchived),
			Value: string(LoanStageArchived),
		},
	})

	// register loan pending status
	enums.Register(LoanPendingStatusName, []enums.EnumOpt{
		{
			Label: string(LoanPendingStatusQuoting),
			Value: string(LoanPendingStatusQuoting),
		},
		{
			Label: string(LoanPendingStatusProcessing),
			Value: string(LoanPendingStatusProcessing),
		},
		{
			Label: string(LoanPendingStatusClosing),
			Value: string(LoanPendingStatusClosing),
		},
	})
	// register loan purpose
	enums.Register(LoanPurposeName, []enums.EnumOpt{
		{
			Label: string(LoanPurposePurchase),
			Value: string(LoanPurposePurchase),
		},
		{
			Label: string(LoanPurposeRateTerm),
			Value: string(LoanPurposeRateTerm),
		},
		{
			Label: string(LoanPurposeCashOut),
			Value: string(LoanPurposeCashOut),
		},
	})
	// register loan type
	enums.Register(LoanTypeName, []enums.EnumOpt{
		{
			Label: string(LoanTypeBridge),
			Value: string(LoanTypeBridge),
		},
		{
			Label: string(LoanTypeRehab),
			Value: string(LoanTypeRehab),
		},
		{
			Label: string(LoanTypeConstruction),
			Value: string(LoanTypeConstruction),
		},
	})
	// register loan source
	enums.Register(LoanSourceName, []enums.EnumOpt{
		{
			Label: string(LoanSourceBorrower),
			Value: string(LoanSourceBorrower),
		},
		{
			Label: string(LoanSourceLender),
			Value: string(LoanSourceLender),
		},
		{
			Label: string(LoanSourceSystem),
			Value: string(LoanSourceSystem),
		},
	})
	// register lien position
	enums.Register(LienPositionName, []enums.EnumOpt{
		{
			Label: string(LienPosition1),
			Value: string(LienPosition1),
		},
		{
			Label: string(LienPosition2),
			Value: string(LienPosition2),
		},
		{
			Label: string(LienPosition3),
			Value: string(LienPosition3),
		},
		{
			Label: string(LienPosition4),
			Value: string(LienPosition4),
		},
		{
			Label: string(LienPositionOther),
			Value: string(LienPositionOther),
		},
	})
	// register business strategy
	enums.Register(BusinessStrategyName, []enums.EnumOpt{
		{
			Label: string(BusinessStrategySale),
			Value: string(BusinessStrategySale),
		},
		{
			Label: string(BusinessStrategyRental),
			Value: string(BusinessStrategyRental),
		},
		{
			Label: string(BusinessStrategyFlip),
			Value: string(BusinessStrategyFlip),
		},
		{
			Label: string(BusinessStrategyADUAddition),
			Value: string(BusinessStrategyADUAddition),
		},
		{
			Label: string(BusinessStrategySingleLotSFR),
			Value: string(BusinessStrategySingleLotSFR),
		},
		{
			Label: string(BusinessStrategySB9LotSplit),
			Value: string(BusinessStrategySB9LotSplit),
		},
		{
			Label: string(BusinessStrategyConformingLotSplit),
			Value: string(BusinessStrategyConformingLotSplit),
		},
		{
			Label: string(BusinessStrategyNewCondoTH),
			Value: string(BusinessStrategyNewCondoTH),
		},
		{
			Label: string(BusinessStrategyCommunityDevelopment),
			Value: string(BusinessStrategyCommunityDevelopment),
		},
		{
			Label: string(BusinessStrategyCommercial),
			Value: string(BusinessStrategyCommercial),
		},
	})

}

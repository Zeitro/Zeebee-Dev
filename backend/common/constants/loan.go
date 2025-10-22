package constants

import (
	"fmt"

	"zeebee-backend/common/enums"
)

// stage: Pending,Active,Archived
type LoanStage string

const LoanStageName = "loan_stage"
const (
	LoanStageDraft LoanStage = "Draft"
	// LoanStagePredraft LoanStage = "Predraft"
	LoanStagePending  LoanStage = "Pending"
	LoanStageActive   LoanStage = "Active"
	LoanStageArchived LoanStage = "Archived"
	LoanStageRejected LoanStage = "Rejected"
)

// pending_status: Quoting,Processing,Closing
type LoanPendingStatus string

const LoanPendingStatusName = "loan_pending_status"
const (
	LoanPendingStatusNone         LoanPendingStatus = "None"
	LoanPendingStatusQuoting      LoanPendingStatus = "Quoting"
	LoanPendingStatusProcessing   LoanPendingStatus = "Processing"
	LoanPendingStatusUnderWriting LoanPendingStatus = "Underwriting"
	LoanPendingStatusClosing      LoanPendingStatus = "Closing"
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
	LoanTypeRental       LoanType = "Rental"
)

// source: borrower,lender,system
type LoanSource string

const LoanSourceName = "loan_source"
const (
	LoanSourceBorrower   LoanSource = "borrower"
	LoanSourceLender     LoanSource = "lender"
	LoanSourceSystem     LoanSource = "system"
	LoanSourcePreApprove LoanSource = "preapprove"
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

type LoanPrepaymentPenalty int32

const LoanPrepaymentPenaltyName = "loan_prepayment_penalty"
const (
	LoanPrepaymentPenaltyNone   LoanPrepaymentPenalty = 0
	LoanPrepaymentPenaltyOne    LoanPrepaymentPenalty = 1
	LoanPrepaymentPenaltyTwo    LoanPrepaymentPenalty = 2
	LoanPrepaymentPenaltyThree  LoanPrepaymentPenalty = 3
	LoanPrepaymentPenaltyFour   LoanPrepaymentPenalty = 4
	LoanPrepaymentPenaltyFive   LoanPrepaymentPenalty = 5
	LoanPrepaymentPenaltySix    LoanPrepaymentPenalty = 6
	LoanPrepaymentPenaltySeven  LoanPrepaymentPenalty = 7
	LoanPrepaymentPenaltyEight  LoanPrepaymentPenalty = 8
	LoanPrepaymentPenaltyNine   LoanPrepaymentPenalty = 9
	LoanPrepaymentPenaltyTen    LoanPrepaymentPenalty = 10
	LoanPrepaymentPenaltyEleven LoanPrepaymentPenalty = 11
	LoanPrepaymentPenaltyTwelve LoanPrepaymentPenalty = 12
)

type LoanExperienceType string

const LoanExperienceTypeName = "loan_experience_type"
const (
	LoanExperienceTypeTransaction LoanExperienceType = "Transaction"
	LoanExperienceTypeRealEstate  LoanExperienceType = "Real Estate"
)

type LoanThirdServiceType string

const LoanThirdServiceTypeName = "loan_third_service_type"
const (
	LoanThirdServiceTypeTitleCo                 LoanThirdServiceType = "Title Co"
	LoanThirdServiceTypeEscrowCo                LoanThirdServiceType = "Escrow Co"
	LoanThirdServiceTypeInsurance               LoanThirdServiceType = "Insurance"
	LoanThirdServiceTypeAppraisalScheduling     LoanThirdServiceType = "Appraisal Scheduling"
	LoanThirdServiceTypeAppraisalAccessContact  LoanThirdServiceType = "Appraisal Access Contact"
	LoanThirdServiceTypeAppraisalPaymentContact LoanThirdServiceType = "Appraisal Payment Contact"
)

type BankruptcyChapter string

const BankruptcyChapterName = "bankruptcy_chapter"
const (
	BankruptcyChapter7  BankruptcyChapter = "Chapter 7"
	BankruptcyChapter11 BankruptcyChapter = "Chapter 11"
	BankruptcyChapter12 BankruptcyChapter = "Chapter 12"
	BankruptcyChapter13 BankruptcyChapter = "Chapter 13"
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
		{
			Label: string(LoanStageRejected),
			Value: string(LoanStageRejected),
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
			Label: string(LoanPendingStatusUnderWriting),
			Value: string(LoanPendingStatusUnderWriting),
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
		{
			Label: string(LoanTypeRental),
			Value: string(LoanTypeRental),
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
		{
			Label: string(LoanSourcePreApprove),
			Value: string(LoanSourcePreApprove),
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
	// register loan prepayment penalty
	enums.Register(LoanPrepaymentPenaltyName, []enums.EnumOpt{
		{
			Label: "None",
			Value: fmt.Sprintf("%d", LoanPrepaymentPenaltyNone),
		},
		{
			Label: fmt.Sprintf("%d month interest", LoanPrepaymentPenaltyOne),
			Value: fmt.Sprintf("%d", LoanPrepaymentPenaltyOne),
		},
		{
			Label: fmt.Sprintf("%d months interest", LoanPrepaymentPenaltyTwo),
			Value: fmt.Sprintf("%d", LoanPrepaymentPenaltyTwo),
		},
		{
			Label: fmt.Sprintf("%d months interest", LoanPrepaymentPenaltyThree),
			Value: fmt.Sprintf("%d", LoanPrepaymentPenaltyThree),
		},
		{
			Label: fmt.Sprintf("%d months interest", LoanPrepaymentPenaltyFour),
			Value: fmt.Sprintf("%d", LoanPrepaymentPenaltyFour),
		},
		{
			Label: fmt.Sprintf("%d months interest", LoanPrepaymentPenaltyFive),
			Value: fmt.Sprintf("%d", LoanPrepaymentPenaltyFive),
		},
		{
			Label: fmt.Sprintf("%d months interest", LoanPrepaymentPenaltySix),
			Value: fmt.Sprintf("%d", LoanPrepaymentPenaltySix),
		},
		{
			Label: fmt.Sprintf("%d months interest", LoanPrepaymentPenaltySeven),
			Value: fmt.Sprintf("%d", LoanPrepaymentPenaltySeven),
		},
		{
			Label: fmt.Sprintf("%d months interest", LoanPrepaymentPenaltyEight),
			Value: fmt.Sprintf("%d", LoanPrepaymentPenaltyEight),
		},
		{
			Label: fmt.Sprintf("%d months interest", LoanPrepaymentPenaltyNine),
			Value: fmt.Sprintf("%d", LoanPrepaymentPenaltyNine),
		},
		{
			Label: fmt.Sprintf("%d months interest", LoanPrepaymentPenaltyTen),
			Value: fmt.Sprintf("%d", LoanPrepaymentPenaltyTen),
		},
		{
			Label: fmt.Sprintf("%d months interest", LoanPrepaymentPenaltyEleven),
			Value: fmt.Sprintf("%d", LoanPrepaymentPenaltyEleven),
		},
		{
			Label: fmt.Sprintf("%d months interest", LoanPrepaymentPenaltyTwelve),
			Value: fmt.Sprintf("%d", LoanPrepaymentPenaltyTwelve),
		},
	})

	// register loan experience type
	enums.Register(LoanExperienceTypeName, []enums.EnumOpt{
		{
			Label: string(LoanExperienceTypeTransaction),
			Value: string(LoanExperienceTypeTransaction),
		},
		{
			Label: string(LoanExperienceTypeRealEstate),
			Value: string(LoanExperienceTypeRealEstate),
		},
	})

	// register loan third service type
	enums.Register(LoanThirdServiceTypeName, []enums.EnumOpt{
		{
			Label: string(LoanThirdServiceTypeTitleCo),
			Value: string(LoanThirdServiceTypeTitleCo),
		},
		{
			Label: string(LoanThirdServiceTypeEscrowCo),
			Value: string(LoanThirdServiceTypeEscrowCo),
		},
		{
			Label: string(LoanThirdServiceTypeInsurance),
			Value: string(LoanThirdServiceTypeInsurance),
		},
		{
			Label: string(LoanThirdServiceTypeAppraisalAccessContact),
			Value: string(LoanThirdServiceTypeAppraisalAccessContact),
		},
		{
			Label: string(LoanThirdServiceTypeAppraisalScheduling),
			Value: string(LoanThirdServiceTypeAppraisalScheduling),
		},
		{
			Label: string(LoanThirdServiceTypeAppraisalPaymentContact),
			Value: string(LoanThirdServiceTypeAppraisalPaymentContact),
		},
	})

	// register bankruptcy chapter
	enums.Register(BankruptcyChapterName, []enums.EnumOpt{
		{
			Label: string(BankruptcyChapter7),
			Value: string(BankruptcyChapter7),
		},
		{
			Label: string(BankruptcyChapter11),
			Value: string(BankruptcyChapter11),
		},
		{
			Label: string(BankruptcyChapter12),
			Value: string(BankruptcyChapter12),
		},
		{
			Label: string(BankruptcyChapter13),
			Value: string(BankruptcyChapter13),
		},
	})

}

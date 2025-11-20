package model

import (
	"time"

	"gorm.io/plugin/soft_delete"

	"zeebee-backend/common/constants"
)

type Base struct {
	Id        int64                 `gorm:"primaryKey;autoIncrement;<-:create" json:"id"`
	CreatedAt time.Time             `gorm:"column:created_at;default:now();<-:create" json:"created_at"`
	UpdatedAt time.Time             `gorm:"column:updated_at;default:now()" json:"updated_at"`
	IsDel     soft_delete.DeletedAt `gorm:"softDelete:flag" json:"-"`
}

type Address struct {
	State    string `json:"state" copier:"State"`
	City     string `json:"city"`
	Address1 string `json:"address1"`
	Address2 string `json:"address2"`
	ZipCode  string `json:"zip_code"`
}

// LoanExtendValues represents extended information related to a loan.
type LoanExtendValues struct {
	Pricing           string `json:"pricing"`
	LoanStatusUpdates string `json:"loan_status_updates"`
	MeetingNotes      string `json:"meeting_notes"`
	ActionList        string `json:"action_list"`
	BusinessLine      string `json:"business_line"`
}

// LoanDemographicProperty
type LoanDemographicProperty struct {
	AsPrimaryResident          bool    `json:"as_primary_resident"`
	HasTransactionRelationship bool    `json:"has_transaction_relationship"`
	HasBorrowedMoney           bool    `json:"has_borrowed_money"`
	BorrowedMoneyAmount        float64 `json:"borrowed_money_amount"`
	HasAnotherMortgageLoan     bool    `json:"has_another_mortgage_loan"`
	HasPriorityOver1st         bool    `json:"has_priority_over_1st"`
	HasOutstandingJudgments    bool    `json:"has_outstanding_judgments"`
	HasFederalDebt             bool    `json:"has_federal_debt"`
	HasALawsuit                bool    `json:"has_a_lawsuit"`
	HasConveyedTitle           bool    `json:"has_conveyed_title"`
	HasAPreForeclosureSale     bool    `json:"has_a_pre_foreclosure_sale"`
	HasForeclosedUpon          bool    `json:"has_foreclosed_upon"`
	HasABankruptcy             bool    `json:"has_a_bankruptcy"`
	BankruptcyChapter          string  `json:"bankruptcy_chapter"`
}

// LoanDemographicEthnicity
type LoanDemographicEthnicity struct {
	Values                    []string `json:"value"`
	OtherHispanicOrLatinoName string   `json:"other_hispanic_or_latino_name"`
}

// LoanDemographicRace
type LoanDemographicRace struct {
	Values                           []string `json:"value"`
	AsianValues                      []string `json:"asian_values"`
	AmericanIndianOrAlaskaNativeName string   `json:"american_indian_or_alaska_native_name"`
	AsianOtherName                   string   `json:"asian_other_name"`
	OtherPacificIslanderName         string   `json:"other_pacific_islander_name"`
}

// LoanDemographicExtraDemographic
type LoanDemographicExtraDemographic struct {
	IsEthnicityCollectedByObservation bool     `json:"is_ethnicity_collected_by_observation"`
	IsSexCollectedByObservation       bool     `json:"is_sex_collected_by_observation"`
	IsRaceCollectedByObservation      bool     `json:"is_race_collected_by_observation"`
	InfoChannels                      []string `json:"info_channels"`
}

type LoanOperationStatus struct {
	Stage         constants.LoanStage         `json:"stage"`
	PendingStatus constants.LoanPendingStatus `json:"pending_status"`
}
type OperationLoanDocument struct {
	Category constants.LoanDocumentCategory `json:"category"`
	FileType constants.LoanDocumentFileType `json:"file_type"`
	Status   constants.LoanDocumentStatus   `json:"status"`
	Files    []*OperationFile               `json:"files,omitempty"`
}

type OperationFile struct {
	EntityId   int64                    `json:"entity_id"`
	EntityType constants.FileEntityType `json:"entity_type"`
	FileUrl    string                   `json:"file_url"`
}

type LoanOperationMetadata struct {
	Comment       string                   `json:"comment,omitempty"`
	LoanDocuments []*OperationLoanDocument `json:"loan_documents,omitempty"`
	// Files         []*OperationFile `json:"files,omitempty"`
}

type LoanBudget struct {
	OriginalBudget        float64 `json:"original_budget"`
	RemainingBudget       float64 `json:"remaining_budget"`
	CostSpentToDateAmount float64 `json:"cost_spent_to_date_amount"`
}

package model

type LoanDemographic struct {
	Base
	LoanId    int64                    `gorm:"column:loan_id;not null;default:0" json:"loan_id"`
	Property  LoanDemographicProperty  `gorm:"column:property;type:jsonb;serializer:json" json:"property"`
	Ethnicity LoanDemographicEthnicity `gorm:"column:ethnicity;type:jsonb;serializer:json" json:"ethnicity"`
	Race      LoanDemographicRace      `gorm:"column:race;type:jsonb;serializer:json" json:"race"`
	Sex       string                   `gorm:"column:sex;size:10;not null;default:'Male'" json:"sex"`
	// Relationships
	Loan *Loan `json:"loan" gorm:"foreignKey:LoanId;references:Id"`
}

func (m LoanDemographic) TableName() string {
	return "loan_demographics"
}

func (m *LoanDemographic) RemoveRelationships() {
	m.Loan = nil
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

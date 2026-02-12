package bo

type LoanDemographicEntityBo struct {
	Id        int64                             `json:"id"`
	LoanId    int64                             `json:"loan_id"`
	Property  *LoanDemographicPropertyEntityBo  `json:"property"`
	Ethnicity *LoanDemographicEthnicityEntityBo `json:"ethnicity"`
	Race      *LoanDemographicRaceEntityBo      `json:"race"`
	Sex       string                            `json:"sex"`
}

type LoanDemographicPropertyEntityBo struct {
	AsPrimaryResident          bool    `json:"as_primary_resident"`
	HasTransactionRelationship bool    `json:"has_transaction_relationship"`
	HasBorrowedMoney           bool    `json:"has_borrowed_money"`
	BorrowedMoneyAmount        float64 `json:"borrowed_money_amount"`
	HasAnotherMortgageLoan     bool    `json:"has_another_mortgage_loan"`
	HasPriorityOver1st         bool    `json:"has_priority_over1st"`
	HasOutstandingJudgments    bool    `json:"has_outstanding_judgments"`
	HasFederalDebt             bool    `json:"has_federal_debt"`
	HasALawsuit                bool    `json:"has_a_lawsuit"`
	HasConveyedTitle           bool    `json:"has_conveyed_title"`
	HasAPreForeclosureSale     bool    `json:"has_a_pre_foreclosure_sale"`
	HasForeclosedUpon          bool    `json:"has_foreclosed_upon"`
	HasABankruptcy             bool    `json:"has_a_bankruptcy"`
	BankruptcyChapter          string  `json:"bankruptcy_chapter"`
}
type LoanDemographicEthnicityEntityBo struct {
	Values                    []string `json:"values"`
	OtherHispanicOrLatinoName string   `json:"other_hispanic_or_latino_name"`
}
type LoanDemographicRaceEntityBo struct {
	Values                           []string `json:"values"`
	AsianValues                      []string `json:"asian_values"`
	AmericanIndianOrAlaskaNativeName string   `json:"american_indian_or_alaska_native_name"`
	AsianOtherName                   string   `json:"asian_other_name"`
	OtherPacificIslanderName         string   `json:"other_pacific_islander_name"`
}

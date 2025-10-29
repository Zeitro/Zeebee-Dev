package types

// LoanDemographicProperty 房产信息条款
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

	// BankruptcyChapter7         bool    `json:"bankruptcy_chapter7"`
	// BankruptcyChapter11        bool    `json:"bankruptcy_chapter11"`
	// BankruptcyChapter12        bool    `json:"bankruptcy_chapter12"`
	// BankruptcyChapter13        bool    `json:"bankruptcy_chapter13"`
}

// LoanDemographicEthnicity 种族条款
type LoanDemographicEthnicity struct {
	// IsHispanicOrLatino        bool   `json:"is_hispanic_or_latino"`
	// IsOtherHispanicOrLatino   bool   `json:"is_other_hispanic_or_latino"`
	// OtherHispanicOrLatinoName string `json:"other_hispanic_or_latino_name"`
	// IsMexican                 bool   `json:"is_mexican"`
	// IsCuban                   bool   `json:"is_cuban"`
	// IsNotHispanicOrLatino     bool   `json:"is_not_hispanic_or_latino"`
	// Hidden                    bool   `json:"hidden"`

	Values                    []string `json:"value"`
	OtherHispanicOrLatinoName string   `json:"other_hispanic_or_latino_name"`
}

// LoanDemographicRace 种族条款
type LoanDemographicRace struct {
	// IsAmericanIndianOrAlaskaNative            bool   `json:"is_american_indian_or_alaska_native"`
	// AmericanIndianOrAlaskaNativeName          string `json:"american_indian_or_alaska_native_name"`
	// IsAsian                                   bool   `json:"is_asian"`
	// IsAsianIndian                             bool   `json:"is_asian_indian"`
	// IsAsianChinese                            bool   `json:"is_asian_chinese"`
	// IsAsianFilipino                           bool   `json:"is_asian_filipino"`
	// IsAsianJapanese                           bool   `json:"is_asian_japanese"`
	// IsAsianKorean                             bool   `json:"is_asian_korean"`
	// IsAsianVietnamese                         bool   `json:"is_asian_vietnamese"`
	// IsAsianOther                              bool   `json:"is_asian_other"`
	// AsianOtherName                            string `json:"asian_other_name"`
	// IsBlackOrAfricanAmerican                  bool   `json:"is_black_or_african_american"`
	// IsNativeHawaiianOrOtherPacificIslander    bool   `json:"is_native_hawaiian_or_other_pacific_islander"`
	// IsNativeHawaiianGuamanianOrChamorroSamoan bool   `json:"is_native_hawaiian_guamanian_or_chamorro_samoan"`
	// IsOtherPacificIslander                    bool   `json:"is_other_pacific_islander"`
	// OtherPacificIslanderName                  string `json:"other_pacific_islander_name"`
	// IsWhite                                   bool   `json:"is_white"`
	// Hidden                                    bool   `json:"hidden"`

	Values                           []string `json:"value"`
	AsianValues                      []string `json:"asian_values"`
	AmericanIndianOrAlaskaNativeName string   `json:"american_indian_or_alaska_native_name"`
	AsianOtherName                   string   `json:"asian_other_name"`
	OtherPacificIslanderName         string   `json:"other_pacific_islander_name"`
}

// LoanDemographicExtraDemographic 获取种族信息的方式
type LoanDemographicExtraDemographic struct {
	IsEthnicityCollectedByObservation bool `json:"is_ethnicity_collected_by_observation"`
	IsSexCollectedByObservation       bool `json:"is_sex_collected_by_observation"`
	IsRaceCollectedByObservation      bool `json:"is_race_collected_by_observation"`
	// IsByFaceToFace                    bool `json:"is_by_face_to_face"`
	// IsByTelephone                     bool `json:"is_by_telephone"`
	// IsByFoxOrMail                     bool `json:"is_by_fox_or_mail"`
	// IsByEmailOrInternet               bool `json:"is_by_email_or_internet"`
	InfoChannels []string `json:"info_channels"`
}

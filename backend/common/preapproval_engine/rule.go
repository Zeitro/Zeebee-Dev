package preapproval_engine

// RuleCap defines the loan cap rules including LTAIV, LTC and LTARV
type RuleCap struct {
	Ltaiv *Values `json:"ltaiv,omitempty"`
	Ltc   *Values `json:"ltc,omitempty"`
	Ltarv *Values `json:"ltarv,omitempty"`
}

// Values represents standard and large unit values
type Values struct {
	Standard  float64 `json:"standard"`
	LargeUnit float64 `json:"large_unit"`
}

// Rules interface for general loan type rules
type Rules interface {
	GetGeneralRules() *RuleCap
	GetHighValueRules() *RuleCap
}

// RehabRules interface for rehab loan type rules
type RehabRules interface {
	Rules
	GetHeavyRules() (general, highValue *RuleCap)
	GetLightRules() (general, highValue *RuleCap)
}

// GeneralLoanTypeRules for general loan types
type GeneralLoanTypeRules struct {
	GeneralRules   *RuleCap
	HighValueRules *RuleCap
}

// GetGeneralRules returns general rules
func (r *GeneralLoanTypeRules) GetGeneralRules() *RuleCap {
	return r.GeneralRules
}

// GetHighValueRules returns high value rules
func (r *GeneralLoanTypeRules) GetHighValueRules() *RuleCap {
	return r.HighValueRules
}

// RehabLoanTypeRules for rehab loan types
type RehabLoanTypeRules struct {
	Heavy RehabLoanTypeRule
	Light RehabLoanTypeRule
}

type RehabLoanTypeRule struct {
	GeneralRule   *RuleCap
	HighValueRule *RuleCap
}

// GetHeavyRules returns heavy rehab rules
func (r *RehabLoanTypeRules) GetHeavyRules() (general, highValue *RuleCap) {
	return r.Heavy.GeneralRule, r.Heavy.HighValueRule
}

// GetLightRules returns light rehab rules
func (r *RehabLoanTypeRules) GetLightRules() (general, highValue *RuleCap) {
	return r.Light.GeneralRule, r.Light.HighValueRule
}

// GetGeneralRules returns general rules (default to heavy rules)
func (r *RehabLoanTypeRules) GetGeneralRules() *RuleCap {
	return r.Heavy.GeneralRule
}

// GetHighValueRules returns high value rules (default to heavy rules)
func (r *RehabLoanTypeRules) GetHighValueRules() *RuleCap {
	return r.Heavy.HighValueRule
}

// ExperienceLevelRule defines rules for determining borrower experience level
type ExperienceLevelRule struct {
	MaxTransactionHistoryCountForNoExperience int `json:"max_transaction_history_count_for_no_experience"`
	MaxTransactionHistoryCountForExperienced  int `json:"max_transaction_history_count_for_experienced"`
}

// getExperienceLevel determines the experience level based on transaction history count
func (r *ExperienceLevelRule) getExperienceLevel(transactionHistoryCount int) ExperienceLevel {
	if transactionHistoryCount <= r.MaxTransactionHistoryCountForNoExperience {
		return ExperienceLevelNoExperienced
	} else if transactionHistoryCount > r.MaxTransactionHistoryCountForNoExperience &&
		transactionHistoryCount <= r.MaxTransactionHistoryCountForExperienced {
		return ExperienceLevelExperienced
	} else {
		return ExperienceLevelInstitutional
	}
}

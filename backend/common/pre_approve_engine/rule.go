package pre_approve_engine

type RuleCap struct {
	LTAIV *Values
	LTC   *Values
	LTARV *Values
}

type Values struct {
	Standard  float32
	LargeUnit float32
}

type Rules interface {
	GetGeneralRules() *RuleCap
	GetHighValueRules() *RuleCap
}

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

// implements Rules
func (r *GeneralLoanTypeRules) GetGeneralRules() *RuleCap {
	return r.GeneralRules
}

func (r *GeneralLoanTypeRules) GetHighValueRules() *RuleCap {
	return r.HighValueRules
}

// RehabLoanTypeRules for rehab loan types
type RehabLoanTypeRules struct {
	Heavy struct {
		GeneralRules   *RuleCap
		HighValueRules *RuleCap
	}
	Light struct {
		GeneralRules   *RuleCap
		HighValueRules *RuleCap
	}
}

func (r *RehabLoanTypeRules) GetHeavyRules() (general, highValue *RuleCap) {
	return r.Heavy.GeneralRules, r.Heavy.HighValueRules
}

func (r *RehabLoanTypeRules) GetLightRules() (general, highValue *RuleCap) {
	return r.Light.GeneralRules, r.Light.HighValueRules
}

func (r *RehabLoanTypeRules) GetGeneralRules() *RuleCap {
	return r.Heavy.GeneralRules
}

func (r *RehabLoanTypeRules) GetHighValueRules() *RuleCap {
	return r.Heavy.HighValueRules
}

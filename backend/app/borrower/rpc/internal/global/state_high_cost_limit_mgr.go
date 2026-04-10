package global

import "zeebee-backend/common/conforming_loan_limit"

var StateHighCostLimitMgr *conforming_loan_limit.StateHighCostLimitMgr

func initStateHighCostLimitMgr(content []byte, limitAmount int64) (err error) {
	StateHighCostLimitMgr, err = conforming_loan_limit.ParseHighCostLimit(content, limitAmount)
	return
}

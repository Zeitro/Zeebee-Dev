package constants

const (
	BorrowerLoginFailedTimes   = "zeebee:borrower:%v:login:failed:times"
	BorrowerRegisterVerifyCode = "zeebee:borrower:%v:register:verify:code"

	BorrowerRegisterVerifyEmail = "zeebee:borrower:%v:register:verify:email"
	BorrowerLoginEmail          = "zeebee:borrower:%v:login:email"
	BorrowerResetPwdEmail       = "zeebee:borrower:%v:reset:pwd:email"
	BorrowerResetPwdCode        = "zeebee:borrower:%v:reset:pwd:code"
	BorrowerPreApproveCacheKey  = "zeebee:borrower:%v:preapprove:loan"

	LenderLoginFailedTimes   = "zeebee:lender:%v:login:failed:times"
	LenderRegisterVerifyCode = "zeebee:lender:%v:register:verify:code"

	LenderRegisterVerifyEmail = "zeebee:lender:%v:register:verify:email"
	LenderLoginEmail          = "zeebee:lender:%v:login:email"
	LenderResetPwdEmail       = "zeebee:lender:%v:reset:pwd:email"
	LenderResetPwdCode        = "zeebee:lender:%v:reset:pwd:code"
)

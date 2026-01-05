package limiter

type LimiterStrategy string

const (
	StrategyLogin              LimiterStrategy = "Login"
	StrategySendLoginEmail     LimiterStrategy = "SendLoginEmail"
	StrategySendVerifyEmail    LimiterStrategy = "SendVerifyEmail"
	StrategySendResetPwdEmail  LimiterStrategy = "SendResetPwdEmail"
	StrategySendLoanQuoteEmail LimiterStrategy = "SendLoanQuoteEmail"
)

package constants

const (
	UserLoginFailedTimes   = "zeebee:%v:%v:login:failed:times"
	UserRegisterVerifyCode = "zeebee:%v:%v:register:verify:code"
	UserInviteCode         = "zeebee:%v:%v:%v:invite:code"
	UserInviteCodeMarker   = "zeebee:%v:%v:%v:invite:code:marker"
	UserInviteUserPwd      = "zeebee:%v:%v:%v:invite:pwd"
	UserInviteEmailSendTag = "zeebee:%v:%v:%v:invite:email:send:tag"

	BorrowerPreApprovalCacheKey            = "zeebee:borrower:%v:preapprovalLoan"
	BorrowerSendAprApprovalNotificationKey = "zeebee:organizationId:%v:requestId:%v:preapprovalLoan"

	RegisterVerifyEmail = "zeebee:%v:%v:register:verify:email"
	LoginEmail          = "zeebee:%v:%v:login:email"
	ResetPwdEmail       = "zeebee:%v:%v:reset:pwd:email"
	ResetPwdCode        = "zeebee:%v:%v:reset:pwd:code"
)

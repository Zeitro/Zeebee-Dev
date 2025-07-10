package constants

// 系统常量
const (
	// 服务名称
	ServiceNameWebserver  = "webserver"
	ServiceNameStreamline = "streamline"
	ServiceNameForms      = "forms"

	// 数据库相关
	DatabaseName = "zeebee"

	// JWT相关
	JWTSecretKey   = "zeebee-jwt-secret-key"
	JWTExpireHours = 24

	// 状态码
	StatusSuccess             = 200
	StatusCreated             = 201
	StatusBadRequest          = 400
	StatusUnauthorized        = 401
	StatusForbidden           = 403
	StatusNotFound            = 404
	StatusInternalServerError = 500
)

// 用户角色 TBD
const (
	RoleAdmin    = "admin"
	RoleLender   = "lender"
	RoleBorrower = "borrower"
	RoleAgent    = "agent"
)

// 贷款状态 TBD
const (
	LoanStatusDraft     = "draft"
	LoanStatusPending   = "pending"
	LoanStatusApproved  = "approved"
	LoanStatusRejected  = "rejected"
	LoanStatusActive    = "active"
	LoanStatusCompleted = "completed"
	LoanStatusDefaulted = "defaulted"
)

package logic

import (
	"context"
	"sync"

	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"zeebee-backend/app/admin/rpc/internal/global"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/rbac"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data"
	"zeebee-backend/domains/zeebee/logic"
	"zeebee-backend/domains/zeebee/model"
)

var ErrorInvalidRequest = errorx.New("Invalid request")

var baseLogic *BaseLogic
var onceBase sync.Once

type BaseLogic struct {
	// ctx    context.Context
	svcCtx *svc.ServiceContext
	// logx.Logger
	*logic.AllLogics
}

func GetBaseLogic(svcCtx *svc.ServiceContext) *BaseLogic {
	onceBase.Do(func() {
		baseLogic = NewBaseLogic(svcCtx)
	})
	return baseLogic
}

func NewBaseLogic(svcCtx *svc.ServiceContext) *BaseLogic {
	redisCli := global.GetRedis()
	jwtConfig := logic.JwtConfig{
		AccessSecret:  svcCtx.Config.JwtAuth.AccessSecret,
		AccessExpire:  svcCtx.Config.JwtAuth.AccessExpire,
		RefreshSecret: svcCtx.Config.JwtAuth.RefreshSecret,
		RefreshExpire: svcCtx.Config.JwtAuth.RefreshExpire,
	}
	emailClient := global.EmailClient
	enforcerMap := map[constants.UserType]*rbac.Enforcer{
		constants.UserTypeAdmin:    global.AdminEnforcer,
		constants.UserTypeLender:   global.LenderEnforcer,
		constants.UserTypeBorrower: global.BorrowerEnforcer,
	}

	currentEnforcer := global.AdminEnforcer
	_ = currentEnforcer.LoadFilteredPolicy(gormadapter.Filter{
		Ptype: []string{"p"},
		V0:    nil,
		V1:    []string{string(constants.UserTypeAdmin)},
		V2:    nil,
		V3:    nil,
		V4:    nil,
		V5:    nil,
	})

	return &BaseLogic{
		// ctx:       ctx,
		svcCtx: svcCtx,
		// Logger:    logx.WithContext(ctx),
		AllLogics: logic.NewAllLogics(svcCtx.AllRepos, redisCli, jwtConfig, emailClient, currentEnforcer, enforcerMap),
	}
}

func (l *BaseLogic) CreateLoanOperationLog(ctx context.Context, tx *gorm.DB, loanId, operatorId int64, from, to data.LoanOperationStatus, action constants.LoanOperationAction, metadata *data.LoanOperationMetadata) (err error) {
	defer func() {
		if err != nil {
			logx.WithContext(ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
				logx.Field("loan_id", loanId),
				logx.Field("operator_id", operatorId),
				logx.Field("from", from),
				logx.Field("to", to),
				logx.Field("action", action),
				logx.Field("metadata", metadata),
				logx.Field("OperatorType", constants.LoanOperationOperatorTypeAdmin),
			).Error("failed to create loan operation log")
		}
	}()
	operator, _ := l.svcCtx.OrganizationUserRepo.FindOneByIdWithTx(ctx, tx, operatorId)

	operatorName := ""
	if operator != nil {
		operatorName = operator.FullName
	}

	operation := &model.LoanOperation{
		LoanId:       loanId,
		OperatorId:   operatorId,
		OperatorName: operatorName,
		OperatorType: constants.LoanOperationOperatorTypeAdmin,
		Action:       action,
		From:         from,
		To:           to,
		Metadata:     metadata,
	}
	_, err = l.svcCtx.LoanOperationRepo.CreateWithTx(ctx, tx, operation, false)
	if err != nil {
		return errorx.Wrap(err, "Failed to create operation log")
	}

	return nil
}

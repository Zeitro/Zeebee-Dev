package logic

import (
	"context"
	"sync"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
	"gorm.io/gorm"

	"zeebee-backend/app/admin/rpc/internal/global"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/logic"
	"zeebee-backend/domains/zeebee/model"
)

var ErrorInvalidRequest = errorx.New("invalid request")

var baseLogic *BaseLogic
var onceBase sync.Once

type BaseLogic struct {
	// ctx    context.Context
	svcCtx *svc.ServiceContext
	//
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
	jwtConfig := dto.JwtConfig{
		AccessSecret:  svcCtx.Config.JwtConf.AccessSecret,
		AccessExpire:  svcCtx.Config.JwtConf.AccessExpire,
		RefreshSecret: svcCtx.Config.JwtConf.RefreshSecret,
		RefreshExpire: svcCtx.Config.JwtConf.RefreshExpire,
	}
	emailClient := global.EmailClient

	currentEnforcer := global.AdminEnforcer

	currentAuthEventTopic := svcCtx.Config.CasbinConf.EventTopic.AdminTopic
	authEventTopicMap := map[constants.UserType]string{
		constants.UserTypeAdmin:    svcCtx.Config.CasbinConf.EventTopic.AdminTopic,
		constants.UserTypeLender:   svcCtx.Config.CasbinConf.EventTopic.LenderTopic,
		constants.UserTypeBorrower: svcCtx.Config.CasbinConf.EventTopic.BorrowerTopic,
	}
	logicConfig := logic.AllLogicConfig{
		Repos:                 svcCtx.AllRepos,
		RedisCli:              redisCli,
		JwtConfig:             jwtConfig,
		EmailClient:           emailClient,
		CurrentEnforcer:       currentEnforcer,
		CurrentAuthEventTopic: currentAuthEventTopic,
		AuthEventTopicMap:     authEventTopicMap,
		Logger:                svcCtx.Logger,
	}
	l := &BaseLogic{
		svcCtx:    svcCtx,
		AllLogics: logic.NewAllLogics(logicConfig),
	}

	threading.GoSafe(func() {
		ctx := context.Background()
		err := l.AuthLogic.RebuildPolicy(ctx, constants.UserTypeAdmin)
		if err != nil {
			l.svcCtx.Logger.WithContext(ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
				logx.Field("user_type", constants.UserTypeAdmin),
			).Error("fail to rebuild admin policy")
		}
	})

	return l
}

func (l *BaseLogic) CreateLoanOperationLog(ctx context.Context, tx *gorm.DB, loanId, operatorId int64, from, to dto.LoanOperationStatus, action constants.LoanOperationAction, metadata *dto.LoanOperationMetadata) (err error) {
	defer func() {
		if err != nil {
			l.svcCtx.Logger.WithContext(ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
				logx.Field("loan_id", loanId),
				logx.Field("operator_id", operatorId),
				logx.Field("from", from),
				logx.Field("to", to),
				logx.Field("action", action),
				logx.Field("metadata", metadata),
				logx.Field("OperatorType", constants.LoanOperationOperatorTypeAdmin),
			).Error("fail to create loan operation log")
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
		return errorx.Wrap(err, "fail to create operation log")
	}

	return nil
}

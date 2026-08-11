package logic

import (
	"context"
	"sync"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	"zeebee-backend/app/lender/rpc/internal/global"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/logic"
	"zeebee-backend/domains/zeebee/model"
)

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
	authEventTopicMap := map[constants.UserType]string{
		constants.UserTypeLender: svcCtx.Config.CasbinConf.EventTopic.LenderTopic,
	}

	logicConfig := logic.AllLogicConfig{
		Repos:       svcCtx.AllRepos,
		RedisCli:    redisCli,
		JwtConfig:   jwtConfig,
		EmailClient: global.EmailClient,
		// US-12: tenant is the sole casbin executor; lender no longer runs a
		// local enforcer, so it subscribes to no reload channel (AuthLogic is
		// never invoked here anymore).
		CurrentEnforcer:   nil,
		AuthEventTopicMap: authEventTopicMap,
		Logger:            svcCtx.Logger,
	}

	l := &BaseLogic{
		svcCtx:    svcCtx,
		AllLogics: logic.NewAllLogics(logicConfig),
	}

	return l
}

func (l *BaseLogic) GetUserClaimFromCtx(ctx context.Context) (utils.JwtUser, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok || len(md.Get(utils.JwtTokenCtxKey)) == 0 || len(md.Get(utils.UserTypeCtxKey)) == 0 {
		return utils.JwtUser{}, status.Error(codes.Unauthenticated, "not login")
	}
	claim, err := l.svcCtx.TenantClient.User.GetClaim(ctx, &tenant.GetClaimRequest{Data: &tenant.GetClaimRequestData{
		AccessToken: md.Get(utils.JwtTokenCtxKey)[0],
		UserType:    md.Get(utils.UserTypeCtxKey)[0],
	}})
	if err != nil {
		return utils.JwtUser{}, err
	}
	if claim.Data == nil {
		return utils.JwtUser{}, status.Error(codes.Unauthenticated, "invalid token")
	}
	return utils.JwtUser{
		UserId: claim.Data.UserId, OrganizationId: claim.Data.OrganizationId,
		OrganizationUserId: claim.Data.OrganizationUserId, Username: claim.Data.Username,
		Role: claim.Data.Role, UserType: constants.UserType(claim.Data.UserType),
	}, nil
}

// get lender loan
func (l *BaseLogic) GetLenderLoan(ctx context.Context, loanId int64, lenderOrganizationId int64, preloads []string) (loan *model.Loan, err error) {
	loan, err = l.svcCtx.LoanRepo.FindOneWithPreloadById(ctx, loanId, preloads)
	if err != nil {
		return nil, err
	}
	if loan.LenderOrganizationId != lenderOrganizationId {
		return nil, errorx.WithStack(errorx.New("invalid lender company"))
	}

	if loan.Source != constants.LoanSourceSystem && loan.Source != constants.LoanSourcePreapproval {
		return nil, errorx.WithStack(errorx.New("invalid loan source"))
	}

	return
}

// CreateLoanOperationLog creates a loan operation log record
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
				logx.Field("OperatorType", constants.LoanOperationOperatorTypeLender),
			).Error("fail to create loan operation log")
		}
	}()
	operatorName := ""
	if claim, claimErr := l.GetUserClaimFromCtx(ctx); claimErr == nil && claim.OrganizationUserId == operatorId {
		operatorName = claim.Username
	}

	operation := &model.LoanOperation{
		LoanId:       loanId,
		OperatorId:   operatorId,
		OperatorName: operatorName,
		OperatorType: constants.LoanOperationOperatorTypeLender,
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

func (l *BaseLogic) GetLoanAndDocumentByDocumentId(ctx context.Context, documentId int64, organization int64) (loan *model.Loan, loanDocument *model.LoanDocument, err error) {
	loanDocument, err = l.svcCtx.LoanDocumentRepo.FindOneWithPreloadById(ctx, documentId, []string{"Files"})
	if err != nil {
		return nil, nil, errorx.Wrap(err, "fail to find loan document")
	}

	loan, err = l.GetLenderLoan(ctx, loanDocument.LoanId, organization, []string{"LoanBorrower"})
	if err != nil {
		return nil, nil, err
	}

	return loan, loanDocument, nil
}

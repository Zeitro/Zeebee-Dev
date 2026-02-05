package logic

import (
	"context"

	"github.com/duke-git/lancet/v2/condition"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"zeebee-backend/app/lender/rpc/internal/global"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/logic"
	"zeebee-backend/domains/zeebee/model"
)

type BaseLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	*logic.AllLogics
}

func NewBaseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BaseLogic {
	redisCli := global.GetRedis()
	jwtConfig := logic.JwtConfig{
		AccessSecret:  svcCtx.Config.JwtAuth.AccessSecret,
		AccessExpire:  svcCtx.Config.JwtAuth.AccessExpire,
		RefreshSecret: svcCtx.Config.JwtAuth.RefreshSecret,
		RefreshExpire: svcCtx.Config.JwtAuth.RefreshExpire,
	}
	return &BaseLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		AllLogics: logic.NewAllLogics(svcCtx.AllRepos, redisCli, jwtConfig, global.EmailClient),
	}
}

// get lender loan
func (l *BaseLogic) GetLenderLoan(loanId int64, lenderOrganizationId int64, preloads []string) (loan *model.Loan, err error) {
	loan, err = l.svcCtx.LoanRepo.FindOneWithPreloadById(l.ctx, loanId, preloads)
	if err != nil {
		return nil, err
	}
	if loan.LenderOrganizationId != lenderOrganizationId {
		return nil, errorx.WithStack(errorx.New("Invalid lender company"))
	}

	if loan.Source != constants.LoanSourceSystem && loan.Source != constants.LoanSourcePreapproval {
		return nil, errorx.WithStack(errorx.New("Invalid loan source"))
	}

	return
}

// CreateLoanOperationLog creates a loan operation log record
func (l *BaseLogic) CreateLoanOperationLog(tx *gorm.DB, loanId, operatorId int64, from, to model.LoanOperationStatus, action constants.LoanOperationAction, metadata *model.LoanOperationMetadata) error {
	operator, err := l.svcCtx.UserRepo.FindOneByIdWithTx(l.ctx, tx, operatorId)
	if err != nil {
		return errorx.Wrap(err, "Failed To Get Operator")
	}

	operatorName := ""
	if operator != nil {
		operatorName = operator.FullName
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
	_, err = l.svcCtx.LoanOperationRepo.CreateWithTx(l.ctx, tx, operation, false)
	if err != nil {
		return errorx.Wrap(err, "Failed To Create Operation Log")
	}

	return nil
}

func (l *BaseLogic) GetLoanAndDocumentByDocumentId(documentId int64, organization int64) (loan *model.Loan, loanDocument *model.LoanDocument, err error) {
	loanDocument, err = l.svcCtx.LoanDocumentRepo.FindOneWithPreloadById(l.ctx, documentId, []string{"Files"})
	if err != nil {
		return nil, nil, errorx.Wrap(err, "Failed To Find Loan Document")
	}

	loan, err = l.GetLenderLoan(loanDocument.LoanId, organization, []string{"LoanBorrower"})
	if err != nil {
		return nil, nil, err
	}

	return loan, loanDocument, nil
}

// GetActiveOrganizationByUser returns the active organization for the given user
func (l *BaseLogic) GetActiveOrganizationByUser(user *model.User) (currentOrganization *model.Organization, isChanged bool, err error) {
	if user == nil {
		return nil, false, errorx.New("User is empty")
	}
	organizations := user.Organizations
	// find current company
	currentOrganization, ok := slice.FindBy(organizations, func(index int, item *model.Organization) bool {
		return item.Id == user.DefaultOrganizationId && item.Status == constants.OrganizationStatusActive
	})
	// find first active company
	if !ok {
		currentOrganization, ok = slice.FindBy(organizations, func(index int, item *model.Organization) bool {
			return item.Status == constants.OrganizationStatusActive
		})
	}
	if !ok {
		// no active company
		isChanged = condition.Ternary(user.DefaultOrganizationId == 0, false, true)
		user.DefaultOrganizationId = 0
	} else if user.DefaultOrganizationId != currentOrganization.Id {
		// company changed
		isChanged = true
		user.DefaultOrganizationId = currentOrganization.Id
	}

	return currentOrganization, isChanged, err
}

// CheckAndUpdateActiveOrganizationByUser checks and updates the active organization for the given user
func (l *BaseLogic) CheckAndUpdateActiveOrganizationByUser(user *model.User) (currentUser *model.User, err error) {
	currentOrganization, isChanged, err := l.GetActiveOrganizationByUser(user)
	if err != nil {
		return currentUser, err
	}
	if isChanged {
		var defaultOrganizationId int64
		if currentOrganization != nil {
			defaultOrganizationId = currentOrganization.Id
		}

		user.DefaultOrganizationId = defaultOrganizationId
		err = l.svcCtx.UserRepo.UpdateColumnsById(l.ctx, user.Id, map[string]interface{}{
			"default_organization_id": user.DefaultOrganizationId,
		})
		if err != nil {
			return currentUser, errorx.WithStack(err)
		}
	}
	return user, err
}

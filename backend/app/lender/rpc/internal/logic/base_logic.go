package logic

import (
	"context"

	"github.com/duke-git/lancet/v2/condition"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/common/constants"
	"zeebee-backend/domains/zeebee/model"
)

type BaseLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBaseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BaseLogic {
	return &BaseLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// get lender loan
func (l *BaseLogic) GetLenderLoan(loanId int64, lenderCompanyId int64, preloads []string) (loan *model.Loan, err error) {
	loan, err = l.svcCtx.LoanRepo.FindOneWithPreloadById(l.ctx, loanId, preloads)
	if err != nil {
		return nil, err
	}
	if loan.LenderCompanyId != lenderCompanyId {
		return nil, errors.WithStack(errors.New("Invalid Lender Company"))
	}

	if loan.Source != constants.LoanSourceSystem && loan.Source != constants.LoanSourcePreApproval {
		return nil, errors.WithStack(errors.New("Invalid Loan Source"))
	}

	return
}

// CreateLoanOperationLog creates a loan operation log record
func (l *BaseLogic) CreateLoanOperationLog(tx *gorm.DB, loanId, operatorId int64, from, to model.LoanOperationStatus, action constants.LoanOperationAction, metadata *model.LoanOperationMetadata) error {
	operator, err := l.svcCtx.LenderEmployeeRepo.FindOneByIdWithTx(l.ctx, tx, operatorId)
	if err != nil {
		return errors.Wrap(err, "Failed To Get Operator")
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
		return errors.Wrap(err, "Failed To Create Operation Log")
	}

	return nil
}

func (l *BaseLogic) GetLoanAndDocumentByDocumentId(documentId int64, lenderCompanyId int64) (loan *model.Loan, loanDocument *model.LoanDocument, err error) {
	loanDocument, err = l.svcCtx.LoanDocumentRepo.FindOneWithPreloadById(l.ctx, documentId, []string{"Files"})
	if err != nil {
		return nil, nil, errors.Wrap(err, "Failed To Find Loan Document")
	}

	loan, err = l.GetLenderLoan(loanDocument.LoanId, lenderCompanyId, []string{"LoanBorrower"})
	if err != nil {
		return nil, nil, err
	}

	return loan, loanDocument, nil
}

// GetActiveLenderCompanyByEmployee returns the active lender company for the given employee
func (l *BaseLogic) GetActiveLenderCompanyByEmployee(employee *model.LenderEmployee) (currentCompany *model.LenderCompany, isChanged bool, err error) {
	if employee == nil {
		return nil, false, errors.New("Employee Is Empty")
	}
	companies := employee.LenderCompanies
	// find current company
	currentCompany, ok := slice.FindBy(employee.LenderCompanies, func(index int, item *model.LenderCompany) bool {
		return item.Id == employee.DefaultLenderCompanyId && item.Status == constants.LenderCompanyStatusActive
	})
	// find first active company
	if !ok {
		currentCompany, ok = slice.FindBy(companies, func(index int, item *model.LenderCompany) bool {
			return item.Status == constants.LenderCompanyStatusActive
		})
	}
	if !ok {
		// no active company
		isChanged = condition.Ternary(employee.DefaultLenderCompanyId == 0, false, true)
		employee.DefaultLenderCompanyId = 0
	} else if employee.DefaultLenderCompanyId != currentCompany.Id {
		// company changed
		isChanged = true
		employee.DefaultLenderCompanyId = currentCompany.Id
	}

	return currentCompany, isChanged, err
}

// CheckAndUpdateActiveLenderCompanyByEmployee checks and updates the active lender company for the given employee
func (l *BaseLogic) CheckAndUpdateActiveLenderCompanyByEmployee(employee *model.LenderEmployee) (currentEmployee *model.LenderEmployee, err error) {
	currentCompany, isChanged, err := l.GetActiveLenderCompanyByEmployee(employee)
	if err != nil {
		return currentEmployee, err
	}
	if isChanged {
		employee.DefaultLenderCompanyId = condition.Ternary(currentCompany == nil, 0, currentCompany.Id)
		err = l.svcCtx.LenderEmployeeRepo.UpdateColumnsById(l.ctx, employee.Id, map[string]interface{}{
			"default_lender_company_id": employee.DefaultLenderCompanyId,
		})
		if err != nil {
			return currentEmployee, errors.WithStack(err)
		}
	}
	return employee, err
}

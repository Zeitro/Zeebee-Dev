package logic

import (
	"context"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"zeebee-backend/app/lender/rpc/internal/repo/model"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/common/constants"
	"zeebee-backend/common/types"
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
		return nil, errors.WithStack(errors.New("invalid lender company"))
	}

	if loan.Source != constants.LoanSourceSystem && loan.Source != constants.LoanSourcePreApprove {
		return nil, errors.WithStack(errors.New("invalid loan source"))
	}

	return
}

// CreateLoanOperationLog creates a loan operation log record
func (l *BaseLogic) CreateLoanOperationLog(tx *gorm.DB, loanId, operatorId int64, from, to types.LoanOperationStatus, action constants.LoanOperationAction, metadata *types.LoanOperationMetadata) error {
	operator, err := l.svcCtx.LenderEmployeeRepo.FindOneByIdWithTx(l.ctx, tx, operatorId)
	if err != nil {
		return errors.Wrap(err, "failed to get operator")
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
		return errors.Wrap(err, "failed to create operation log")
	}

	return nil
}

func (l *BaseLogic) GetLoanAndDocumentByDocumentId(documentId int64, lenderCompanyId int64) (loan *model.Loan, loanDocument *model.LoanDocument, err error) {
	loanDocument, err = l.svcCtx.LoanDocumentRepo.FindOneWithPreloadById(l.ctx, documentId, []string{"Files"})
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to find loan document")
	}

	loan, err = l.GetLenderLoan(loanDocument.LoanId, lenderCompanyId, []string{"LoanBorrower"})
	if err != nil {
		return nil, nil, err
	}

	return loan, loanDocument, nil
}

package loanlogic

import (
	"context"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/common/constants"
	"zeebee-backend/domains/zeebee/model"
)

type LoanBaseLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoanBaseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoanBaseLogic {
	return &LoanBaseLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// get lender loan
func (l *LoanBaseLogic) GetLenderLoan(loanId int64, lenderCompanyId int64, preloads []string) (loan *model.Loan, err error) {
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
func (l *LoanBaseLogic) CreateLoanOperationLog(tx *gorm.DB, loanId, operatorId int64, from, to model.LoanOperationStatus, action constants.LoanOperationAction, metadata *model.LoanOperationMetadata) error {
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

type loanPendingStatusTransition struct {
	rejectTarget  constants.LoanPendingStatus
	approveTarget constants.LoanPendingStatus
}

var loanPendingStatusTransMap = map[constants.LoanPendingStatus]loanPendingStatusTransition{
	constants.LoanPendingStatusQuoting: {
		rejectTarget:  constants.LoanPendingStatusNone,
		approveTarget: constants.LoanPendingStatusProcessing,
	},
	constants.LoanPendingStatusProcessing: {
		rejectTarget:  constants.LoanPendingStatusNone,
		approveTarget: constants.LoanPendingStatusUnderWriting,
	},
	constants.LoanPendingStatusUnderWriting: {
		rejectTarget:  constants.LoanPendingStatusProcessing,
		approveTarget: constants.LoanPendingStatusClosing,
	},
	constants.LoanPendingStatusClosing: {
		rejectTarget:  constants.LoanPendingStatusUnderWriting,
		approveTarget: constants.LoanPendingStatusNone,
	},
}

func (l *LoanBaseLogic) GetNextPendingStatus(currentStatus constants.LoanPendingStatus, action constants.LoanOperationAction) (nextStatus constants.LoanPendingStatus, err error) {

	trans, ok := loanPendingStatusTransMap[currentStatus]
	if !ok {
		return constants.LoanPendingStatusNone, errors.Errorf("Invalid Current Status %v", currentStatus)
	}
	switch action {
	case constants.LoanOperationActionApprove:
		nextStatus = trans.approveTarget
	case constants.LoanOperationActionReject:
		nextStatus = trans.rejectTarget
	default:
		nextStatus = constants.LoanPendingStatusNone
		err = errors.Errorf("Invalid Action %v", action)
	}
	if err != nil {
		return constants.LoanPendingStatusNone, err
	}
	if currentStatus != constants.LoanPendingStatusClosing && nextStatus == constants.LoanPendingStatusNone {
		return constants.LoanPendingStatusNone, errors.Errorf("Invalid Action %v For Current Status %v", action, currentStatus)
	}

	return nextStatus, err
}

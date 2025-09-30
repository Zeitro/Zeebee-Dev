package loanlogic

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
		return nil, errors.WithStack(errors.New("invalid lender company"))
	}

	if loan.Source != constants.LoanSourceSystem {
		return nil, errors.WithStack(errors.New("invalid loan source"))
	}

	return
}

// CreateLoanOperationLog creates a loan operation log record
func (l *LoanBaseLogic) CreateLoanOperationLog(tx *gorm.DB, loanId, operatorId int64, from, to types.LoanOperationStatus, action constants.LoanOperationAction, metadata *types.LoanOperationMetadata) error {
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
		return constants.LoanPendingStatusNone, errors.Errorf("Invalid current status %v", currentStatus)
	}
	switch action {
	case constants.LoanOperationActionApprove:
		nextStatus = trans.approveTarget
	case constants.LoanOperationActionReject:
		nextStatus = trans.rejectTarget
	default:
		nextStatus = constants.LoanPendingStatusNone
		err = errors.Errorf("Invalid action %v", action)
	}
	if err != nil {
		return constants.LoanPendingStatusNone, err
	}
	if currentStatus != constants.LoanPendingStatusClosing && nextStatus == constants.LoanPendingStatusNone {
		return constants.LoanPendingStatusNone, errors.Errorf("Invalid action %v for current status %v", action, currentStatus)
	}

	return nextStatus, err
}

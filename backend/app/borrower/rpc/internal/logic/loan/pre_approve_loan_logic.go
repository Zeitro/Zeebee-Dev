package loanlogic

import (
	"context"

	"github.com/pkg/errors"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
)

type PreApproveLoanLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPreApproveLoanLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PreApproveLoanLogic {
	return &PreApproveLoanLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PreApproveLoanLogic) PreApproveLoan(in *borrower.LoanPreApproveRequest) (*borrower.LoanPreApproveResponse, error) {
	// todo: add your logic here and delete this line

	return &borrower.LoanPreApproveResponse{}, nil
}

type CalculatePreApprovedLoanAmountResult struct {
	TotalAmount float32
	Rate        float32
}

func calculatePreApprovedLoanAmount(originalData *borrower.LoanPreApproveRequestData) (result *CalculatePreApprovedLoanAmountResult, err error) {
	if originalData == nil {
		err = errors.WithStack(errors.New("pre approve data is empty"))
		return
	}

	if originalData.Property == nil {
		err = errors.WithStack(errors.New("property information is empty"))
		return
	}

	if originalData.Borrower == nil {
		err = errors.WithStack(errors.New("borrower information is empty"))
	}

	if originalData.Loan == nil {
		err = errors.WithStack(errors.New("loan information is empty"))
	}

	return
}

func calculateTotalTotalAmount(originalData *borrower.LoanPreApproveRequestData) (result float32, err error) {

	return
}

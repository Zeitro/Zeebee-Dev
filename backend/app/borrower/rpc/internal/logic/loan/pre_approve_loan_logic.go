package loanlogic

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/duke-git/lancet/v2/mathutil"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/rpc/internal/global"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/constants"
	"zeebee-backend/common/pre_approve_engine"
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

	engines := pre_approve_engine.GetDefaultPreApproveEngines()
	input := pre_approve_engine.InputLoanInfo{
		LoanType:                constants.LoanType(in.Data.Loan.LoanType),
		LoanPurpose:             constants.LoanPurpose(in.Data.Loan.LoanPurpose),
		FicoScore:               int(in.Data.Borrower.EstCreditScore),
		TransactionHistoryCount: int(in.Data.Borrower.FlipPropertyCount + in.Data.Borrower.PurchaseRentalPropertyCount),
		Unit:                    int(in.Data.Property.Unit),
		LoanTerm:                int(in.Data.Loan.LoanTerm),
		EstCurrentPropertyValue: in.Data.Loan.EstPropertyValueAmount,
		EstAfterRepairValue:     in.Data.Loan.Arv,
		TotalAmount:             in.Data.Loan.TotalAmount,
		OriginalBudgetAmount:    in.Data.Loan.BudgetAmount,
		PurchaseAmount:          in.Data.Loan.PurchaseAmount,
		InitialAmount:           in.Data.Loan.InitialAmount,
		HoldbackAmount:          in.Data.Loan.HoldbackAmount,
		PrincipalAmount:         in.Data.Loan.PrincipalAmount,
		State:                   in.Data.Property.Address.State,
		City:                    in.Data.Property.Address.City,
		HasAssignmentFee:        in.Data.Loan.HasAssignmentFee,
		AssignmentFeeAmount:     in.Data.Loan.AssignmentFeeAmount,
		InterestReserveMonths:   int(in.Data.Loan.InterestReserveMonths),
	}
	result, err := engines["zeebee"].CalPreApprovedLoanAmount(input)
	if err != nil {
		return nil, err
	}

	responseData := &borrower.LoanPreApproveResponseData{
		Primary:     &borrower.PrimaryLoanPreApprove{},
		Secondaries: []*borrower.SecondaryLoanPreApprove{},
	}
	_ = copier.CopyWithOption(responseData.Primary, result, copier.Option{
		DeepCopy:      true,
		CaseSensitive: true,
	})
	if result.IsApproved {
		responseData.Primary.TotalAmount = mathutil.RoundToFloat(responseData.Primary.TotalAmount, 0)
		responseData.Primary.MonthlyPayment = mathutil.RoundToFloat(responseData.Primary.MonthlyPayment, 0)
		responseData.Primary.InitialAmount = mathutil.RoundToFloat(responseData.Primary.InitialAmount, 0)
		responseData.Primary.HoldbackAmount = mathutil.RoundToFloat(responseData.Primary.HoldbackAmount, 0)
		responseData.Primary.CashToClose = mathutil.RoundToFloat(responseData.Primary.CashToClose, 0)

		responseData.Primary.Rate = float32(mathutil.RoundToFloat(responseData.Primary.Rate*100.0, 2))
		responseData.Primary.Tlc = float32(mathutil.RoundToFloat(responseData.Primary.Tlc*100.0, 2))
		responseData.Primary.Ltv = float32(mathutil.RoundToFloat(responseData.Primary.Ltv*100.0, 2))
		responseData.Primary.Ltarv = float32(mathutil.RoundToFloat(responseData.Primary.Ltarv*100.0, 2))
	}
	_ = copier.CopyWithOption(&responseData.Secondaries, &global.Config.PreApproveLoanConf, copier.Option{
		DeepCopy:      true,
		CaseSensitive: true,
	})

	// save original data and result to redis
	cacheData := types.PreApproveLoanCacheData{}
	_ = copier.CopyWithOption(&cacheData.RequestData, in.Data, copier.Option{
		DeepCopy:      true,
		CaseSensitive: true,
	})
	_ = copier.CopyWithOption(&cacheData.ResultData, result, copier.Option{
		DeepCopy:      true,
		CaseSensitive: true,
	})
	key := fmt.Sprintf(constants.BorrowerPreApproveCacheKey, result.RequestId)
	contents, _ := json.Marshal(cacheData)
	err = global.GetRedis().SetExpire(l.ctx, key, string(contents), global.Config.PreApproveLoanExpireSeconds)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &borrower.LoanPreApproveResponse{
		Data: responseData,
	}, nil
}

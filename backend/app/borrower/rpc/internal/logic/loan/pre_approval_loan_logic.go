package loanlogic

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/duke-git/lancet/v2/mathutil"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"zeebee-backend/app/borrower/rpc/internal/email_tpl"
	"zeebee-backend/app/borrower/rpc/internal/global"
	"zeebee-backend/app/borrower/rpc/internal/repo"
	"zeebee-backend/app/borrower/rpc/internal/repo/model"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/constants"
	"zeebee-backend/common/email"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/preapproval_engine"
	"zeebee-backend/common/uuid"
)

type PreApprovalLoanLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPreApprovalLoanLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PreApprovalLoanLogic {
	return &PreApprovalLoanLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PreApprovalLoanLogic) PreApprovalLoan(in *borrower.LoanPreApprovalRequest) (*borrower.LoanPreApprovalResponse, error) {

	if in.Data.Property == nil || in.Data.Property.Address == nil {
		return nil, errors.New("Invalid Property Address")
	}

	// get geo data
	geo, err := l.svcCtx.UsGeographyRepo.FindOneByColumns(l.ctx, map[string]interface{}{
		"state_usps": in.Data.Property.Address.State,
		"city_name":  in.Data.Property.Address.City,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.Errorf("property address is not found")
		}

		return nil, errors.WithStack(err)
	}

	engines := preapproval_engine.GetDefaultPreApprovalEngines()
	input := preapproval_engine.InputLoanInfo{
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
		IsHighValueZone:         global.StateHighCostLimitMgr.IsHighCostZone(geo.StateUsps, geo.CityName),
	}
	result, err := engines["zeebee"].CalPreApprovalLoanAmount(input)
	if err != nil {
		return nil, err
	}
	requestId := in.Data.RequestId
	if len(requestId) == 0 {
		requestId, _ = uuid.GetSonyflakeIDString()
	}

	responseData := &borrower.LoanPreApprovalResponseData{
		Primary:     &borrower.PrimaryLoanPreApproval{},
		Secondaries: []*borrower.SecondaryLoanPreApproval{},
	}
	_ = copier.CopyWithOption(responseData.Primary, result, copier.Option{
		DeepCopy:      true,
		CaseSensitive: true,
	})
	responseData.Primary.RequestId = requestId
	if result.IsApproved {
		responseData.Primary.TotalAmount = mathutil.RoundToFloat(responseData.Primary.TotalAmount, 0)
		responseData.Primary.MonthlyPayment = mathutil.RoundToFloat(responseData.Primary.MonthlyPayment, 0)
		responseData.Primary.InitialAmount = mathutil.RoundToFloat(responseData.Primary.InitialAmount, 0)
		responseData.Primary.HoldbackAmount = mathutil.RoundToFloat(responseData.Primary.HoldbackAmount, 0)
		responseData.Primary.CashToClose = mathutil.RoundToFloat(responseData.Primary.CashToClose, 0)

		responseData.Primary.Rate = float32(mathutil.RoundToFloat(responseData.Primary.Rate*100.0, 2))
		responseData.Primary.Ltc = float32(mathutil.RoundToFloat(responseData.Primary.Ltc*100.0, 2))
		responseData.Primary.Ltv = float32(mathutil.RoundToFloat(responseData.Primary.Ltv*100.0, 2))
		responseData.Primary.Ltarv = float32(mathutil.RoundToFloat(responseData.Primary.Ltarv*100.0, 2))
	}
	_ = copier.CopyWithOption(&responseData.Secondaries, &global.Config.PreApprovalLoanConf, copier.Option{
		DeepCopy:      true,
		CaseSensitive: true,
	})

	// save original data and result to redis
	cacheData := types.PreApproveLoanCacheData{
		RequestData: types.LoanPreApprovalRequest{},
		ResultData:  &types.LoanPreApprovalResult{},
	}
	_ = copier.CopyWithOption(&cacheData.RequestData, in.Data, copier.Option{
		DeepCopy:      true,
		CaseSensitive: true,
	})
	_ = copier.CopyWithOption(&cacheData.ResultData, result, copier.Option{
		DeepCopy:      true,
		CaseSensitive: true,
	})
	key := fmt.Sprintf(constants.BorrowerPreApprovalCacheKey, requestId)
	contents, _ := json.Marshal(cacheData)
	err = global.GetRedis().SetExpire(l.ctx, key, string(contents), global.Config.PreApprovalLoanExpireSeconds)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	selfLenderCompany, err := l.svcCtx.LenderCompanyRepo.FindOneByColumns(l.ctx, map[string]interface{}{"owner_type": constants.LenderCompanyOwnerTypeSelf})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	go func() {
		ctx := context.WithoutCancel(l.ctx)
		err := sendPreApprovalNotification(ctx, l.svcCtx.LenderCompanyRepo, requestId, selfLenderCompany.Id)
		if err != nil {
			logx.WithContext(ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	return &borrower.LoanPreApprovalResponse{
		Data: responseData,
	}, nil
}

func sendPreApprovalNotification(ctx context.Context, repo *repo.LenderCompanyRepo, requestId string, lenderCompanyId int64) (err error) {
	var company *model.LenderCompany
	// check if send already
	sendKey := fmt.Sprintf(constants.BorrowerSendAprApprovalNotificationKey, lenderCompanyId, requestId)
	isSend, err := global.GetRedis().IsExisted(ctx, sendKey)
	if err != nil {
		return errors.WithStack(err)
	}
	if isSend {
		return nil
	}
	// check company status
	company, err = repo.FindOneById(ctx, lenderCompanyId)
	if err != nil {
		return errors.WithStack(err)
	}
	if len(company.ContactEmail) == 0 || company.Status != constants.LenderCompanyStatusActive {
		return nil
	}

	// get cache data
	cacheKey := fmt.Sprintf(constants.BorrowerPreApprovalCacheKey, requestId)
	cacheDataStr, err := global.GetRedis().Get(ctx, cacheKey)
	if err != nil {
		return errors.WithStack(err)
	}
	cacheData := types.PreApproveLoanCacheData{}
	err = json.Unmarshal([]byte(cacheDataStr), &cacheData)
	if err != nil {
		return errors.WithStack(err)
	}
	if cacheData.ResultData == nil {
		cacheData.ResultData = &types.LoanPreApprovalResult{}
	}
	// parse template
	tplVars := map[string]interface{}{
		"request": cacheData.RequestData,
		"result":  cacheData.ResultData,
		"company": company.Name,
	}
	emailContent := ""

	switch company.OwnerType {
	case constants.LenderCompanyOwnerTypeSelf:
		emailContent, err = email.ParseTemplate(email_tpl.SelfPreApprovalNotification, tplVars)
	case constants.LenderCompanyOwnerTypeThirdParty:
		emailContent, err = email.ParseTemplate(email_tpl.ThirdPartyPreApprovalNotification, tplVars)
	default:
		err = errors.New("Invalid Company Owner Type")
	}
	if err != nil {
		return errors.WithStack(err)
	}

	// send email
	emailMsg := email.Message{
		Subject: "Pre Approval Loan Notification",
		From:    global.Config.EmailConf.Username,
		To:      []string{company.ContactEmail},
		Body:    emailContent,
		IsHTML:  true,
	}

	err = global.EmailClient.SendEmail(emailMsg)
	if err != nil {
		return errors.WithStack(err)
	}
	// save send record
	_ = global.GetRedis().SetExpire(ctx, sendKey, "1", global.Config.PreApprovalLoanExpireSeconds)

	return err
}

package loanlogic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/duke-git/lancet/v2/mathutil"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"zeebee-backend/app/borrower/config"
	"zeebee-backend/app/borrower/rpc/internal/email_tpl"
	"zeebee-backend/app/borrower/rpc/internal/global"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/email"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/preapproval_engine"
	"zeebee-backend/common/uuid"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/repo"
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

func (l *PreApprovalLoanLogic) PreApprovalLoan(in *borrower.LoanPreApprovalRequest) (resp *borrower.LoanPreApprovalResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
				logx.Field("service_name", l.svcCtx.Config.Name),
			).Error(err)
		}
	}()

	if in.Data.Property == nil || in.Data.Property.Address == nil {
		return nil, errorx.New("Invalid property address")
	}

	// get geo data
	geo, err := l.svcCtx.UsGeographyRepo.FindOneByColumns(l.ctx, map[string]interface{}{
		"state_usps": in.Data.Property.Address.State,
		"city_name":  in.Data.Property.Address.City,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorx.Errorf("Property address is not found")
		}

		return nil, errorx.WithStack(err)
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
	_ = copierx.Assign[preapproval_engine.CalculatePreApprovalLoanAmountResult, borrower.PrimaryLoanPreApproval](result, responseData.Primary)
	if err != nil {
		return nil, err
	}

	result.RequestId = requestId
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
		responseData.Primary.LenderCompanyId = 1
		responseData.Primary.LenderCompanyName = "Zeebee"
	}

	_ = copierx.Assign[[]config.PreApprovalLoanConf, []*borrower.SecondaryLoanPreApproval](&global.Config.PreApprovalLoanConf, &responseData.Secondaries)

	// save original data and result to redis
	cacheData := types.PreApproveLoanCacheData{
		RequestData: types.LoanPreApprovalRequest{},
		ResultData:  types.LoanPreApprovalResult{},
	}
	_ = copierx.Assign[borrower.LoanPreApprovalRequestData, types.LoanPreApprovalRequest](in.Data, &cacheData.RequestData)
	_ = copierx.Assign[preapproval_engine.CalculatePreApprovalLoanAmountResult, types.LoanPreApprovalResult](result, &cacheData.ResultData)
	cacheData.RequestData.Borrower.FicoScore = in.Data.Borrower.EstCreditScore
	key := fmt.Sprintf(constants.BorrowerPreApprovalCacheKey, requestId)
	contents, _ := json.Marshal(cacheData)
	err = global.GetRedis().SetExpire(l.ctx, key, string(contents), global.Config.PreApprovalLoanExpireSeconds)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	selfOrganization, err := l.svcCtx.OrganizationRepo.FindOneByColumns(l.ctx, map[string]interface{}{
		"owner_type": constants.OrganizationOwnerTypeSelf,
		"type":       constants.UserTypeLender,
	})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	go func() {
		ctx := context.WithoutCancel(l.ctx)
		err := sendPreApprovalNotification(ctx, l.svcCtx.OrganizationRepo, requestId, selfOrganization.Id, false)
		if err != nil {
			logx.WithContext(ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	return &borrower.LoanPreApprovalResponse{
		Data: responseData,
	}, nil
}

func sendPreApprovalNotification(ctx context.Context, repo *repo.OrganizationRepo, requestId string, organizationId int64, isThirdParty bool) (err error) {
	var organization *model.Organization
	// check if send already
	sendKey := fmt.Sprintf(constants.BorrowerSendAprApprovalNotificationKey, organizationId, requestId)
	isSend, err := global.GetRedis().IsExisted(ctx, sendKey)
	if err != nil {
		return errorx.WithStack(err)
	}
	if isSend {
		return nil
	}
	// check company status
	organization, err = repo.FindOneById(ctx, organizationId)
	if err != nil {
		return errorx.WithStack(err)
	}
	if len(organization.ContactEmail) == 0 || organization.Status != constants.OrganizationStatusActive {
		return nil
	}

	// get cache data
	cacheKey := fmt.Sprintf(constants.BorrowerPreApprovalCacheKey, requestId)
	cacheDataStr, err := global.GetRedis().Get(ctx, cacheKey)
	if err != nil {
		return errorx.WithStack(err)
	}
	cacheData := types.PreApproveLoanCacheData{}
	err = json.Unmarshal([]byte(cacheDataStr), &cacheData)
	if err != nil {
		return errorx.WithStack(err)
	}
	// if cacheData.ResultData == nil {
	// 	cacheData.ResultData = types.LoanPreApprovalResult{}
	// }
	cacheData.ResultData.TotalAmount = mathutil.RoundToFloat(cacheData.ResultData.TotalAmount, 0)
	cacheData.ResultData.MonthlyPayment = mathutil.RoundToFloat(cacheData.ResultData.MonthlyPayment, 0)
	cacheData.ResultData.InitialAmount = mathutil.RoundToFloat(cacheData.ResultData.InitialAmount, 0)
	cacheData.ResultData.HoldbackAmount = mathutil.RoundToFloat(cacheData.ResultData.HoldbackAmount, 0)
	cacheData.ResultData.CashToClose = mathutil.RoundToFloat(cacheData.ResultData.CashToClose, 0)

	cacheData.ResultData.Rate = mathutil.RoundToFloat(cacheData.ResultData.Rate*100.0, 2)
	cacheData.ResultData.Ltc = float32(mathutil.RoundToFloat(cacheData.ResultData.Ltc*100.0, 2))
	cacheData.ResultData.Ltv = float32(mathutil.RoundToFloat(cacheData.ResultData.Ltv*100.0, 2))
	cacheData.ResultData.Ltarv = float32(mathutil.RoundToFloat(cacheData.ResultData.Ltarv*100.0, 2))

	// parse template
	tplVars := map[string]interface{}{
		"request":        cacheData.RequestData,
		"result":         cacheData.ResultData,
		"company":        organization.Name,
		"url":            global.Config.LenderUrl,
		"is_third_party": isThirdParty,
	}
	emailContent := ""
	emailContent, err = email.ParseTemplate(email_tpl.SelfPreApprovalNotification, tplVars)

	if err != nil {
		return errorx.WithStack(err)
	}

	targetEmail := organization.ContactEmail
	if isThirdParty {
		targetEmail = global.Config.ZeebeeAgentEmail
	}

	// send email
	emailMsg := email.Message{
		Subject: "New Pre-Approval-Ready for Review in Zeebee",
		From:    global.Config.EmailConf.Username,
		To:      []string{targetEmail},
		Body:    emailContent,
		IsHTML:  true,
	}

	err = global.EmailClient.SendEmail(emailMsg)
	if err != nil {
		return errorx.WithStack(err)
	}
	// save send record
	_ = global.GetRedis().SetExpire(ctx, sendKey, "1", global.Config.PreApprovalLoanExpireSeconds)

	return err
}

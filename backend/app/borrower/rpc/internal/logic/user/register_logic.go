package userlogic

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/duke-git/lancet/v2/random"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"zeebee-backend/app/borrower/rpc/internal/email_tpl"
	"zeebee-backend/app/borrower/rpc/internal/global"
	"zeebee-backend/app/borrower/rpc/internal/logic/loan"
	"zeebee-backend/app/borrower/rpc/internal/repo/model"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/constants"
	email2 "zeebee-backend/common/email"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/limiter"
	"zeebee-backend/common/stringx"
	commontypes "zeebee-backend/common/types"
	"zeebee-backend/common/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *borrower.UserRegisterRequest) (resp *borrower.UserRegisterResponse, err error) {

	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	// check if user already exists
	var existedUser *model.User
	existedUser, _ = l.svcCtx.UserRepo.GetBorrowerUserByEmail(l.ctx, in.Data.Email)

	if existedUser != nil && existedUser.Id > 0 {
		ctx := context.WithoutCancel(l.ctx)
		if existedUser.Status == constants.UserStatusUnverified {
			// send verify email
			go func() {
				err := sendRegisterVerifyUrlEmail(ctx, in.Data.Email, existedUser.Borrower.FirstName)
				if err != nil {
					err = errors.Wrap(err, "send verify email error")
					l.WithContext(ctx).Error(err)
				}
			}()
		} else {
			// send login email
			go func() {
				err := sendLoginUrlEmail(ctx, in.Data.Email, existedUser.Borrower.FirstName)
				if err != nil {
					err = errors.Wrap(err, "send login email error")
					l.WithContext(ctx).Error(err)
				}
			}()
		}

		return &borrower.UserRegisterResponse{
			Data: &borrower.UserRegisterResponseData{},
		}, nil
	}
	var pwd string
	pwd, err = utils.HashPassword(in.Data.Password)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	userModel := &model.User{
		Email:        in.Data.Email,
		PasswordHash: pwd,
		Status:       constants.UserStatusUnverified,
		Type:         constants.UserTypeBorrower,
	}
	borrowerModel := &model.Borrower{
		FirstName: in.Data.FirstName,
		LastName:  in.Data.LastName,
		FullName:  stringx.FormatFullName(in.Data.FirstName, "", in.Data.LastName),
	}
	err = l.svcCtx.UserRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		// create user
		userModel, err = l.svcCtx.UserRepo.CreateWithTx(l.ctx, tx, userModel, false)
		if err != nil {
			return err
		}
		// create borrower
		borrowerModel.UserId = userModel.Id
		_, err = l.svcCtx.BorrowerRepo.CreateWithTx(l.ctx, tx, borrowerModel, false)
		return err
	})
	if err != nil {
		return nil, err
	}

	go func() {
		// send account validate email
		ctx := context.WithoutCancel(l.ctx)
		err := sendRegisterVerifyUrlEmail(ctx, in.Data.Email, borrowerModel.FirstName)
		if err != nil {
			l.WithContext(ctx).Error(err)
		}
		if err != nil {
			err = errors.Wrap(err, "send register verify email error")
			l.WithContext(ctx).Error(err)
		}
	}()

	go func() {
		if len(in.Data.PreApprovalRequestId) == 0 || in.Data.TargetLenderCompanyId == 0 {
			return
		}
		// create draft loan
		err := l.handlePreapproveLoan(context.WithoutCancel(l.ctx), borrowerModel.Id, in.Data.PreApprovalRequestId, in.Data.TargetLenderCompanyId)
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	return &borrower.UserRegisterResponse{
		Data: &borrower.UserRegisterResponseData{},
	}, err
}

func (l *RegisterLogic) handlePreapproveLoan(ctx context.Context, borrowerId int64, requestId string, lenderCompanyId int64) (err error) {
	// get pre-approve loan data from cache
	key := fmt.Sprintf(constants.BorrowerPreApprovalCacheKey, requestId)
	cacheDataStr, err := global.GetRedis().Get(ctx, key)
	if err != nil {
		return errors.WithStack(err)
	}

	if len(cacheDataStr) == 0 {
		return errors.WithStack(errors.New("pre-approve loan data not found in cache"))
	}

	defer func() {
		if err != nil {
			// delete cache data
			_ = global.GetRedis().Del(ctx, key)
		}
	}()

	// parse cache data
	var cacheData types.PreApproveLoanCacheData
	err = json.Unmarshal([]byte(cacheDataStr), &cacheData)
	if err != nil {
		return errors.WithStack(err)
	}
	preApproveRequest := cacheData.RequestData
	calcResult := cacheData.ResultData
	if calcResult.LenderCompanyId != lenderCompanyId {
		calcResult = nil
	}

	// get lender company
	lenderCompany, err := l.svcCtx.LenderCompanyRepo.FindOneById(ctx, lenderCompanyId)
	if err != nil {
		return errors.WithStack(err)
	}

	if lenderCompany.Status != constants.LenderCompanyStatusActive {
		return errors.WithStack(errors.New("lender company is not active"))
	}
	// prepare data
	var borrowerProperty *model.BorrowerProperty
	var borrowerApplicant *model.BorrowerApplicant
	err = l.svcCtx.UserRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		var innerErr error

		eg := errorx.NewErrorGroup()
		eg.Do(
			// get BorrowerProperty
			func() error {
				var err error
				borrowerProperty, err = l.svcCtx.BorrowerPropertyRepo.GetOrCreateBorrowerPropertyWithTx(l.ctx, tx, borrowerId, commontypes.Address{
					State:    preApproveRequest.Property.Address.State,
					City:     preApproveRequest.Property.Address.City,
					Address1: preApproveRequest.Property.Address.Address1,
					Address2: preApproveRequest.Property.Address.Address2,
					ZipCode:  preApproveRequest.Property.Address.ZipCode,
				})
				return err
			},
			// create BorrowerApplicant
			func() error {
				var err error
				borrowerApplicant = &model.BorrowerApplicant{
					BorrowerId:     borrowerId,
					FirstName:      preApproveRequest.Borrower.Name,
					Email:          preApproveRequest.Borrower.Email,
					Phone:          preApproveRequest.Borrower.Phone,
					EstCreditScore: preApproveRequest.Borrower.FicoScore,
				}
				borrowerApplicant, err = l.svcCtx.BorrowerApplicantRepo.CreateWithTx(l.ctx, tx, borrowerApplicant, false)
				return err
			},
		)
		innerErr = eg.Error()
		return innerErr
	})
	if err != nil {
		return err
	}

	// create draft loan
	saveDraftLogic := loanlogic.NewSaveDraftLoanLogic(ctx, l.svcCtx)
	// build LoanApplyRequest
	loanApplyRequest := &borrower.LoanApplyRequest{
		Data: &borrower.LoanApplyRequestData{
			IsPreapproveLoan: true,
			LoanPurpose:      preApproveRequest.Loan.LoanPurpose,
			LoanType:         preApproveRequest.Loan.LoanType,
			TotalAmount:      preApproveRequest.Loan.TotalAmount,
			InitialAmount:    preApproveRequest.Loan.InitialAmount,
			LoanProperty: &borrower.LoanPropertyApplyEntry{
				BorrowerPropertyId:  borrowerProperty.Id,
				FuturePropertyType:  preApproveRequest.Property.PropertyType,
				FutureUnit:          int32(preApproveRequest.Property.Unit),
				PurchaseAmount:      preApproveRequest.Loan.PurchaseAmount,
				PurchaseDate:        preApproveRequest.Loan.PurchaseDate,
				HasAssignmentFee:    preApproveRequest.Loan.HasAssignmentFee,
				AssignmentFeeAmount: preApproveRequest.Loan.AssignmentFeeAmount,
				Budget: &borrower.LoanBudgetEntry{
					OriginalBudget: preApproveRequest.Loan.BudgetAmount,
				},
				EstPropertyValueAmount: preApproveRequest.Loan.EstPropertyValueAmount,
			},
			HoldbackAmount:  preApproveRequest.Loan.HoldbackAmount,
			PrincipalAmount: preApproveRequest.Loan.PrincipalAmount,
			LenderCompanyId: lenderCompanyId,
			BorrowerId:      borrowerId,
			Applicants: []*borrower.LoanApplicantApplyEntry{
				{
					IsPrimary:           true,
					BorrowerApplicantId: borrowerApplicant.Id,
				},
			},
			LoanTerm:       preApproveRequest.Loan.LoanTerm,
			Arv:            preApproveRequest.Loan.Arv,
			EstCreditScore: preApproveRequest.Borrower.FicoScore,
		},
	}
	if calcResult != nil {
		loanApplyRequest.Data.ReservesAmount = calcResult.InterestReservedAmount
		loanApplyRequest.Data.TotalAmount = calcResult.TotalAmount
		loanApplyRequest.Data.InitialAmount = calcResult.InitialAmount
		loanApplyRequest.Data.Ltv = float64(calcResult.Ltv)
	}

	// save draft loan
	_, err = saveDraftLogic.SaveDraftLoan(loanApplyRequest)
	if err != nil {
		return errors.WithStack(err)
	}

	//  send email to lender
	// err = sendPreApprovalLoanEmail(lenderCompany)
	// if err != nil {
	// 	return errors.WithStack(err)
	// }

	return err
}

// sendPreApprovalLoanEmail send pre approval loan email
func sendPreApprovalLoanEmail(lenderCompany *model.LenderCompany) (err error) {
	// check lender company
	if lenderCompany == nil || lenderCompany.Id == 0 {
		return err
	}
	if len(lenderCompany.ContactEmail) == 0 || lenderCompany.Status != constants.LenderCompanyStatusActive {
		return err
	}
	// send email
	data := map[string]interface{}{}
	content, err := email2.ParseTemplate(email_tpl.SelfPreApprovalNotification, data)
	if err != nil {
		return errors.WithStack(err)
	}
	msg := email2.Message{
		Subject: "PreApproval Loan Notification",
		From:    global.Config.EmailConf.Username,
		To:      []string{lenderCompany.ContactEmail},
		Body:    content,
		IsHTML:  true,
	}
	err = global.EmailClient.SendEmail(msg)

	return errors.WithStack(err)
}

// sendLoginUrlEmail send login url email
func sendLoginUrlEmail(ctx context.Context, email string, firstName string) (err error) {
	key := stringx.Format(constants.BorrowerLoginEmail, email)
	limitOpt := limiter.WithRedisLimit(
		global.Config.LimiterConf.SendLoginEmail.Rate,
		time.Second*time.Duration(global.Config.LimiterConf.SendLoginEmail.Period),
	)
	ok, err := global.SendLoginEmailLimiter.Limit(ctx, key, limitOpt)
	if err != nil {
		return errors.WithStack(err)
	}
	if !ok {
		return errors.WithStack(errors.Errorf("Email %s try too many times", email))
	}

	// send login email
	data := map[string]interface{}{
		"url":        fmt.Sprintf("%s", global.Config.LoginUrl),
		"first_name": firstName,
	}
	content, err := email2.ParseTemplate(email_tpl.LoginUrlEmail, data)
	if err != nil {
		return errors.WithStack(err)
	}
	msg := email2.Message{
		Subject: "Login to your account",
		From:    global.Config.EmailConf.Username,
		To:      []string{email},
		Body:    content,
		IsHTML:  true,
	}
	err = global.EmailClient.SendEmail(msg)

	return errors.WithStack(err)
}

// sendRegisterVerifyUrlEmail send register verify url email
func sendRegisterVerifyUrlEmail(ctx context.Context, email string, firstName string) (err error) {
	var code string
	// generate verify code
	codeKey := stringx.Format(constants.BorrowerRegisterVerifyCode, email)
	code, _ = global.GetRedis().Get(ctx, codeKey)
	if len(code) == 0 {
		code = random.RandNumeral(10)
		err = global.GetRedis().SetExpire(ctx, codeKey, code, global.Config.VerifyEmailExpireSeconds)
		if err != nil {
			return errors.WithStack(err)
		}
	}
	// send verify email
	data := map[string]interface{}{
		"url":        fmt.Sprintf("%s?code=%s&email=%s", global.Config.VerifyEmailUrl, code, email),
		"first_name": firstName,
	}
	content, err := email2.ParseTemplate(email_tpl.RegisterVerifyEmail, data)
	if err != nil {
		return errors.WithStack(err)
	}
	msg := email2.Message{
		Subject: "Verify your email",
		From:    global.Config.EmailConf.Username,
		To:      []string{email},
		Body:    content,
		IsHTML:  true,
	}
	err = global.EmailClient.SendEmail(msg)

	return errors.WithStack(err)
}

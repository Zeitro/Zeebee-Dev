package userlogic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/duke-git/lancet/v2/random"
	"gorm.io/gorm"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/rpc/internal/email_tpl"
	"zeebee-backend/app/borrower/rpc/internal/global"
	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/logic/loan"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/constants"
	email2 "zeebee-backend/common/email"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/limiter"
	"zeebee-backend/common/stringx"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/model"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.NewBaseLogic(ctx, svcCtx),
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
	existedUser, _ = l.svcCtx.UserRepo.GetUserByEmail(l.ctx, in.Data.Email, constants.UserTypeBorrower)

	if existedUser != nil && existedUser.Id > 0 {
		ctx := context.WithoutCancel(l.ctx)
		if existedUser.Status == constants.UserStatusUnverified {
			// send verify email
			go func() {
				err := sendRegisterVerifyUrlEmail(ctx, in.Data.Email, existedUser.FirstName)
				if err != nil {
					err = errorx.Wrap(err, "Send Verify Email Error")
					l.WithContext(ctx).Error(err)
				}
			}()
		} else {
			// send login email
			go func() {
				err := sendLoginUrlEmail(ctx, in.Data.Email, existedUser.FirstName)
				if err != nil {
					err = errorx.Wrap(err, "Send Login Email Error")
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
		return nil, errorx.WithStack(err)
	}
	userModel := &model.User{
		Email:        in.Data.Email,
		PasswordHash: pwd,
		Status:       constants.UserStatusUnverified,
		Type:         constants.UserTypeBorrower,
		FirstName:    in.Data.FirstName,
		LastName:     in.Data.LastName,
		FullName:     stringx.FormatFullName(in.Data.FirstName, "", in.Data.LastName),
	}
	organizationModel := &model.Organization{
		Name:      in.Data.OrganizationName,
		Type:      constants.UserTypeBorrower,
		Status:    constants.OrganizationStatusActive,
		OwnerType: constants.OrganizationOwnerTypeNone,
	}

	userModel, err = l.baseLogic.UserLogic.RegisterWithOrganization(l.ctx, userModel, organizationModel, constants.UserTypeBorrower)
	if err != nil {
		return nil, err
	}

	newCxt := context.WithoutCancel(l.ctx)
	go func() {
		// send account validate email
		ctx := context.WithoutCancel(newCxt)
		err := sendRegisterVerifyUrlEmail(ctx, in.Data.Email, userModel.FirstName)
		if err != nil {
			l.WithContext(ctx).Error(err)
		}
		if err != nil {
			err = errorx.Wrap(err, "Send Register Verify Email Error")
			l.WithContext(ctx).Error(err)
		}
	}()

	go func() {
		if len(in.Data.PreApprovalRequestId) == 0 || in.Data.LenderOrganizationId == 0 {
			return
		}
		// create draft loan
		ctx := context.WithoutCancel(newCxt)
		err := l.handlePreapproveLoan(newCxt, userModel.Id, userModel.FullName, in.Data.PreApprovalRequestId, in.Data.LenderOrganizationId)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				l.WithContext(ctx).
					WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).
					Info("Pre-Approve Loan Processing Canceled")
			} else {
				l.WithContext(ctx).
					WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).
					Error(err)
			}
		}
	}()

	return &borrower.UserRegisterResponse{
		Data: &borrower.UserRegisterResponseData{},
	}, err
}

func (l *RegisterLogic) handlePreapproveLoan(ctx context.Context, creatorId int64, creatorName string, requestId string, organizationId int64) (err error) {
	// get pre-approve loan data from cache
	key := fmt.Sprintf(constants.BorrowerPreApprovalCacheKey, requestId)
	cacheDataStr, err := global.GetRedis().Get(ctx, key)
	if err != nil {
		return errorx.WithStack(err)
	}

	if len(cacheDataStr) == 0 {
		return errorx.WithStack(errorx.New("Pre-Approve Loan Data Not Found In Cache"))
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
		return errorx.WithStack(err)
	}
	preApproveRequest := cacheData.RequestData
	calcResult := cacheData.ResultData
	if calcResult.LenderCompanyId != organizationId {
		calcResult = nil
	}

	// get lender company
	lenderOrganization, err := l.svcCtx.OrganizationRepo.FindOneById(ctx, organizationId)
	if err != nil {
		return errorx.WithStack(err)
	}

	if lenderOrganization.Status != constants.OrganizationStatusActive {
		return errorx.WithStack(errorx.New("Lender Organization Is Not Active"))
	}
	// prepare data
	var borrowerProperty *model.BorrowerProperty
	var borrowerApplicant *model.BorrowerApplicant
	err = l.svcCtx.UserRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var innerErr error

		eg := errorx.NewErrorGroup()
		eg.Do(
			// get BorrowerProperty
			func() error {
				var err error
				borrowerProperty, err = l.svcCtx.BorrowerPropertyRepo.GetOrCreateBorrowerPropertyWithTxForBorrower(ctx, tx, creatorId, creatorName, model.Address{
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
					CreatorId:      creatorId,
					CreatorName:    creatorName,
					OrganizationId: organizationId,
					FirstName:      preApproveRequest.Borrower.Name,
					Email:          preApproveRequest.Borrower.Email,
					Phone:          preApproveRequest.Borrower.Phone,
					EstCreditScore: preApproveRequest.Borrower.FicoScore,
				}
				borrowerApplicant, err = l.svcCtx.BorrowerApplicantRepo.CreateWithTx(ctx, tx, borrowerApplicant, false)
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
				Arv:                    preApproveRequest.Loan.Arv,
			},
			HoldbackAmount:  preApproveRequest.Loan.HoldbackAmount,
			PrincipalAmount: preApproveRequest.Loan.PrincipalAmount,
			LenderCompanyId: organizationId,
			UserId:          creatorId,
			Applicants: []*borrower.LoanApplicantApplyEntry{
				{
					IsPrimary:           true,
					BorrowerApplicantId: borrowerApplicant.Id,
				},
			},
			LoanTerm:       preApproveRequest.Loan.LoanTerm,
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
		return errorx.WithStack(err)
	}
	return err
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
		return errorx.WithStack(err)
	}
	if !ok {
		return errorx.WithStack(errorx.Errorf("Email %s Try Too Many Times", email))
	}

	// send login email
	data := map[string]interface{}{
		"url":        fmt.Sprintf("%s", global.Config.LoginUrl),
		"first_name": firstName,
	}
	content, err := email2.ParseTemplate(email_tpl.LoginUrlEmail, data)
	if err != nil {
		return errorx.WithStack(err)
	}
	msg := email2.Message{
		Subject: "Login to your account",
		From:    global.Config.EmailConf.Username,
		To:      []string{email},
		Body:    content,
		IsHTML:  true,
	}
	err = global.EmailClient.SendEmail(msg)

	return errorx.WithStack(err)
}

// sendRegisterVerifyUrlEmail send register verify url email
func sendRegisterVerifyUrlEmail(ctx context.Context, email string, firstName string) (err error) {
	var code string
	// generate verify code
	codeKey := stringx.Format(constants.UserRegisterVerifyCode, constants.UserTypeBorrower, email)
	code, _ = global.GetRedis().Get(ctx, codeKey)
	if len(code) == 0 {
		code = random.RandNumeral(10)
		err = global.GetRedis().SetExpire(ctx, codeKey, code, global.Config.VerifyEmailExpireSeconds)
		if err != nil {
			return errorx.WithStack(err)
		}
	}
	// send verify email
	data := map[string]interface{}{
		"url":        fmt.Sprintf("%s?code=%s&email=%s", global.Config.VerifyEmailUrl, code, email),
		"first_name": firstName,
	}
	content, err := email2.ParseTemplate(email_tpl.RegisterVerifyEmail, data)
	if err != nil {
		return errorx.WithStack(err)
	}
	msg := email2.Message{
		Subject: "Verify your email",
		From:    global.Config.EmailConf.Username,
		To:      []string{email},
		Body:    content,
		IsHTML:  true,
	}
	err = global.EmailClient.SendEmail(msg)

	return errorx.WithStack(err)
}

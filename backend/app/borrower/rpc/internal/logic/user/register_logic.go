package userlogic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/duke-git/lancet/v2/random"
	"google.golang.org/grpc/metadata"
	"gorm.io/gorm"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/rpc/internal/email_tpl"
	"zeebee-backend/app/borrower/rpc/internal/global"
	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/logic/loan"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	email2 "zeebee-backend/common/email"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/stringx"
	"zeebee-backend/common/timex"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/limiter"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/proto"
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
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *RegisterLogic) Register(in *borrower.UserRegisterRequest) (resp *proto.EmptyResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
				logx.Field("service_name", l.svcCtx.Config.Name),
			).Error(err)
		}
	}()
	// check if user already exists
	var existedUser *model.User
	// var organizationUserModel *model.OrganizationUser
	in.Data.Email = strings.ToLower(in.Data.Email)
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

		return &proto.EmptyResponse{}, nil
	}
	var pwd string
	pwd, err = utils.HashPassword(in.Data.Password)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	userInput := dto.UserRegisterRequestDto{
		Email:        in.Data.Email,
		PasswordHash: pwd,
		Status:       constants.UserStatusUnverified,
		Type:         constants.UserTypeBorrower,
		FirstName:    in.Data.FirstName,
		LastName:     in.Data.LastName,
		FullName:     stringx.FormatFullName(in.Data.FirstName, "", in.Data.LastName),
	}
	organizationInput := dto.OrganizationRegisterRequestDto{
		Name:      in.Data.OrganizationName,
		Type:      constants.UserTypeBorrower,
		Status:    constants.OrganizationStatusActive,
		OwnerType: constants.OrganizationOwnerTypeNone,
	}

	userModel, organizationUserModel, err := l.baseLogic.UserLogic.RegisterWithOrganization(l.ctx, userInput, organizationInput, constants.UserTypeBorrower, false)
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
		var err error
		var tokenInfo *dto.UserJwtTokenDto
		ctx := context.WithoutCancel(newCxt)
		md, _ := metadata.FromIncomingContext(ctx)
		newMd := md.Copy()
		logx.WithContext(ctx).WithFields(
			logx.Field("pre_approval_request_id", in.Data.PreApprovalRequestId),
			logx.Field("lender_organization_id", in.Data.LenderOrganizationId),
		).Info("create draft loan for preapproval loan")
		defer func() {
			panicInfo := recover()
			if panicInfo != nil {
				err = errorx.Wrap(panicInfo.(error), "Send Register Verify Email Panic Error")
				l.WithContext(ctx).Error(err)
			}

			if err != nil {
				err = errorx.Wrap(err, "Send Register Verify Email Error")
				l.WithContext(ctx).Error(err)
			}
		}()
		if len(in.Data.PreApprovalRequestId) == 0 || in.Data.LenderOrganizationId == 0 {
			return
		}

		tokenInfo, err = l.baseLogic.UserLogic.GenerateJwtToken(userModel, organizationUserModel, constants.UserTypeBorrower)
		if err != nil {
			return
		}
		newMd.Set(utils.JwtTokenCtxKey, tokenInfo.AccessToken)
		ctx = metadata.NewIncomingContext(ctx, newMd)
		// create a draft loan
		err = l.handlePreapproveLoan(
			ctx,
			organizationUserModel.Id,
			organizationUserModel.FullName,
			in.Data.PreApprovalRequestId,
			organizationUserModel.OrganizationId,
			in.Data.LenderOrganizationId)
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

	return &proto.EmptyResponse{}, err
}

func (l *RegisterLogic) handlePreapproveLoan(ctx context.Context, creatorId int64, creatorName string, requestId string, borrowerOrganizationId int64, lenderOrganizationId int64) (err error) {
	// get pre-approve loan data from cache
	key := fmt.Sprintf(constants.BorrowerPreApprovalCacheKey, requestId)
	cacheDataStr, err := global.GetRedis().Get(ctx, key)
	if err != nil {
		return errorx.WithStack(err)
	}

	if len(cacheDataStr) == 0 {
		return errorx.WithStack(errorx.New("Pre-approve loan data not found in cache"))
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
	if calcResult.LenderOrganizationId != lenderOrganizationId {
		calcResult = types.LoanPreApprovalResult{}
	}

	// get lender company
	lenderOrganization, err := l.svcCtx.OrganizationRepo.FindOneById(ctx, lenderOrganizationId)
	if err != nil {
		return errorx.WithStack(err)
	}

	if lenderOrganization.Status != constants.OrganizationStatusActive {
		return errorx.WithStack(errorx.New("Lender organization is not active"))
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
				borrowerProperty, err = l.svcCtx.BorrowerPropertyRepo.GetOrCreateBorrowerPropertyWithTxForBorrower(ctx, tx, borrowerOrganizationId, creatorId, creatorName, model.Address{
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
					OrganizationId: borrowerOrganizationId,
					FirstName:      preApproveRequest.Borrower.Name,
					Email:          preApproveRequest.Borrower.Email,
					Phone:          preApproveRequest.Borrower.Phone,
					EstCreditScore: preApproveRequest.Borrower.FicoScore,
				}
				borrowerApplicant, err = l.svcCtx.BorrowerApplicantRepo.CreateWithTx(ctx, tx, borrowerApplicant, false)
				return err
			},
		)
		innerErr = eg.Wait()
		return innerErr
	})
	if err != nil {
		return err
	}

	loanPreapprovalInfo := &borrower.LoanPreapproveInfo{}
	loanPreapprovalInfo.Input, _ = copierx.Trans[types.LoanPreApprovalRequest, borrower.LoanPreApprovalRequestData](&preApproveRequest)
	loanPreapprovalInfo.Result, _ = copierx.Trans[types.LoanPreApprovalResult, borrower.LoanPreapproveResult](&calcResult)

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
			LoanProperty: &borrower.LoanPropertyApplyEntity{
				BorrowerPropertyId:  borrowerProperty.Id,
				FuturePropertyType:  preApproveRequest.Property.PropertyType,
				FutureUnit:          int32(preApproveRequest.Property.Unit),
				PurchaseAmount:      preApproveRequest.Loan.PurchaseAmount,
				PurchaseDate:        preApproveRequest.Loan.PurchaseDate,
				HasAssignmentFee:    preApproveRequest.Loan.HasAssignmentFee,
				AssignmentFeeAmount: preApproveRequest.Loan.AssignmentFeeAmount,
				Budget: &proto.LoanBudgetEntity{
					OriginalBudget: preApproveRequest.Loan.BudgetAmount,
				},
				EstPropertyValueAmount: preApproveRequest.Loan.EstPropertyValueAmount,
				Arv:                    preApproveRequest.Loan.Arv,
			},
			HoldbackAmount:  preApproveRequest.Loan.HoldbackAmount,
			PrincipalAmount: preApproveRequest.Loan.PrincipalAmount,
			LenderCompanyId: lenderOrganizationId,
			// UserId:          creatorId,
			Applicants: []*borrower.LoanApplicantApplyEntity{
				{
					IsPrimary:           true,
					BorrowerApplicantId: borrowerApplicant.Id,
				},
			},
			LoanTerm:           preApproveRequest.Loan.LoanTerm,
			EstCreditScore:     preApproveRequest.Borrower.FicoScore,
			LoanPreapproveInfo: loanPreapprovalInfo,
		},
	}
	if len(calcResult.RequestId) > 0 || calcResult.TotalAmount > 0 {
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
	key := stringx.Format(constants.LoginEmail, constants.UserTypeBorrower, email)
	ok, err := global.LimiterMgr.Limit(ctx, limiter.StrategySendLoginEmail, key)

	if err != nil {
		return errorx.WithStack(err)
	}
	if !ok {
		return errorx.WithStack(errorx.Errorf("Email %s try too many times", email))
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
		"url":         fmt.Sprintf("%s?code=%s&email=%s", global.Config.VerifyEmailUrl, code, email),
		"first_name":  firstName,
		"expire_time": timex.DurationToString(time.Second*time.Duration(global.Config.VerifyEmailExpireSeconds), true),
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

package userlogic

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/metadata"
	"gorm.io/gorm"

	"zeebee-backend/app/borrower/rpc/internal/global"
	"zeebee-backend/app/borrower/rpc/internal/logic"
	loanlogic "zeebee-backend/app/borrower/rpc/internal/logic/loan"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/proto"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

// Register proxies identity creation to tenant (the sole identity authority)
// and keeps the borrower-only pre-approval side effect local: after tenant
// creates the account, borrower asynchronously builds the draft loan from the
// cached pre-approval result. A side-effect failure never fails registration.
func (l *RegisterLogic) Register(in *borrower.UserRegisterRequest) (resp *proto.EmptyResponse, err error) {
	preApprovalJSON, _ := json.Marshal(map[string]any{
		"pre_approval_request_id": in.Data.PreApprovalRequestId,
		"lender_organization_id":  in.Data.LenderOrganizationId,
	})
	registerResp, err := l.svcCtx.TenantClient.User.Register(l.ctx, &tenant.RegisterRequest{
		Data: &tenant.RegisterRequestData{
			UserType:                string(constants.UserTypeBorrower),
			Email:                   in.Data.Email,
			Password:                in.Data.Password,
			FirstName:               in.Data.FirstName,
			LastName:                in.Data.LastName,
			OrganizationName:        in.Data.OrganizationName,
			OwnerType:               string(constants.OrganizationOwnerTypeThirdParty),
			BorrowerPreApprovalJson: string(preApprovalJSON),
		},
	})
	if err != nil {
		return nil, err
	}

	if len(in.Data.PreApprovalRequestId) > 0 && in.Data.LenderOrganizationId > 0 && registerResp.Data.AccessToken != "" {
		go l.handlePreapproveLoanSideEffect(
			registerResp.Data.AccessToken,
			in.Data.PreApprovalRequestId,
			in.Data.LenderOrganizationId,
		)
	}

	return &proto.EmptyResponse{}, nil
}

// handlePreapproveLoanSideEffect fetches the just-created identity and creates
// the draft pre-approval loan. Fire-and-forget; errors are logged only.
func (l *RegisterLogic) handlePreapproveLoanSideEffect(accessToken, requestId string, lenderOrganizationId int64) {
	ctx := context.WithoutCancel(l.ctx)
	defer func() {
		if panicInfo := recover(); panicInfo != nil {
			l.svcCtx.Logger.WithContext(ctx).
				WithFields(logx.Field("stacks", fmt.Sprint(panicInfo))).
				Error("pre-approval loan side effect panic")
		}
	}()

	md, _ := metadata.FromIncomingContext(ctx)
	newMd := md.Copy()
	newMd.Set(utils.JwtTokenCtxKey, accessToken)
	newMd.Set(utils.UserTypeCtxKey, string(constants.UserTypeBorrower))
	ctx = metadata.NewIncomingContext(ctx, newMd)
	claim, err := l.baseLogic.GetUserClaimFromCtx(ctx)
	if err != nil {
		l.svcCtx.Logger.WithContext(ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		return
	}

	if err = l.handlePreapproveLoan(ctx, claim.OrganizationUserId, claim.Username, requestId, claim.OrganizationId, lenderOrganizationId); err != nil {
		l.svcCtx.Logger.WithContext(ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
	}
}

func (l *RegisterLogic) handlePreapproveLoan(ctx context.Context, creatorId int64, creatorName string, requestId string, borrowerOrganizationId int64, lenderOrganizationId int64) (err error) {
	// get pre-approve loan data from cache
	key := fmt.Sprintf(constants.BorrowerPreApprovalCacheKey, requestId)
	cacheDataStr, err := global.GetRedis().Get(ctx, key)
	if err != nil {
		return errorx.WithStack(err)
	}

	if len(cacheDataStr) == 0 {
		return errorx.WithStack(errorx.New("pre-approve loan data not found in cache"))
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
	lenderOrganization, err := l.svcCtx.TenantClient.Organization.DetailWithType(ctx, &tenant.OrganizationDetailWithTypeRequest{
		Data: &tenant.OrganizationDetailWithTypeRequestData{
			UserType: string(constants.UserTypeLender),
			Id:       lenderOrganizationId,
		}})
	if err != nil {
		return err
	}

	if lenderOrganization.Data == nil || lenderOrganization.Data.Status != string(constants.OrganizationStatusActive) {
		return errorx.WithStack(errorx.New("lender organization is not active"))
	}
	// prepare data
	var borrowerProperty *model.BorrowerProperty
	var borrowerApplicant *model.BorrowerApplicant
	err = l.svcCtx.BorrowerPropertyRepo.Transaction(ctx, func(tx *gorm.DB) error {
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
	loanPreapprovalInfo.Input, err = copierx.Trans[types.LoanPreApprovalRequest, borrower.LoanPreApprovalRequestData](&preApproveRequest)
	if err != nil {
		return err
	}
	loanPreapprovalInfo.Result, err = copierx.Trans[types.LoanPreApprovalResult, borrower.LoanPreapproveResult](&calcResult)
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

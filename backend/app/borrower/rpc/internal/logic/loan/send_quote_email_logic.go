package loanlogic

import (
	"context"
	"fmt"
	"strings"

	"github.com/duke-git/lancet/v2/slice"

	"zeebee-backend/app/borrower/rpc/internal/email_tpl"
	"zeebee-backend/app/borrower/rpc/internal/global"
	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	commonemail "zeebee-backend/common/email"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/email/client"
	"zeebee-backend/domains/zeebee/limiter"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/proto"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type SendQuoteEmailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewSendQuoteEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendQuoteEmailLogic {
	return &SendQuoteEmailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *SendQuoteEmailLogic) SendQuoteEmail(in *borrower.LoanSendQuoteEmailRequest) (resp *proto.EmptyResponse, err error) {

	if len(in.Data.LenderIds) == 0 {
		return nil, errorx.New("Lender is required")
	}

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	borrowerOrganization, err := l.svcCtx.OrganizationRepo.GetOrganization(l.ctx, claim.OrganizationId, constants.UserTypeBorrower, false)
	if err != nil {
		return nil, err
	}
	if borrowerOrganization.Status != constants.OrganizationStatusActive {
		return nil, errorx.New("Borrower organization is not active")
	}
	borrowerApplicantIds := slice.Map(in.Data.Applicants, func(index int, item *borrower.ThirdPartyLoanApplicantEditEntity) int64 {
		return item.BorrowerApplicantId
	})
	if len(borrowerApplicantIds) == 0 {
		return nil, errorx.New("Borrowers is required")
	}

	borrowerApplicantNames, err := l.svcCtx.BorrowerApplicantRepo.GetFullNames(l.ctx, borrowerApplicantIds, claim.OrganizationId)
	if err != nil {
		return nil, err
	}

	senderEmail, emailClient, err := l.getEmailClient(l.ctx, claim.UserId, claim.OrganizationId)
	if err != nil {
		return nil, err
	}
	borrowerProperty, err := l.svcCtx.BorrowerPropertyRepo.FindOneById(l.ctx, in.Data.LoanProperty.BorrowerPropertyId)
	if err != nil {
		return nil, err
	}

	go func() {
		newCtx := context.WithoutCancel(l.ctx)
		l.svcCtx.Logger.WithContext(newCtx).WithFields(logx.Field("lender ids", in.Data.LenderIds)).Info("Start send quote email")
		tplVars := map[string]interface{}{
			"borrower_name":           borrowerOrganization.Name,
			"borrower_applicant_name": strings.Join(borrowerApplicantNames, ","),
			"holdback_amount":         in.Data.HoldbackAmount,
			"closing_date":            in.Data.ClosingDate,
			"business_strategy":       strings.Join(in.Data.BusinessStrategy, ""),
			"total_loan_amount":       in.Data.TotalAmount,
			"loan_type":               in.Data.LoanType,
			"loan_purpose":            in.Data.LoanPurpose,
			"purchase_amount":         in.Data.LoanProperty.PurchaseAmount,
			"arv":                     in.Data.LoanProperty.Arv,
			"property":                borrowerProperty,
		}
		for _, id := range in.Data.LenderIds {
			input := quoteEmailInput{
				senderEmail: senderEmail,
				emailClient: emailClient,
				lenderId:    id,
				tplVars:     tplVars,
			}
			err := l.sendEmail(newCtx, input, claim)
			if err != nil {
				l.svcCtx.Logger.WithContext(newCtx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
			}
		}
		l.svcCtx.Logger.WithContext(newCtx).Info("End send quote email")
	}()
	l.svcCtx.Logger.WithContext(l.ctx).Info("Finish request")
	return &proto.EmptyResponse{}, nil
}

type quoteEmailInput struct {
	senderEmail string
	emailClient client.Client
	lenderId    int64
	tplVars     map[string]interface{}
}

// send email
func (l *SendQuoteEmailLogic) sendEmail(ctx context.Context, input quoteEmailInput, claim utils.JwtUser) (err error) {
	lender, err := l.svcCtx.OrganizationRepo.GetOrganization(ctx, input.lenderId, constants.UserTypeLender, false)
	if err != nil {
		return err
	}
	if lender.Status != constants.OrganizationStatusActive {
		return errorx.New("Lender organization is not active")
	}
	tplVars := input.tplVars
	tplVars["lender_name"] = lender.Name
	property, ok := tplVars["property"].(*model.BorrowerProperty)
	if !ok {
		return errorx.New("Invalid property")
	}

	key := global.EmailAuthManager.BuildKey(limiter.StrategySendLoanQuoteEmail, input.senderEmail, lender.ContactEmail, claim.OrganizationId, claim.UserType)
	// rate limit
	ok, err = global.LimiterMgr.Limit(ctx, limiter.StrategySendLoanQuoteEmail, key)
	if err != nil {
		return errorx.WithStack(err)
	}
	if !ok {
		return errorx.New("Rate limit")
	}

	emailContent := ""
	emailContent, err = commonemail.ParseTemplate(email_tpl.QuoteEmail, tplVars)
	if err != nil {
		return errorx.WithStack(err)
	}
	msg := commonemail.Message{
		Subject: fmt.Sprintf("%v lender - %v quote - %v", tplVars["lender_name"], tplVars["loan_purpose"], property.Address.Address1),
		From:    input.senderEmail,
		To:      []string{lender.ContactEmail},
		Body:    emailContent,
		IsHTML:  true,
	}
	err = input.emailClient.SendEmail(ctx, msg)
	if err != nil {
		return errorx.WithStack(err)
	}
	l.svcCtx.Logger.WithContext(ctx).WithFields(
		logx.Field("sender", input.senderEmail),
		logx.Field("receiver", lender.ContactEmail),
	).Info("Send quote email success")
	return err
}

// get borrower email
func (l *SendQuoteEmailLogic) getEmailClient(ctx context.Context, userId int64, organizationId int64) (senderEmail string, emailClient client.Client, err error) {
	var emailAuthorization *model.EmailAuthorization
	organizationEmail, err := l.svcCtx.OrganizationEmailRepo.GetOrganizationEmailByOrganizationId(ctx, organizationId, constants.EmailPurposeContact)
	if err != nil {
		return "", nil, err
	}
	switch organizationEmail.AuthStrategy {
	case constants.EmailAuthStrategyOrganization:
		emailAuthorization, err = l.svcCtx.EmailAuthorizationRepo.GetEmailAuthorizationByOrganizationEmailId(ctx, organizationEmail.Id, 0, constants.EmailPurposeContact)
	case constants.EmailAuthStrategyUser:
		emailAuthorization, err = l.svcCtx.EmailAuthorizationRepo.GetEmailAuthorizationByOrganizationEmailId(ctx, organizationEmail.Id, userId, constants.EmailPurposeContact)
	default:
		err = errorx.New("Invalid email auth strategy")
	}
	if err != nil {
		return "", nil, err
	}

	if !emailAuthorization.IsAuthorized {
		err = errorx.New("Borrower email is not authorized")
		return "", nil, err
	}

	var existed bool
	emailClient, existed = global.EmailAuthManager.GetClient(emailAuthorization.Email)
	if !existed {
		emailProvider, ok := global.EmailAuthUrlGenerator[emailAuthorization.EmailProvider]
		if !ok {
			return "", nil, errorx.New("Invalid email provider")
		}
		oauth2Config := emailProvider.GetConfig()
		switch emailAuthorization.EmailProvider {
		case constants.EmailProviderGmail:
			emailClient, err = client.NewGmailClientFromToken(oauth2Config, emailAuthorization.Token.Token)
		default:
			err = errorx.New("Invalid email provider")
		}
		if err != nil {
			return "", nil, err
		}
		global.EmailAuthManager.Register(emailAuthorization.Email, emailClient)
	}

	return emailAuthorization.Email, emailClient, nil
}

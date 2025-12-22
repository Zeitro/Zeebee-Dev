package loanlogic

import (
	"context"
	"strings"

	"zeebee-backend/app/borrower/rpc/internal/email_tpl"
	"zeebee-backend/app/borrower/rpc/internal/global"
	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	commonemail "zeebee-backend/common/email"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/email/client"
	"zeebee-backend/domains/zeebee/model"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type SendQuoteEmailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewSendQuoteEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendQuoteEmailLogic {
	return &SendQuoteEmailLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.NewBaseLogic(ctx, svcCtx),
	}
}

func (l *SendQuoteEmailLogic) SendQuoteEmail(in *borrower.LoanSendQuoteEmailRequest) (resp *borrower.EmptyResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	borrowerOrganization, err := l.svcCtx.OrganizationRepo.GetOrganization(l.ctx, claim.OrganizationId, constants.UserTypeBorrower, false)
	if err != nil {
		return nil, err
	}
	if borrowerOrganization.Status != constants.OrganizationStatusActive {
		return nil, errorx.New("Borrower Organization Not Active")
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
		tplVars := map[string]interface{}{
			"borrower_name":     borrowerOrganization.Name,
			"holdback_amount":   in.Data.HoldbackAmount,
			"closing_date":      in.Data.ClosingDate,
			"business_strategy": strings.Join(in.Data.BusinessStrategy, ""),
			"total_loan_amount": in.Data.TotalAmount,
			"loan_type":         in.Data.LoanType,
			"loan_purpose":      in.Data.LoanPurpose,
			"purchase_amount":   in.Data.LoanProperty.PurchaseAmount,
			"arv":               in.Data.LoanProperty.Arv,
			"property":          borrowerProperty,
		}
		for _, id := range in.Data.LenderIds {
			input := quoteEmailInput{
				senderEmail: senderEmail,
				emailClient: emailClient,
				lenderId:    id,
				tplVars:     tplVars,
			}
			err := l.sendEmail(newCtx, input)
			if err != nil {
				l.WithContext(newCtx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
			}

		}
	}()

	return &borrower.EmptyResponse{}, nil
}

type quoteEmailInput struct {
	senderEmail string
	emailClient client.Client
	lenderId    int64
	tplVars     map[string]interface{}
}

// send email
func (l *SendQuoteEmailLogic) sendEmail(ctx context.Context, input quoteEmailInput) (err error) {
	lender, err := l.svcCtx.OrganizationRepo.GetOrganization(ctx, input.lenderId, constants.UserTypeLender, false)
	if err != nil {
		return err
	}
	if lender.Status != constants.OrganizationStatusActive {
		return errorx.New("Lender Organization Not Active")
	}
	tplVars := input.tplVars
	tplVars["lender_name"] = lender.Name

	emailContent := ""
	emailContent, err = commonemail.ParseTemplate(email_tpl.QuoteEmail, tplVars)
	if err != nil {
		return errorx.WithStack(err)
	}
	msg := commonemail.Message{
		Subject: "Quoting",
		From:    input.senderEmail,
		To:      []string{lender.ContactEmail},
		Body:    emailContent,
		IsHTML:  true,
	}
	err = input.emailClient.SendEmail(ctx, msg)

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
		err = errorx.New("Invalid Email Auth Strategy")
	}
	if err != nil {
		return "", nil, err
	}

	if !emailAuthorization.IsAuthorized {
		err = errorx.New("Borrower Email Not Authorized")
		return "", nil, err
	}

	var existed bool
	emailClient, existed = global.EmailAuthManager.GetClient(emailAuthorization.Email)
	if !existed {
		emailProvider, ok := global.EmailAuthUrlGenerator[emailAuthorization.EmailProvider]
		if !ok {
			return "", nil, errorx.New("Invalid Email Provider")
		}
		oauth2Config := emailProvider.GetConfig()
		switch emailAuthorization.EmailProvider {
		case constants.EmailProviderGmail:
			emailClient, err = client.NewGmailClientFromToken(oauth2Config, emailAuthorization.Token.Token)
		default:
			err = errorx.New("Invalid Email Provider")

		}
		if err != nil {
			return "", nil, err
		}
		global.EmailAuthManager.Register(emailAuthorization.Email, emailClient)
	}

	return emailAuthorization.Email, emailClient, nil
}

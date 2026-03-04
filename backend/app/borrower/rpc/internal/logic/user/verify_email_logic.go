package userlogic

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/duke-git/lancet/v2/mathutil"
	"github.com/jinzhu/copier"

	"zeebee-backend/app/borrower/rpc/internal/global"
	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/email"
	"zeebee-backend/common/stringx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/templater"
	"zeebee-backend/domains/zeebee/templater/loader"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type VerifyEmailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewVerifyEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerifyEmailLogic {
	return &VerifyEmailLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *VerifyEmailLogic) VerifyEmail(in *borrower.VerifyEmailRequest) (resp *borrower.VerifyEmailResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	var user *model.User
	email := in.Data.Email
	code := in.Data.Code
	user, err = l.baseLogic.UserLogic.VerifyUserEmail(l.ctx, email, code, constants.UserTypeBorrower)
	if err != nil {
		return nil, err
	}
	organizationUser, err := l.svcCtx.OrganizationUserRepo.GetOrganizationUserByUserId(l.ctx, user.DefaultOrganizationId, user.Id)
	if err != nil {
		return nil, err
	}

	tokenInfo, err := l.baseLogic.UserLogic.GenerateJwtToken(user, organizationUser, constants.UserTypeBorrower)
	if err != nil {
		return nil, err
	}
	data := &borrower.UserLoginResponseData{}
	err = copier.CopyWithOption(data, tokenInfo, copier.Option{CaseSensitive: true})
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	go func() {
		newCtx := context.WithoutCancel(l.ctx)
		err := l.handlePreapprovalLoan(newCtx, user, organizationUser.OrganizationId)
		if err != nil {
			err = errorx.Wrap(err, "Failed to handle preapproval loan")
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	return &borrower.VerifyEmailResponse{
		Data: data,
	}, nil
}

func (l *VerifyEmailLogic) handlePreapprovalLoan(ctx context.Context, user *model.User, borrowerOrganizationId int64) (err error) {
	borrowerOrganization, err := l.svcCtx.OrganizationRepo.FindOneById(ctx, borrowerOrganizationId)
	if err != nil {
		return err
	}
	if borrowerOrganization == nil || borrowerOrganization.Id == 0 || borrowerOrganization.Status != constants.OrganizationStatusActive {
		l.WithContext(ctx).Info("Borrower organization is inactive")
		return nil
	}

	draftLoan, err := l.svcCtx.LoanRepo.FindOneByColumns(ctx, map[string]interface{}{
		"borrower_organization_id": borrowerOrganizationId,
		"stage":                    constants.LoanStageDraft,
	})
	if err != nil {
		return err
	}
	if draftLoan == nil || draftLoan.Id == 0 {
		l.WithContext(ctx).Info("Draft loan not found")
		return nil
	}
	lenderOrganizationId := draftLoan.LenderOrganizationId

	lenderOrganization, err := l.svcCtx.OrganizationRepo.FindOneById(ctx, lenderOrganizationId)
	if err != nil {
		return err
	}
	if lenderOrganization == nil || lenderOrganization.Id == 0 || lenderOrganization.Status != constants.OrganizationStatusActive {
		l.WithContext(ctx).Info("Lender organization is inactive")
		return nil
	}
	preapprovalLetterSetting, err := l.baseLogic.OrganizationSettingLogic.DetailByPurpose(ctx, lenderOrganization.Id, constants.OrganizationPurposePreapprovalLetter)
	if err != nil {
		return err
	}
	if !preapprovalLetterSetting.Enabled {
		l.WithContext(ctx).Info("Preapproval letter is not enabled")
		return nil
	}

	preapprovalInfo := draftLoan.Extend.LoanPreapprovalInfo
	if preapprovalInfo == nil {
		return nil
	}
	preapprovalInput := preapprovalInfo.Input
	if preapprovalInput == nil {
		return nil
	}

	preapprovalResult := preapprovalInfo.Result
	if preapprovalResult == nil {
		return nil
	}

	// Get PDF
	// now := time.Now()
	renderData := map[string]interface{}{
		// "pdfCreateTime": now.Month().String() + now.Format(" 02,2006"),
		"loanBorrower":  user.FirstName,
		"salesPrice":    stringx.FormatUsDollarCurrency(preapprovalInput.Loan.PurchaseAmount),
		"maxLoanAmount": stringx.FormatUsDollarCurrency(preapprovalResult.TotalAmount),
		"terms":         preapprovalResult.LoanTerm,
		"interestRate":  fmt.Sprintf("%.2f%%", mathutil.RoundToFloat(preapprovalResult.Rate*100.0, 2)),
	}

	pdfConfig, ok := global.Config.TemplateConf[constants.BorrowerPreapprovalLetterPdfName]
	if !ok {
		return errorx.New("Preapproval letter PDF template is not set")
	}
	pdfTemplateConfig := templater.TemplateConf{
		Purpose:          pdfConfig.Purpose,
		Scope:            templater.Scope(pdfConfig.Scope),
		Category:         templater.Category(pdfConfig.Category),
		OrganizationType: string(constants.UserTypeBorrower),
		OrganizationId:   lenderOrganizationId,
		FileType:         templater.FileTypePDF,
	}
	pdfContent, pdfContentType, err := global.Templater.Export(ctx, pdfTemplateConfig, renderData)
	if errors.Is(err, loader.ErrTemplateNotActive) {
		l.WithContext(ctx).Info("PDF template is inactive")
		return nil
	}
	if err != nil {
		return errorx.WithStack(err)
	}
	// Get email
	renderData = map[string]interface{}{
		"borrowerName": user.FirstName,
	}
	emailConfig, ok := global.Config.TemplateConf[constants.BorrowerPreapprovalLetterEmailName]
	if !ok {
		return errorx.New("Preapproval letter email template is not set")
	}
	emailTemplateConfig := templater.TemplateConf{
		Purpose:          emailConfig.Purpose,
		Scope:            templater.Scope(emailConfig.Scope),
		Category:         templater.Category(emailConfig.Category),
		OrganizationType: string(constants.UserTypeBorrower),
		OrganizationId:   lenderOrganizationId,
		FileType:         templater.FileTypeHTML,
	}
	emailContent, _, err := global.Templater.Export(ctx, emailTemplateConfig, renderData)
	if errors.Is(err, loader.ErrTemplateNotActive) {
		l.WithContext(ctx).Info("Email template is inactive")
		return nil
	}
	if err != nil {
		return errorx.WithStack(err)
	}
	emailMsg := email.Message{
		Subject: emailConfig.Title,
		From:    global.Config.EmailConf.Username,
		To:      []string{user.Email},
		Body:    string(emailContent),
		Cc:      nil,
		Attachments: []email.Attachment{
			{
				Filename: pdfConfig.Title,
				Reader:   bytes.NewReader(pdfContent),
				MimeType: pdfContentType,
			},
		},
		IsHTML: true,
	}

	err = global.EmailClient.SendEmail(emailMsg)
	if err != nil {
		return errorx.WithStack(err)
	}

	return err
}

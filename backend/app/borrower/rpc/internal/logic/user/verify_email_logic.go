package userlogic

import (
	"bytes"
	"context"
	"errors"
	"time"

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
		baseLogic: logic.NewBaseLogic(ctx, svcCtx),
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
			err = errorx.Wrap(err, "Failed To Handle Preapproval Loan")
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	return &borrower.VerifyEmailResponse{
		Data: data,
	}, nil
}

func (l *VerifyEmailLogic) handlePreapprovalLoan(ctx context.Context, user *model.User, borrowerOrganizationId int64) (err error) {
	borrowerOrganization, err := l.svcCtx.OrganizationRepo.FindOneById(ctx, borrowerOrganizationId)
	if borrowerOrganization == nil || borrowerOrganization.Id == 0 || borrowerOrganization.Status != constants.OrganizationStatusActive {
		l.WithContext(ctx).Info("Borrower Organization Is Not Active")
		return nil
	}

	draftLoan, _ := l.svcCtx.LoanRepo.FindOneByColumns(ctx, map[string]interface{}{
		"borrower_organization_id": borrowerOrganizationId,
		"stage":                    constants.LoanStageDraft,
	})
	if draftLoan == nil || draftLoan.Id == 0 {
		l.WithContext(ctx).Info("Draft Loan Is Not Found")
		return nil
	}
	lenderOrganizationId := draftLoan.LenderOrganizationId

	lenderOrganization, _ := l.svcCtx.OrganizationRepo.FindOneById(ctx, lenderOrganizationId)
	if lenderOrganization == nil || lenderOrganization.Id == 0 || lenderOrganization.Status != constants.OrganizationStatusActive {
		l.WithContext(ctx).Info("Lender Organization Is Not Active")
		return nil
	}
	if lenderOrganization.OwnerType != constants.OrganizationOwnerTypeSelf {
		l.WithContext(ctx).Info("Lender Organization Is Not Self Owned")
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

	// Get pdf
	now := time.Now()
	renderData := map[string]interface{}{
		"pdfCreateTime": now.Month().String() + now.Format(" 02,2006"),
		"loanBorrower":  user.FirstName,
		"salesPrice":    stringx.FormatUsDollarCurrency(preapprovalInput.Loan.PurchaseAmount),
		"maxLoanAmount": stringx.FormatUsDollarCurrency(preapprovalResult.TotalAmount),
		"terms":         preapprovalResult.LoanTerm,
		"interestRate":  mathutil.RoundToFloat(preapprovalResult.Rate*100.0, 2),
	}

	// for _, groupDto := range req.VariableGroupSchema {
	//
	// 	for _, variableDto := range groupDto.Variables {
	// 		if variableDto.Source == constants.TemplateVariableSourceFixed || variableDto.Source == constants.TemplateVariableSourceTemplate {
	// 			k := strings.Trim(variableDto.Key, ".")
	// 			data[k] = variableDto.FixedValue
	// 		}
	// 	}
	// }

	pdfConfig, ok := global.Config.TemplateConf[constants.BorrowerPreapprovalLetterPdfName]
	if !ok {
		return errorx.New("Preapproval Letter Pdf Template Is Not Set")
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
		l.WithContext(ctx).Info("Pdf Template Is Not Active")
		return nil
	}
	if err != nil {
		return errorx.WithStack(err)
	}
	// get email
	renderData = map[string]interface{}{
		"borrowerName": user.FirstName,
	}
	emailConfig, ok := global.Config.TemplateConf[constants.BorrowerPreapprovalLetterEmailName]
	if !ok {
		return errorx.New("Preapproval Letter Email Template Is Not Set")
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
		l.WithContext(ctx).Info("Email Template Is Not Active")
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

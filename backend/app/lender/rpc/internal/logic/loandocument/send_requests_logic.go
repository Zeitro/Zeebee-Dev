package loandocumentlogic

import (
	"context"

	"github.com/duke-git/lancet/v2/slice"
	"gorm.io/gorm"

	"zeebee-backend/app/lender/rpc/internal/email_tpl"
	"zeebee-backend/app/lender/rpc/internal/global"
	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/constants"
	email2 "zeebee-backend/common/email"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/model"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type SendRequestsLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
	logx.Logger
}

func NewSendRequestsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendRequestsLogic {
	return &SendRequestsLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.NewBaseLogic(ctx, svcCtx),
		Logger:    logx.WithContext(ctx),
	}
}

func (l *SendRequestsLogic) SendRequests(in *lender.LoanDocumentSendRequestsRequest) (resp *lender.EmptyResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	var loan *model.Loan

	var loanDocuments []*model.LoanDocument
	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	loan, err = l.prepare(in, claim)
	if err != nil {
		return nil, err
	}

	err = l.svcCtx.LoanDocumentRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		var innerErr error
		for _, documentType := range in.Data.Documents {
			var document *model.LoanDocument
			document, innerErr = l.svcCtx.LoanDocumentRepo.FindOneByColumnsWithTx(l.ctx, tx, map[string]interface{}{
				"id":      documentType.Id,
				"loan_id": loan.Id,
			})
			if innerErr != nil {
				return innerErr
			}
			if document.Status == constants.LoanDocumentStatusRequested {
				continue
			}
			if document.LoanId != loan.Id {
				return errorx.WithStack(errorx.New("Document Loan Id Not Match"))
			}

			loanDocuments = append(loanDocuments, document)

			innerErr = l.svcCtx.LoanDocumentRepo.UpdateColumnsByIdWithTx(l.ctx, tx, documentType.Id, map[string]interface{}{
				"status":  constants.LoanDocumentStatusRequested,
				"comment": documentType.Comment,
			})
			if innerErr != nil {
				return innerErr
			}
		}
		// add loan operation
		comments := slice.FlatMap[*lender.SendDocumentRequestEntry, string](in.Data.Documents, func(index int, item *lender.SendDocumentRequestEntry) []string {
			return []string{item.Comment}
		})
		comments = slice.Filter(comments, func(index int, item string) bool {
			return len(item) > 0
		})

		innerErr = l.baseLogic.CreateLoanOperationLog(
			tx,
			loan.Id,
			claim.OrganizationUserId,
			model.LoanOperationStatus{
				Stage:         loan.Stage,
				PendingStatus: loan.PendingStatus,
			},
			model.LoanOperationStatus{
				Stage:         loan.Stage,
				PendingStatus: loan.PendingStatus,
			},
			constants.LoanOperationActionSendDocumentRequest,
			&model.LoanOperationMetadata{
				Comment: slice.Join(comments, "/"),
			},
		)

		return innerErr
	})

	if len(loanDocuments) == 0 {
		return &lender.EmptyResponse{}, nil
	}

	// send email
	go func() {
		ctx := context.WithoutCancel(l.ctx)
		err := l.sendDocumentRequestEmail(ctx, loan, loanDocuments)
		if err != nil {
			err = errorx.Wrap(err, "Send Document Request Email Error")
			l.WithContext(ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	return &lender.EmptyResponse{}, nil
}

// prepare
func (l *SendRequestsLogic) prepare(in *lender.LoanDocumentSendRequestsRequest, claim utils.JwtUser) (loan *model.Loan, err error) {
	if len(in.Data.Documents) == 0 {
		return nil, errorx.WithStack(errorx.New("No Document Type Found"))
	}
	// var document *model.LoanDocument
	loan, _, err = l.baseLogic.GetLoanAndDocumentByDocumentId(in.Data.Documents[0].Id, claim.OrganizationId)
	return loan, err
}

func (l *SendRequestsLogic) sendDocumentRequestEmail(ctx context.Context, loan *model.Loan, loanDocuments []*model.LoanDocument) error {

	if loan.LoanBorrower == nil || len(loan.LoanBorrower.ContactEmail) == 0 {
		return errorx.New("No Contact Email Found For The Loan")
	}

	emailData := map[string]interface{}{
		"first_name":         loan.LoanBorrower.ContactName,
		"borrower_login_url": global.Config.BorrowerLoginUrl,
		"loan_documents":     loanDocuments,
	}

	// parse email template
	content, err := email2.ParseTemplate(email_tpl.SendLoanDocumentRequestEmail, emailData)
	if err != nil {
		return errorx.Wrap(err, "Failed To Parse Email Template")
	}

	// send email
	msg := email2.Message{
		Subject: "Loan Document Request",
		From:    global.Config.EmailConf.Username,
		To:      []string{loan.LoanBorrower.ContactEmail},
		Body:    content,
		IsHTML:  true,
	}

	err = global.EmailClient.SendEmail(msg)
	if err != nil {
		return errorx.Wrap(err, "Failed To Send Email")
	}

	return nil
}

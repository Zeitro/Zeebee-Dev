package loandocumentlogic

import (
	"context"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"zeebee-backend/app/lender/rpc/internal/email_tpl"
	"zeebee-backend/app/lender/rpc/internal/global"
	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/repo/model"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/constants"
	email2 "zeebee-backend/common/email"
	"zeebee-backend/common/types"

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

	loan, err = l.prepare(in)
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
				return errors.WithStack(errors.New("document loan id not match"))
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
			in.Data.LenderEmployeeId,
			types.LoanOperationStatus{
				Stage:         loan.Stage,
				PendingStatus: loan.PendingStatus,
			},
			types.LoanOperationStatus{
				Stage:         loan.Stage,
				PendingStatus: loan.PendingStatus,
			},
			constants.LoanOperationActionSendDocumentRequest,
			&types.LoanOperationMetadata{
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
			err = errors.Wrap(err, "send document request email error")
			l.WithContext(ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	return &lender.EmptyResponse{}, nil
}

// prepare
func (l *SendRequestsLogic) prepare(in *lender.LoanDocumentSendRequestsRequest) (loan *model.Loan, err error) {
	if len(in.Data.Documents) == 0 {
		return nil, errors.WithStack(errors.New("no document type found"))
	}
	// var document *model.LoanDocument
	loan, _, err = l.baseLogic.GetLoanAndDocumentByDocumentId(in.Data.Documents[0].Id, in.Data.LenderCompanyId)
	return loan, err
}

func (l *SendRequestsLogic) sendDocumentRequestEmail(ctx context.Context, loan *model.Loan, loanDocuments []*model.LoanDocument) error {

	if loan.LoanBorrower == nil || len(loan.LoanBorrower.ContactEmail) == 0 {
		return errors.New("no contact email found for the loan")
	}

	emailData := map[string]interface{}{
		"first_name":     loan.LoanBorrower.ContactName,
		"loan_documents": loanDocuments,
	}

	// parse email template
	content, err := email_tpl.GetParsedContent(email_tpl.SendLoanDocumentRequestEmail, emailData)
	if err != nil {
		return errors.Wrap(err, "failed to parse email template")
	}

	// send email
	msg := email2.EmailMessage{
		Subject: "Loan Document Request",
		From:    global.Config.EmailConf.Username,
		To:      []string{loan.LoanBorrower.ContactEmail},
		Body:    content,
		IsHTML:  true,
	}

	err = global.EmailClient.SendEmail(msg)
	if err != nil {
		return errors.Wrap(err, "failed to send email")
	}

	return nil
}

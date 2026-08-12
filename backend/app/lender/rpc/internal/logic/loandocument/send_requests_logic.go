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
	email2 "zeebee-backend/common/email"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/proto"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type SendRequestsLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewSendRequestsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendRequestsLogic {
	return &SendRequestsLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *SendRequestsLogic) SendRequests(in *lender.LoanDocumentSendRequestsRequest) (resp *proto.EmptyResponse, err error) {

	var loan *model.Loan

	var loanDocuments []*model.LoanDocument
	claim, err := l.baseLogic.GetUserClaimFromCtx(l.ctx)
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
				return errorx.WithStack(errorx.New("document loan id not match"))
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
		comments := slice.FlatMap[*lender.SendDocumentRequestEntity, string](in.Data.Documents, func(index int, item *lender.SendDocumentRequestEntity) []string {
			return []string{item.Comment}
		})
		comments = slice.Filter(comments, func(index int, item string) bool {
			return len(item) > 0
		})

		innerErr = l.baseLogic.CreateLoanOperationLog(
			l.ctx,
			tx,
			loan.Id,
			claim.OrganizationUserId,
			dto.LoanOperationStatus{
				Stage:         loan.Stage,
				PendingStatus: loan.PendingStatus,
			},
			dto.LoanOperationStatus{
				Stage:         loan.Stage,
				PendingStatus: loan.PendingStatus,
			},
			constants.LoanOperationActionSendDocumentRequest,
			&dto.LoanOperationMetadata{
				Comment: slice.Join(comments, "/"),
			},
		)

		return innerErr
	})

	if len(loanDocuments) == 0 {
		return &proto.EmptyResponse{}, nil
	}

	// send email
	go func() {
		ctx := context.WithoutCancel(l.ctx)
		err := l.sendDocumentRequestEmail(ctx, loan, loanDocuments)
		if err != nil {
			err = errorx.Wrap(err, "send document request email error")
			l.svcCtx.Logger.WithContext(ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
				logx.Field("service_name", l.svcCtx.Config.Name),
				logx.Field("type", "send_email"),
			).Error(err)
		}
	}()

	return &proto.EmptyResponse{}, nil
}

// prepare
func (l *SendRequestsLogic) prepare(in *lender.LoanDocumentSendRequestsRequest, claim utils.JwtUser) (loan *model.Loan, err error) {
	if len(in.Data.Documents) == 0 {
		return nil, errorx.WithStack(errorx.New("no document type found"))
	}
	// var document *model.LoanDocument
	loan, _, err = l.baseLogic.GetLoanAndDocumentByDocumentId(l.ctx, in.Data.Documents[0].Id, claim.OrganizationId)
	return loan, err
}

func (l *SendRequestsLogic) sendDocumentRequestEmail(ctx context.Context, loan *model.Loan, loanDocuments []*model.LoanDocument) error {
	if loan.LoanBorrower == nil || len(loan.LoanBorrower.ContactEmail) == 0 {
		return errorx.New("no contact email found for the loan")
	}

	emailData := map[string]any{
		"first_name":         loan.LoanBorrower.ContactName,
		"borrower_login_url": global.Config.BorrowerLoginUrl,
		"loan_documents":     loanDocuments,
	}

	// parse email template
	content, err := email2.ParseTemplate(email_tpl.SendLoanDocumentRequestEmail, emailData)
	if err != nil {
		return errorx.Wrap(err, "fail to parse email template")
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
		return errorx.Wrap(err, "fail to send email")
	}

	return nil
}

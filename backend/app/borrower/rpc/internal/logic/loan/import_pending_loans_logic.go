package loanlogic

import (
	"context"

	"github.com/duke-git/lancet/v2/slice"
	"gorm.io/gorm"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/stringx"
	"zeebee-backend/common/timex"
	"zeebee-backend/common/utils"
	"zeebee-backend/common/uuid"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
)

type importPendingLoanModelData struct {
	loan             *model.Loan
	loanBorrower     *model.LoanBorrower
	loanProperty     *model.LoanProperty
	borrowerProperty *model.BorrowerProperty
	applicants       []*model.LoanApplicant
	// loanExtend   *model.LoanExtend
}

type ImportPendingLoansLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewImportPendingLoansLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ImportPendingLoansLogic {
	return &ImportPendingLoansLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *ImportPendingLoansLogic) ImportPendingLoans(in *borrower.ImportPendingLoansRequest) (resp *borrower.ImportPendingLoansResponse, err error) {

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx, constants.UserTypeBorrower)
	if err != nil {
		return nil, err
	}

	if len(in.Data) == 0 {
		return &borrower.ImportPendingLoansResponse{
			Data: &borrower.ImportPendingLoansResponseData{
				RecordCount: 0,
			},
		}, nil
	}

	// build saved data
	models, err := l.buildImportModels(in.Data, claim)
	if err != nil {
		return nil, err
	}

	// create  loan
	err = l.svcCtx.LoanRepo.Transaction(l.ctx, func(tx *gorm.DB) (innerErr error) {

		var index int
		defer func() {
			if innerErr != nil {
				innerErr = errorx.Wrapf(innerErr, "fail to import pending loan at row:%d", index)
			}
		}()

		for i := range models {
			index = i + 1
			var borrowerProperty *model.BorrowerProperty
			var uid string
			data := models[i]
			// get BorrowerProperty
			borrowerProperty, innerErr = l.svcCtx.BorrowerPropertyRepo.GetOrCreateBorrowerPropertyWithTxForBorrower(l.ctx, tx, claim.OrganizationId, claim.OrganizationUserId, claim.Username, data.borrowerProperty.Address)
			if innerErr != nil {
				return innerErr
			}

			// Loan
			uid, innerErr = uuid.GetSonyflakeIDString()
			if innerErr != nil {
				return innerErr
			}
			loan := data.loan
			loan.Code = uid
			loan.RemoveRelationships()
			loan, innerErr = l.svcCtx.LoanRepo.CreateWithTx(l.ctx, tx, loan, false)
			if innerErr != nil {
				return innerErr
			}

			// LoanProperty
			loanProperty := data.loanProperty
			loanProperty.LoanId = loan.Id
			loanProperty.BorrowerPropertyId = borrowerProperty.Id
			loanProperty, innerErr = l.svcCtx.LoanPropertyRepo.CreateWithTx(l.ctx, tx, loanProperty, false)
			if innerErr != nil {
				return innerErr
			}

			// LoanBorrower
			loanBorrower := data.loanBorrower
			loanBorrower.LoanId = loan.Id
			loanBorrower, innerErr = l.svcCtx.LoanBorrowerRepo.CreateWithTx(l.ctx, tx, loanBorrower, false)
			if innerErr != nil {
				return innerErr
			}

			// LoanApplicants
			applicants := data.applicants
			slice.ForEach(applicants, func(index int, item *model.LoanApplicant) {
				item.LoanId = loan.Id
				item.LoanPropertyId = loanProperty.Id
				item.LoanBorrowerId = loanBorrower.Id
				item.FullName = stringx.FormatFullName(item.FirstName, item.MidName, item.LastName)
			})
			_, innerErr = l.svcCtx.LoanApplicantRepo.BatchCreateWithTx(l.ctx, tx, applicants, false)
			if innerErr != nil {
				return innerErr
			}

			// sync total amount and loan count
			innerErr = l.svcCtx.BorrowerPropertyRepo.SyncTotalAmountAndLoanCountWithTx(l.ctx, tx, borrowerProperty.Id)
			if innerErr != nil {
				return innerErr
			}
		}

		return innerErr
	})
	if err != nil {
		return nil, errorx.Wrap(err, "fail to import pending loan")
	}

	return &borrower.ImportPendingLoansResponse{
		Data: &borrower.ImportPendingLoansResponseData{
			RecordCount: int64(len(models)),
		},
	}, nil
}

// buildImportModels gets the models for importing a pending loan.
func (l *ImportPendingLoansLogic) buildImportModels(in []*borrower.ImportPendingLoansRequestData, claim utils.JwtUser) (data []*importPendingLoanModelData, err error) {
	data = make([]*importPendingLoanModelData, 0, len(in))
	for item := range in {
		item, err := l.buildImportModel(in[item], claim)
		if err != nil {
			return nil, err
		}
		data = append(data, item)
	}
	return data, nil
}

func (l *ImportPendingLoansLogic) buildImportModel(in *borrower.ImportPendingLoansRequestData, claim utils.JwtUser) (data *importPendingLoanModelData, err error) {
	// loan
	loan := &model.Loan{}
	loan, err = copierx.Trans[borrower.ImportPendingLoansRequestData, model.Loan](in, copierx.StrToSqlTimePtr(timex.LayoutUsStyle))
	if err != nil {
		return nil, err
	}
	loan.CreatorId = claim.OrganizationUserId
	loan.CreatorName = claim.Username
	loan.BorrowerOrganizationId = claim.OrganizationId
	loan.LoanType = constants.LoanType(in.LoanType)
	loan.Stage = constants.LoanStagePending
	loan.PendingStatus = constants.LoanPendingStatusQuoting
	loan.Source = constants.LoanSourceBorrower

	// loan borrower
	loanBorrower := &model.LoanBorrower{}
	loanBorrower, err = copierx.Trans[borrower.ImportPendingLoansRequestData, model.LoanBorrower](in, copierx.StrToSqlTimePtr(timex.LayoutUsStyle))
	if err != nil {
		return nil, err
	}

	// applicants
	var applicants []*model.LoanApplicant
	var totalOwnPercent int32
	applicantsPtr, err := copierx.Trans[[]*borrower.ImportPendingLoanApplicantEntity, []*model.LoanApplicant](&in.Applicants, copierx.StrToSqlTimePtr(timex.LayoutUsStyle))
	if err != nil {
		return nil, err
	}
	applicants = *applicantsPtr

	if len(applicants) == 0 {
		return nil, errorx.New("no applicant")
	}
	primaryCount := slice.CountBy(applicants, func(index int, item *model.LoanApplicant) bool {
		return item.IsPrimary
	})
	if primaryCount == 0 {
		return nil, errorx.New("no primary applicant")
	}
	slice.ForEach(applicants, func(index int, item *model.LoanApplicant) {
		item.BorrowerOrganizationId = claim.OrganizationId
		totalOwnPercent += item.OwnPercent
	})
	if totalOwnPercent > 100 {
		return nil, errorx.New("total own percent is greater than 100")
	}

	// loan property: address/owner/unit live in the nested loan_property entity, not the top-level request
	var loanProperty *model.LoanProperty
	loanProperty, err = copierx.Trans[borrower.ImportPendingLoanPropertyEntity, model.LoanProperty](in.LoanProperty, copierx.StrToSqlTimePtr(timex.LayoutUsStyle))
	if err != nil {
		return nil, err
	}
	if loanProperty == nil {
		loanProperty = &model.LoanProperty{}
	}

	// BorrowerProperty
	borrowerProperty := &model.BorrowerProperty{}
	borrowerProperty.CreatorId = claim.OrganizationUserId
	borrowerProperty.CreatorName = claim.Username
	borrowerProperty.Address = loanProperty.Address
	borrowerProperty.OrganizationId = claim.OrganizationId

	return &importPendingLoanModelData{
		loan:             loan,
		loanBorrower:     loanBorrower,
		loanProperty:     loanProperty,
		borrowerProperty: borrowerProperty,
		applicants:       applicants,
	}, nil
}

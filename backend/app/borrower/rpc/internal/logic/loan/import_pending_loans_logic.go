package loanlogic

import (
	"context"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/constants"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/stringx"
	"zeebee-backend/common/timex"
	"zeebee-backend/common/utils"
	"zeebee-backend/common/uuid"
	"zeebee-backend/domains/zeebee/model"

	"github.com/zeromicro/go-zero/core/logx"
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
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewImportPendingLoansLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ImportPendingLoansLogic {
	return &ImportPendingLoansLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.NewBaseLogic(ctx, svcCtx),
	}
}

func (l *ImportPendingLoansLogic) ImportPendingLoans(in *borrower.ImportPendingLoansRequest) (resp *borrower.ImportPendingLoansResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	// build saved data
	models, err := l.buildImportModels(in.Data, claim)
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
	creatorId := claim.OrganizationUserId
	creator, err := l.svcCtx.UserRepo.GetUser(l.ctx, creatorId, constants.UserTypeBorrower)
	if err != nil {
		return nil, err
	}

	// create  loan
	err = l.svcCtx.LoanRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		var innerErr error

		for i := range models {
			var borrowerProperty *model.BorrowerProperty
			var uid string
			data := models[i]
			// get BorrowerProperty
			borrowerProperty, innerErr = l.svcCtx.BorrowerPropertyRepo.GetOrCreateBorrowerPropertyWithTxForBorrower(l.ctx, tx, creator.Id, creator.FullName, data.borrowerProperty.Address)
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
		}

		return innerErr
	})
	if err != nil {
		return nil, errorx.Wrap(err, "Failed To Import Pending Loan")
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

	err = copier.CopyWithOption(loan, &in, copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
		Converters: []copier.TypeConverter{
			copierx.StrToSqlTimePtr(timex.LayoutUsStyle),
		},
	})
	if err != nil {
		return nil, errorx.GetFailedTransErr(err)
	}
	// creator, err := l.svcCtx.UserRepo.GetUser(l.ctx, claim.UserId, constants.UserTypeBorrower)
	// if err != nil {
	// 	return nil, err
	// }
	loan.CreatorId = claim.OrganizationUserId
	loan.CreatorName = claim.Username
	loan.BorrowerOrganizationId = claim.OrganizationId
	loan.LoanType = constants.LoanType(in.LoanType)
	loan.Stage = constants.LoanStagePending
	loan.PendingStatus = constants.LoanPendingStatusQuoting
	loan.Source = constants.LoanSourceBorrower

	// loan borrower
	loanBorrower := &model.LoanBorrower{}
	err = copier.CopyWithOption(loanBorrower, in.LoanBorrower, copier.Option{
		IgnoreEmpty:   true,
		CaseSensitive: true,
		DeepCopy:      true,
		Converters: []copier.TypeConverter{
			copierx.StrToSqlTimePtr(timex.LayoutUsStyle),
		},
	})
	if err != nil {
		return nil, errorx.GetFailedTransErr(err)
	}

	// applicants
	var applicants []*model.LoanApplicant
	var totalOwnPercent int32

	err = copier.CopyWithOption(&applicants, in.Applicants, copier.Option{
		IgnoreEmpty:   false,
		CaseSensitive: true,
		DeepCopy:      true,
		Converters: []copier.TypeConverter{
			copierx.StrToSqlTimePtr(timex.LayoutUsStyle),
		},
	})
	if err != nil {
		return nil, errorx.GetFailedTransErr(err)
	}
	if len(applicants) == 0 {
		return nil, errorx.New("No Applicant")
	}
	primaryCount := slice.CountBy(applicants, func(index int, item *model.LoanApplicant) bool {
		return item.IsPrimary
	})
	if primaryCount == 0 {
		return nil, errorx.New("No Primary Applicant")
	}
	slice.ForEach(applicants, func(index int, item *model.LoanApplicant) {
		item.BorrowerOrganizationId = claim.OrganizationId
		totalOwnPercent += item.OwnPercent
	})
	if totalOwnPercent > 100 {
		return nil, errorx.New("Total Own Percent Greater Than 100")
	}

	// loan property
	loanProperty := &model.LoanProperty{}
	err = copier.CopyWithOption(loanProperty, in.LoanProperty, copier.Option{
		IgnoreEmpty:   true,
		CaseSensitive: true,
		DeepCopy:      true,
		Converters: []copier.TypeConverter{
			copierx.StrToSqlTimePtr(timex.LayoutUsStyle),
		},
	})
	if err != nil {
		return nil, errorx.GetFailedTransErr(err)
	}
	// loanProperty.Unit = loanProperty.BorrowerProperty.Unit

	// BorrowerProperty
	borrowerProperty := &model.BorrowerProperty{}
	borrowerProperty.CreatorId = claim.OrganizationUserId
	borrowerProperty.CreatorName = claim.Username
	borrowerProperty.Address = loanProperty.Address

	return &importPendingLoanModelData{
		loan:             loan,
		loanBorrower:     loanBorrower,
		loanProperty:     loanProperty,
		borrowerProperty: borrowerProperty,
		applicants:       applicants,
	}, nil
}

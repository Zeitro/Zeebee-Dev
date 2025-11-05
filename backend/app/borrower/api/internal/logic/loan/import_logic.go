package loan

import (
	"bytes"
	"context"
	"strconv"
	"time"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/xuri/excelize/v2"
	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/constants"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/utils"
	"zeebee-backend/common/validator"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/common/timex"
)

const (
	pendingColPropertyAddress  = "Property Address"
	pendingColCity             = "City"
	pendingColState            = "State"
	pendingColZip              = "Zip"
	pendingColBorrower         = "Borrower"
	pendingColTargetCloseDate  = "Target Close Date"
	pendingColPurchasePrice    = "Purchase price"
	pendingColLoanMaturityDate = "Loan Maturity Date (1st)"
	pendingColEstClosingFunds  = "Est. Closing Funds"
	pendingColTargetLender     = "Target Lender"
	pendingColLoanType         = "Loan Type"
	pendingColEstimatedLoanAmt = "Estimated loan Amt"
	pendingColPricing          = "Pricing"
	pendingColLoanStatus       = "Loan Status Updates"
	pendingColMeetingNotes     = "Meeting Notes"
	pendingColActionList       = "Action List"
	pendingColExistingLender   = "Existing Lender (1st)"
	pendingColRehabLoan        = "Rehab Loan(1st）"
	pendingColOriginalDate     = "Original Contract Date"
	pendingColExtensionDate    = "Extension Date"
	pendingColLoanProcessor    = "Loan Processor (Person )"
	pendingColPurchaseDate     = "Purchase Date"
	pendingColOwner            = "Owner"
	pendingColBusinessLine     = "Business Line"
)

var pendingLoanFields = map[string]struct{}{
	pendingColPropertyAddress:  {},
	pendingColCity:             {},
	pendingColState:            {},
	pendingColZip:              {},
	pendingColBorrower:         {},
	pendingColTargetCloseDate:  {},
	pendingColPurchasePrice:    {},
	pendingColLoanMaturityDate: {},
	pendingColEstClosingFunds:  {},
	pendingColTargetLender:     {},
	pendingColLoanType:         {},
	pendingColPricing:          {},
	pendingColLoanStatus:       {},
	pendingColMeetingNotes:     {},
	pendingColActionList:       {},
	pendingColExistingLender:   {},
	pendingColRehabLoan:        {},
	pendingColOriginalDate:     {},
	pendingColExtensionDate:    {},
	pendingColLoanProcessor:    {},
	pendingColPurchaseDate:     {},
	pendingColOwner:            {},
	pendingColBusinessLine:     {},
	pendingColEstimatedLoanAmt: {},
}

type ImportLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewImportLogic import
func NewImportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ImportLogic {
	return &ImportLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ImportLogic) Import(fileContents []byte, stage string) (resp *types.ImportLoanResponse, err error) {
	var (
		borrowerId  int64
		recordCount int64
	)
	// get borrower id
	borrowerId, err = utils.GetAccountIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	switch constants.LoanStage(stage) {
	case constants.LoanStagePending:
		recordCount, err = l.importPendingLoans(fileContents, borrowerId)
	case constants.LoanStageActive:
		// todo import active loans
		fallthrough
	default:
		err = errors.WithStack(errors.Errorf("Invalid Loan Stage: %s", stage))
	}
	return &types.ImportLoanResponse{
		RecordCount: recordCount,
	}, err
}

func (l *ImportLogic) importPendingLoans(content []byte, borrowerId int64) (recordCount int64, err error) {
	r := bytes.NewReader(content)
	// create excelize
	f, err := excelize.OpenReader(r)
	if err != nil {
		return 0, errors.Wrap(err, "Failed To Open Excel File")
	}
	defer f.Close()

	// get rows
	rows, err := f.GetRows(f.GetSheetName(0))
	if err != nil {
		return 0, errors.Wrap(err, "Failed To Get Rows")
	}

	if len(rows) <= 1 {
		return 0, nil
	}
	// validate headers
	headerMap := make(map[string]struct{})
	for _, header := range rows[0] {
		headerMap[header] = struct{}{}
	}
	for field := range pendingLoanFields {
		if _, ok := headerMap[field]; !ok {
			return 0, errors.Errorf("Field %s Not Found In Header", field)
		}
	}
	// get loans
	importLoans, err := parsePendingLoans(rows, borrowerId)
	if err != nil {
		return 0, err
	}

	rpcReq := &borrower.ImportPendingLoansRequest{
		Data: []*borrower.ImportPendingLoansRequestData{},
	}
	err = copier.CopyWithOption(&rpcReq.Data, &importLoans, copier.Option{
		IgnoreEmpty:   false,
		CaseSensitive: true,
		DeepCopy:      true,
	})
	if err != nil {
		return 0, errorx.GetFailedTransErr(err)
	}
	var rpcResp *borrower.ImportPendingLoansResponse

	rpcResp, err = l.svcCtx.LoanRpcClient.ImportPendingLoans(l.ctx, rpcReq)
	if err != nil {
		return
	}

	return rpcResp.Data.RecordCount, err
}

func parsePendingLoans(rows [][]string, borrowerId int64) ([]*types.ImportPendingLoanEntry, error) {

	importLoans := make([]*types.ImportPendingLoanEntry, 0, len(rows))
	// Create a map for column indices
	columnIndex := make(map[string]int)
	for i, colName := range rows[0] {
		columnIndex[colName] = i
	}
	v, err := validator.GetEnValidator()
	if err != nil {
		return nil, errors.Wrap(err, "Failed To Get Validator")
	}

	for i := 1; i < len(rows); i++ {
		row := rows[i]
		if len(row) == 0 {
			continue
		}

		loanData := &types.ImportPendingLoanEntry{
			BorrowerId: borrowerId,
		}
		// property address
		addr := &types.AddressEntry{}
		if idx, ok := columnIndex[pendingColPropertyAddress]; ok && idx < len(row) {
			addr.Address1 = row[idx]
		}
		if idx, ok := columnIndex[pendingColCity]; ok && idx < len(row) {
			addr.City = row[idx]
		}
		if idx, ok := columnIndex[pendingColState]; ok && idx < len(row) {
			addr.State = row[idx]
		}
		if idx, ok := columnIndex[pendingColZip]; ok && idx < len(row) {
			addr.ZipCode = row[idx]
		}
		// handle LoanProperty
		loanData.LoanProperty = &types.ImportLoanApplyPropertyEntry{
			Unit:                 1,
			Bedrooms:             1,
			Bathrooms:            1,
			FutureUnit:           1,
			FutureLivingAreaSqft: 1,
			FutureBedrooms:       1,
			FutureBathrooms:      1,
			Address:              addr,
		}
		if idx, ok := columnIndex[pendingColOwner]; ok && idx < len(row) {
			loanData.LoanProperty.Owner = row[idx]
		}
		if idx, ok := columnIndex[pendingColPurchasePrice]; ok && idx < len(row) {
			amount, err := parseFloat(row[idx])
			if err != nil {
				return nil, errors.Errorf("Invalid Purchase Price At Row %d", i)
			}
			loanData.LoanProperty.PurchaseAmount = amount
		}
		if idx, ok := columnIndex[pendingColPurchaseDate]; ok && idx < len(row) {
			if !validateDate(row[idx]) {
				return nil, errors.Errorf("Invalid Date Format For Purchase Date At Row %d", i)
			}
			loanData.LoanProperty.PurchaseDate = row[idx]
		}

		// loan borrower
		loanBorrower := &types.ImportLoanBorrowerEntry{}
		if idx, ok := columnIndex[pendingColOwner]; ok && idx < len(row) {
			loanBorrower.Name = row[idx]
			// if strings.Contains(row[idx], "LLC") {
			// 	loanBorrower.Type = string(constants.ApplicantTypeLLC)
			// } else {
			// 	loanBorrower.Type = string(constants.ApplicantTypeIndividual)
			// }
			loanBorrower.Type = string(constants.ApplicantTypeLLC)
			loanBorrower.OwnPercent = 100
		}

		loanData.LoanBorrower = loanBorrower
		// applicant
		applicant := &types.ImportLoanApplicantEntry{}
		if idx, ok := columnIndex[pendingColBorrower]; ok && idx < len(row) {
			applicant.FirstName = row[idx]
			// applicant.Type = loanBorrower.Type
			applicant.IsPrimary = true
		}
		loanData.Applicants = append(loanData.Applicants, applicant)

		// extend
		extend := &types.ImportLoanExtendEntry{}
		if idx, ok := columnIndex[pendingColLoanStatus]; ok && idx < len(row) {
			extend.LoanStatusUpdates = row[idx]
		}
		if idx, ok := columnIndex[pendingColActionList]; ok && idx < len(row) {
			extend.ActionList = row[idx]
		}
		if idx, ok := columnIndex[pendingColBusinessLine]; ok && idx < len(row) {
			extend.BusinessLine = row[idx]
		}
		if idx, ok := columnIndex[pendingColMeetingNotes]; ok && idx < len(row) {
			extend.MeetingNotes = row[idx]
		}
		if idx, ok := columnIndex[pendingColPricing]; ok && idx < len(row) {
			extend.Pricing = row[idx]
		}
		loanData.Extend = extend

		// loan
		if idx, ok := columnIndex[pendingColTargetCloseDate]; ok && idx < len(row) {
			if !validateDate(row[idx]) {
				return nil, errors.Errorf("Invalid Date Format For Target Close Date At Row %d", i)
			}
			loanData.ClosingDate = row[idx]
		}

		if idx, ok := columnIndex[pendingColLoanMaturityDate]; ok && idx < len(row) {
			if !validateDate(row[idx]) {
				return nil, errors.Errorf("Invalid Date Format For Loan Maturity Date At Row %d", i)
			}
			loanData.MaturityDate = row[idx]
		}
		if idx, ok := columnIndex[pendingColEstClosingFunds]; ok && idx < len(row) {
			amount, err := parseFloat(row[idx])
			if err != nil {
				return nil, errors.Errorf("Invalid Closing Funds At Row %d", i)
			}
			loanData.ClosingFund = amount
		}
		if idx, ok := columnIndex[pendingColTargetLender]; ok && idx < len(row) {
			loanData.LenderCompany = row[idx]
		}
		if idx, ok := columnIndex[pendingColLoanType]; ok && idx < len(row) {
			loanData.LoanType = row[idx]
		}
		if idx, ok := columnIndex[pendingColEstimatedLoanAmt]; ok && idx < len(row) {
			amount, err := parseFloat(row[idx])
			if err != nil {
				return nil, errors.Errorf("Invalid Estimated Loan Amount At Row %d", i)
			}
			loanData.TotalAmount = amount
		}
		if idx, ok := columnIndex[pendingColRehabLoan]; ok && idx < len(row) {
			amount, err := parseFloat(row[idx])
			if err != nil {
				return nil, errors.Errorf("Invalid Rehab Loan Amount At Row %d", i)
			}
			loanData.HoldbackAmount = amount
		}

		if idx, ok := columnIndex[pendingColExistingLender]; ok && idx < len(row) {
			loanData.ExistingLenderCompany = row[idx]
		}
		if idx, ok := columnIndex[pendingColOriginalDate]; ok && idx < len(row) {
			loanData.OriginalContractDate = row[idx]
		}
		if idx, ok := columnIndex[pendingColExtensionDate]; ok && idx < len(row) {
			if !validateDate(row[idx]) {
				return nil, errors.Errorf("Invalid Date Format For Extension Date At Row %d", i)
			}
			loanData.ExtensionDate = row[idx]
		}
		if idx, ok := columnIndex[pendingColLoanProcessor]; ok && idx < len(row) {
			loanData.LoanProcessor = row[idx]
		}

		err = v.Validate(loanData)
		if err != nil {
			return nil, errors.Wrapf(err, "Invalid Data At Row %d", i)
		}
		importLoans = append(importLoans, loanData)
	}
	return importLoans, nil
}

func parseInteger(s string) (int64, error) {
	if len(s) == 0 {
		return 0, nil
	}
	return strconv.ParseInt(s, 10, 64)
}

func parseFloat(s string) (float64, error) {
	if len(s) == 0 {
		return 0, nil
	}
	return strconv.ParseFloat(s, 64)
}

func validateDate(s string) bool {
	if len(s) == 0 {
		return true
	}
	_, err := time.Parse(timex.LayoutUsStyle, s)
	return err == nil
}

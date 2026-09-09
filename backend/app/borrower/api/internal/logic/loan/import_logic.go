package loan

import (
	"bytes"
	"context"
	"strconv"
	"time"

	"github.com/xuri/excelize/v2"
	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/timex"
	"zeebee-backend/common/validator"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/proto"
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
		recordCount int64
	)
	// get borrower id
	// claim, err := utils.GetUserClaimFromCtx(l.ctx)
	// if err != nil {
	// 	return nil, err
	// }

	userUserRpcRsp, err := l.svcCtx.UserRpcClient.UserInfo(l.ctx, &proto.EmptyRequest{})
	if err != nil {
		return nil, err
	}
	if userUserRpcRsp.Data == nil {
		return nil, errorx.New("invalid user")
	}
	userId := userUserRpcRsp.Data.Id
	organizationId := userUserRpcRsp.Data.OrganizationId

	switch constants.LoanStage(stage) {
	case constants.LoanStagePending:
		recordCount, err = l.importPendingLoans(fileContents, userId, organizationId)
	case constants.LoanStageActive:
		// todo import active loans
		fallthrough
	default:
		err = errorx.WithStack(errorx.Errorf("invalid loan stage: %s", stage))
	}
	return &types.ImportLoanResponse{
		RecordCount: recordCount,
	}, err
}

func (l *ImportLogic) importPendingLoans(content []byte, userId int64, organizationId int64) (recordCount int64, err error) {
	r := bytes.NewReader(content)
	// create excelize
	f, err := excelize.OpenReader(r)
	if err != nil {
		return 0, errorx.Wrap(err, "fail to open excel file")
	}
	defer f.Close()

	// get rows
	rows, err := f.GetRows(f.GetSheetName(0))
	if err != nil {
		return 0, errorx.Wrap(err, "fail to get rows")
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
			return 0, errorx.Errorf("field %s not found in header", field)
		}
	}
	// get loans
	importLoans, err := parsePendingLoans(rows, userId, organizationId)
	if err != nil {
		return 0, err
	}

	rpcReq := &borrower.ImportPendingLoansRequest{
		Data: []*borrower.ImportPendingLoansRequestData{},
	}
	dataPtr, err := copierx.Trans[[]*types.ImportPendingLoanEntity, []*borrower.ImportPendingLoansRequestData](&importLoans)
	if err != nil {
		return 0, err
	}
	rpcReq.Data = *dataPtr

	var rpcResp *borrower.ImportPendingLoansResponse

	rpcResp, err = l.svcCtx.LoanRpcClient.ImportPendingLoans(l.ctx, rpcReq)
	if err != nil {
		return
	}

	return rpcResp.Data.RecordCount, err
}

func parsePendingLoans(rows [][]string, borrowerId int64, organizationId int64) ([]*types.ImportPendingLoanEntity, error) {
	importLoans := make([]*types.ImportPendingLoanEntity, 0, len(rows))
	// Create a map for column indices
	columnIndex := make(map[string]int)
	for i, colName := range rows[0] {
		columnIndex[colName] = i
	}
	v, err := validator.GetEnValidator()
	if err != nil {
		return nil, errorx.Wrap(err, "fail to get validator")
	}

	for i := 1; i < len(rows); i++ {
		row := rows[i]
		if len(row) == 0 {
			continue
		}

		loanData := &types.ImportPendingLoanEntity{
			UserId:         borrowerId,
			OrganizationId: organizationId,
		}
		// property address
		addr := &types.AddressEntity{
			Address1: cell(row, pendingColPropertyAddress, columnIndex),
			City:     cell(row, pendingColCity, columnIndex),
			State:    cell(row, pendingColState, columnIndex),
			ZipCode:  cell(row, pendingColZip, columnIndex),
		}
		// handle LoanProperty
		loanData.LoanProperty = &types.ImportLoanApplyPropertyEntity{
			Unit:                 1,
			Bedrooms:             1,
			Bathrooms:            1,
			FutureUnit:           1,
			FutureLivingAreaSqft: 1,
			FutureBedrooms:       1,
			FutureBathrooms:      1,
			Address:              addr,
			Owner:                cell(row, pendingColOwner, columnIndex),
		}
		if amount, err := cellAmount(row, pendingColPurchasePrice, columnIndex); err != nil {
			return nil, errorx.Errorf("invalid purchase price at row %d", i+1)
		} else {
			loanData.LoanProperty.PurchaseAmount = amount
		}
		if v, ok := cellDate(row, pendingColPurchaseDate, columnIndex); ok {
			loanData.LoanProperty.PurchaseDate = v
		} else {
			return nil, errorx.Errorf("invalid date format for purchase date at row %d", i+1)
		}

		// loan borrower
		loanBorrower := &types.ImportLoanBorrowerEntity{}
		if idx, ok := columnIndex[pendingColOwner]; ok && idx < len(row) {
			loanBorrower.Name = row[idx]
			// if strings.Contains(row[idx], "LLC") {
			// 	loanBorrower.Type = string(constants.ApplicantTypeLLC)
			// } else {
			// 	loanBorrower.Type = string(constants.ApplicantTypeIndividual)
			// }
			loanBorrower.Type = string(constants.ApplicantTypeLLC)
		}

		loanData.LoanBorrower = loanBorrower
		// applicant
		applicant := &types.ImportLoanApplicantEntity{}
		if idx, ok := columnIndex[pendingColBorrower]; ok && idx < len(row) {
			applicant.FirstName = row[idx]
			// applicant.Type = loanBorrower.Type
			applicant.IsPrimary = true
			applicant.OwnPercent = 100
		}
		loanData.Applicants = append(loanData.Applicants, applicant)

		// extend
		loanData.Extend = &types.ImportLoanExtendEntity{
			LoanStatusUpdates: cell(row, pendingColLoanStatus, columnIndex),
			ActionList:        cell(row, pendingColActionList, columnIndex),
			BusinessLine:      cell(row, pendingColBusinessLine, columnIndex),
			MeetingNotes:      cell(row, pendingColMeetingNotes, columnIndex),
			Pricing:           cell(row, pendingColPricing, columnIndex),
		}

		// loan
		if v, ok := cellDate(row, pendingColTargetCloseDate, columnIndex); ok {
			loanData.ClosingDate = v
		} else {
			return nil, errorx.Errorf("invalid date format for target close date at row %d", i+1)
		}

		if v, ok := cellDate(row, pendingColLoanMaturityDate, columnIndex); ok {
			loanData.MaturityDate = v
		} else {
			return nil, errorx.Errorf("invalid date format for loan maturity date at row %d", i+1)
		}
		if amount, err := cellAmount(row, pendingColEstClosingFunds, columnIndex); err != nil {
			return nil, errorx.Errorf("invalid closing funds at row %d", i+1)
		} else {
			loanData.ClosingFund = amount
		}
		loanData.LenderCompany = cell(row, pendingColTargetLender, columnIndex)
		loanData.LoanType = cell(row, pendingColLoanType, columnIndex)
		if amount, err := cellAmount(row, pendingColEstimatedLoanAmt, columnIndex); err != nil {
			return nil, errorx.Errorf("invalid estimated loan amount at row %d", i+1)
		} else {
			loanData.TotalAmount = amount
		}
		if amount, err := cellAmount(row, pendingColRehabLoan, columnIndex); err != nil {
			return nil, errorx.Errorf("invalid rehab loan amount at row %d", i+1)
		} else {
			loanData.HoldbackAmount = amount
		}

		loanData.ExistingLenderCompany = cell(row, pendingColExistingLender, columnIndex)
		loanData.OriginalContractDate = cell(row, pendingColOriginalDate, columnIndex)
		if v, ok := cellDate(row, pendingColExtensionDate, columnIndex); ok {
			loanData.ExtensionDate = v
		} else {
			return nil, errorx.Errorf("invalid date format for extension date at row %d", i+1)
		}
		loanData.LoanProcessor = cell(row, pendingColLoanProcessor, columnIndex)

		err = v.Validate(loanData)
		if err != nil {
			return nil, errorx.Wrapf(err, "invalid data at row %d", i+1)
		}
		importLoans = append(importLoans, loanData)
	}
	return importLoans, nil
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

// cell returns the value of col from row, or "" when the column is absent or
// the row has no cell at that index.
func cell(row []string, col string, cols map[string]int) string {
	if i, ok := cols[col]; ok && i < len(row) {
		return row[i]
	}
	return ""
}

// cellDate returns the date string of col and whether it is valid; empty is valid.
func cellDate(row []string, col string, cols map[string]int) (string, bool) {
	v := cell(row, col, cols)
	return v, v == "" || validateDate(v)
}

// cellAmount returns the parsed float of col; empty parses as 0.
func cellAmount(row []string, col string, cols map[string]int) (float64, error) {
	return parseFloat(cell(row, col, cols))
}

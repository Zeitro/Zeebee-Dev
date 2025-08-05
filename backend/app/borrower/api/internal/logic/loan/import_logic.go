package loan

import (
	"bytes"
	"context"
	"strconv"
	"strings"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/xuri/excelize/v2"
	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/client/loan"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/constants"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/utils"
	"zeebee-backend/common/validator"
)

const (
	// CSV列名常量定义
	ColBorrower              = "borrower"
	ColLoanPurpose           = "loan purpose"
	ColLoanType              = "loan type"
	ColTotalLoanAmount       = "total loan amount"
	ColInitialLoanAmount     = "initial loan amount"
	ColOutstandingLoanAmount = "outstanding loan amount"
	ColPrincipalBalance      = "principal balance"
	ColMonthlyPayment        = "monthly payment"
	ColPaidOffDate           = "paid off date"
	ColStage                 = "stage"
	ColPendingStatus         = "pending status"
	ColPropertyType          = "property type"
	ColNumberOfUnit          = "number of unit"
	ColState                 = "state"
	ColCity                  = "city"
	ColZipCode               = "zip code"
	ColStreetAddress1        = "street address1"
	ColStreetAddress2        = "street address2"
	ColLender                = "lender"
	ColClosingDate           = "closing date"
	ColClosingFund           = "closing fund"
	ColHoldback              = "holdback"
	ColRate                  = "rate"
	ColLienPosition          = "lien position"
	ColMaturityDate          = "maturity date"
	ColPrepaymentPenalty     = "prepayment penalty(month interest)"
	ColServiceCompany        = "service company"
	ColBusinessStrategy      = "business strategy"
)

type ImportLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// import
func NewImportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ImportLogic {
	return &ImportLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ImportLogic) Import(fileContents []byte) (resp *types.EmptyResponse, err error) {
	rpcReq := &loan.BatchCreateLoanRequest{
		Data: make([]*loan.LoanCreateRequestData, 0),
	}

	// get borrower id
	borrowerId, err := utils.GetAccountIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	r := bytes.NewReader(fileContents)
	// create excelize
	f, err := excelize.OpenReader(r)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open excel file")
	}
	defer f.Close()

	// get rows
	rows, err := f.GetRows(f.GetSheetName(0))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get rows")
	}

	//
	if len(rows) <= 1 {
		return &types.EmptyResponse{}, nil
	}

	// parse row header
	header := rows[0]
	columnIndex := make(map[string]int)
	for i, colName := range header {
		columnIndex[strings.ToLower(colName)] = i
	}

	loans := make([]*types.LoanCreateRequest, 0)

	// parse rows
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		if len(row) == 0 {
			continue
		}

		// create  loan
		loanData := &types.LoanCreateRequest{}

		// parse row data
		if idx, ok := columnIndex[ColBorrower]; ok && idx < len(row) {
			applicantNames := strings.Split(row[idx], ",")

			loanData.Applicants = make([]*types.LoanEditApplicantEntry, len(applicantNames))
			for j, applicantName := range applicantNames {
				loanData.Applicants[j] = &types.LoanEditApplicantEntry{
					Name: applicantName,
					Type: "Individual",
				}
				if j == 0 {
					loanData.Applicants[j].IsPrimary = true
				}
			}
		}
		if idx, ok := columnIndex[ColLoanPurpose]; ok && idx < len(row) {
			loanData.LoanPurpose = row[idx]
		}
		if idx, ok := columnIndex[ColLoanType]; ok && idx < len(row) {
			loanData.LoanType = row[idx]
		}
		if idx, ok := columnIndex[ColTotalLoanAmount]; ok && idx < len(row) {
			if val := row[idx]; val != "" {
				if amount, err := strconv.ParseFloat(val, 64); err == nil {
					loanData.TotalAmount = amount
				}
			}
		}
		if idx, ok := columnIndex[ColInitialLoanAmount]; ok && idx < len(row) {
			if val := row[idx]; val != "" {
				if amount, err := strconv.ParseFloat(val, 64); err == nil {
					loanData.InitialAmount = amount
				}
			}
		}
		if idx, ok := columnIndex[ColStage]; ok && idx < len(row) {
			loanData.Stage = row[idx]
		}
		if idx, ok := columnIndex[ColPendingStatus]; ok && idx < len(row) {
			loanData.PendingStatus = row[idx]
		}
		if idx, ok := columnIndex[ColLender]; ok && idx < len(row) {
			loanData.LenderCompany = row[idx]
		}
		if idx, ok := columnIndex[ColClosingDate]; ok && idx < len(row) {
			loanData.ClosingDate = row[idx]
		}
		if idx, ok := columnIndex[ColClosingFund]; ok && idx < len(row) {
			if val := row[idx]; val != "" {
				if amount, err := strconv.ParseFloat(val, 64); err == nil {
					loanData.ClosingFund = amount
				}
			}
		}
		if idx, ok := columnIndex[ColHoldback]; ok && idx < len(row) {
			if val := row[idx]; val != "" {
				if amount, err := strconv.ParseFloat(val, 64); err == nil {
					loanData.HoldbackAmount = amount
				}
			}
		}
		if idx, ok := columnIndex[ColRate]; ok && idx < len(row) {
			if val := row[idx]; val != "" {
				if rate, err := strconv.ParseFloat(val, 32); err == nil {
					loanData.Rate = float32(rate)
				}
			}
		}
		if idx, ok := columnIndex[ColMonthlyPayment]; ok && idx < len(row) {
			if val := row[idx]; val != "" {
				if amount, err := strconv.ParseFloat(val, 64); err == nil {
					loanData.MonthlyPayment = amount
				}
			}
		}
		if idx, ok := columnIndex[ColPaidOffDate]; ok && idx < len(row) {
			loanData.PaidOffDate = row[idx]
		}
		if idx, ok := columnIndex[ColLienPosition]; ok && idx < len(row) {
			loanData.LienPosition = row[idx]
		}
		if idx, ok := columnIndex[ColMaturityDate]; ok && idx < len(row) {
			loanData.MaturityDate = row[idx]
		}
		if idx, ok := columnIndex[ColPrepaymentPenalty]; ok && idx < len(row) {
			if val := row[idx]; val != "" {
				if penalty, err := strconv.ParseInt(val, 10, 32); err == nil {
					loanData.PrepaymentPenalty = int32(penalty)
				}
			}
		}
		if idx, ok := columnIndex[ColServiceCompany]; ok && idx < len(row) {
			loanData.ServiceCompany = row[idx]
		}

		if idx, ok := columnIndex[ColBusinessStrategy]; ok && idx < len(row) {
			loanData.BusinessStrategy = strings.Split(row[idx], ",")
		}

		stage := ""
		if idx, ok := columnIndex[ColStage]; ok && idx < len(row) {
			stage = row[idx]
		}

		// when stage is archived, use principal_balance
		// when stage is active, use outstanding_loan_amount
		if strings.ToLower(stage) == strings.ToLower(string(constants.LoanStageArchived)) {
			if idx, ok := columnIndex[ColPrincipalBalance]; ok && idx < len(row) {
				if val := row[idx]; val != "" {
					if amount, err := strconv.ParseFloat(val, 64); err == nil {
						loanData.PrincipalAmount = amount
					}
				}
			}
		} else if strings.ToLower(stage) == strings.ToLower(string(constants.LoanStageActive)) {
			if idx, ok := columnIndex[ColOutstandingLoanAmount]; ok && idx < len(row) {
				if val := row[idx]; val != "" {
					if amount, err := strconv.ParseFloat(val, 64); err == nil {
						loanData.PrincipalAmount = amount
					}
				}
			}
		}

		// parse property
		property := &types.PropertyEntry{}
		hasPropertyData := false

		if idx, ok := columnIndex[ColPropertyType]; ok && idx < len(row) {
			property.PropertyType = row[idx]
			if property.PropertyType != "" {
				hasPropertyData = true
			}
		}
		if idx, ok := columnIndex[ColNumberOfUnit]; ok && idx < len(row) {
			if val := row[idx]; val != "" {
				if unit, err := strconv.ParseInt(val, 10, 32); err == nil {
					property.Unit = int32(unit)
					hasPropertyData = true
				}
			}
		}

		// parse  address
		address := &types.AddressEntry{}
		hasAddressData := false

		if idx, ok := columnIndex[ColState]; ok && idx < len(row) {
			address.State = row[idx]
			if address.State != "" {
				hasAddressData = true
			}
		}
		if idx, ok := columnIndex[ColCity]; ok && idx < len(row) {
			address.City = row[idx]
			if address.City != "" {
				hasAddressData = true
			}
		}
		if idx, ok := columnIndex[ColZipCode]; ok && idx < len(row) {
			address.ZipCode = row[idx]
			if address.ZipCode != "" {
				hasAddressData = true
			}
		}
		if idx, ok := columnIndex[ColStreetAddress1]; ok && idx < len(row) {
			address.Address1 = row[idx]
			if address.Address1 != "" {
				hasAddressData = true
			}
		}
		if idx, ok := columnIndex[ColStreetAddress2]; ok && idx < len(row) {
			address.Address2 = row[idx]
			// street address2 can be empty
		}

		if hasAddressData {
			property.Address = address
			hasPropertyData = true
		}

		if hasPropertyData {
			loanData.Property = property
		}
		loans = append(loans, loanData)
	}
	if len(loans) == 0 {
		return nil, nil
	}

	v, err := validator.GetEnValidator()
	if err != nil {
		return nil, err
	}

	// validate
	for i, data := range loans {
		err = v.Validate(data)
		if err != nil {
			err = errors.Errorf("Invalid argument at row %d: %s", i+1, err.Error())
			return nil, err
		}
		// trans
		rpcCreateLoan := &borrower.LoanCreateRequestData{}
		err = copier.Copy(rpcCreateLoan, data)
		if err != nil {
			return nil, errorx.GetFailedTransReqErr(err)
		}
		rpcCreateLoan.BorrowerId = borrowerId
		rpcReq.Data = append(rpcReq.Data, rpcCreateLoan)
	}

	_, err = l.svcCtx.LoanRpcClient.BatchCreateBorrowerLoan(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	return &types.EmptyResponse{}, nil
}

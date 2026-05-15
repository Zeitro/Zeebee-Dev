package conforming_loan_limit

import (
	"bytes"
	"strings"

	"github.com/duke-git/lancet/v2/convertor"
	"github.com/xuri/excelize/v2"

	"zeebee-backend/common/errorx"
	"zeebee-backend/common/stringx"
)

const (
	filedStateCode    = "FIPS State Code"
	filedCountyCode   = "FIPS County Code"
	filedCountyName   = "County Name"
	filedState        = "State"
	filedCbsaNum      = "CBSA Number"
	filedOneUnit      = "One-Unit\nLimit"
	filedTwoUnit      = "Two-Unit\nLimit"
	filedThreeUnit    = "Three-Unit\nLimit"
	filedFourPlusUnit = "Four-Unit\nLimit"
)

var loanLimitFields = map[string]struct{}{
	filedStateCode:    {},
	filedCountyCode:   {},
	filedCountyName:   {},
	filedState:        {},
	filedCbsaNum:      {},
	filedOneUnit:      {},
	filedTwoUnit:      {},
	filedThreeUnit:    {},
	filedFourPlusUnit: {},
}

type ConformingLoanLimit struct {
	StateCode         string
	CountyCode        string
	CountyName        string
	State             string
	CbsaNum           string
	OneUnitLimit      int64
	TwoUnitLimit      int64
	ThreeUnitLimit    int64
	FourPlusUnitLimit int64
}

type StateHighCostLimitMgr struct {
	StateHighCostLimits map[string]*StateHighCostLimit
}

type StateHighCostLimit struct {
	State        string
	CountyLimits map[string]*ConformingLoanLimit
}

// IsHighCostZone return true if the county is high cost zone
func (mgr *StateHighCostLimitMgr) IsHighCostZone(state string, county string) (isHighCostZone bool) {
	stateHighCostLimit, isExisted := mgr.StateHighCostLimits[state]
	if !isExisted {
		return false
	}
	_, isExisted = stateHighCostLimit.CountyLimits[county]

	return isExisted
}

func ParseHighCostLimit(content []byte, limitAmount int64) (mgr *StateHighCostLimitMgr, err error) {
	limits, err := getConformingLoanLimits(content)
	if err != nil {
		return
	}
	mgr = &StateHighCostLimitMgr{
		StateHighCostLimits: make(map[string]*StateHighCostLimit),
	}
	for index := range limits {
		limit := limits[index]
		// fmt.Printf("[%d]=>%v\n",index,limit)
		stateHighCostLimit, ok := mgr.StateHighCostLimits[limit.State]
		if !ok {
			stateHighCostLimit = &StateHighCostLimit{
				State:        limit.State,
				CountyLimits: make(map[string]*ConformingLoanLimit),
			}
			mgr.StateHighCostLimits[limit.State] = stateHighCostLimit
		}
		if limit.OneUnitLimit > limitAmount {
			stateHighCostLimit.CountyLimits[limit.CountyName] = limit
		}
	}

	return
}

func getConformingLoanLimits(content []byte) (limits []*ConformingLoanLimit, err error) {
	r := bytes.NewReader(content)
	f, err := excelize.OpenReader(r)

	// f, err := excelize.OpenFile(filePath)
	if err != nil {
		err = errorx.Wrap(err, "Failed To Read Excel Content")
		return
	}

	defer f.Close()

	rows, err := f.GetRows(f.GetSheetName(0))
	if err != nil {
		err = errorx.Wrap(err, "Failed To Get Rows")
		return
	}
	if len(rows) < 3 {
		err = errorx.New("No Data")
		return
	}
	// check header fileds
	headerMap := make(map[string]struct{})
	for _, header := range rows[1] {
		headerMap[header] = struct{}{}
	}
	// fmt.Println(headerMap)
	for field := range loanLimitFields {
		if _, ok := headerMap[field]; !ok {
			return nil, errorx.Errorf("Field %s Not Found In Header", field)
		}
	}

	columnIndex := make(map[string]int)
	for i, colName := range rows[1] {
		columnIndex[colName] = i
	}
	// get and parse rows
	limits = make([]*ConformingLoanLimit, 0, len(rows)-2)
	for i := 2; i < len(rows); i++ {
		row := rows[i]
		limit := &ConformingLoanLimit{
			State:      trim(row[columnIndex[filedState]]),
			CountyName: trim(row[columnIndex[filedCountyName]]),
			CountyCode: trim(row[columnIndex[filedCountyCode]]),
			StateCode:  trim(row[columnIndex[filedStateCode]]),
			CbsaNum:    trim(row[columnIndex[filedCbsaNum]]),
		}
		limit.CountyName = handleCountyName(limit.CountyName)

		limit.OneUnitLimit, err = parseUsaCurrency(row[columnIndex[filedOneUnit]])
		if err != nil {
			return
		}
		limit.TwoUnitLimit, err = parseUsaCurrency(row[columnIndex[filedTwoUnit]])
		if err != nil {
			return
		}
		limit.ThreeUnitLimit, err = parseUsaCurrency(row[columnIndex[filedThreeUnit]])
		if err != nil {
			return
		}
		limit.FourPlusUnitLimit, err = parseUsaCurrency(row[columnIndex[filedFourPlusUnit]])

		limits = append(limits, limit)
	}

	return
}

// handleCountyName
func handleCountyName(s string) (result string) {
	suffixes := []string{"COUNTY", "EAST BOROUGH", "CITY AND BOROUGH", "BOROUGH", "WEST CENSUS AREA", "CENSUS AREA"}
	result = s
	for _, suffix := range suffixes {
		if strings.HasSuffix(result, suffix) {
			result = strings.TrimSuffix(s, suffix)
			break
		}
	}
	result = trim(result)
	result = stringx.ToTitle(result)
	return
}

// trim
func trim(s string) string {
	s = strings.TrimPrefix(s, "'")
	return strings.TrimSpace(s)
}

// parseUsaCurrency
func parseUsaCurrency(s string) (result int64, err error) {
	// $8500,000 => 8500000
	s = trim(s)
	s = strings.ReplaceAll(s, "$", "")
	s = strings.ReplaceAll(s, ",", "")
	if s == "" {
		return 0, nil
	}
	val, err := convertor.ToInt(s)
	if err != nil {
		return 0, errorx.Wrapf(err, "failed to parse int %s", s)
	}
	return val, nil
}

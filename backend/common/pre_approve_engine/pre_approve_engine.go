package pre_approve_engine

import (
	"math"

	"zeebee-backend/common/constants"
)

type ExperienceLevel int

const (
	ExperienceLevelNoExperienced ExperienceLevel = iota
	ExperienceLevelExperienced
	ExperienceLevelInstitutional
)

type InputLoanInfo struct {
	LoanType                constants.LoanType
	LoanPurpose             constants.LoanPurpose
	FicoScore               int32
	FlipPropertyCount       int32
	PurchaseRentalCount     int32
	UnitCount               int32
	LoanTerm                int32
	EstCurrentPropertyValue float64
	EstAfterRepairValue     float64
	OriginalBudget          float64
	PurchasePrice           float64
	OutstandingLoanAmount   float64
	IsHighValueArea         bool
	State                   string
	County                  string
	HasAssignmentFee        bool
	AssignmentFeeAmount     float64
	InterestReserveMonths   int32
}

type CalculatePreApprovedLoanAmountResult struct {
	IsApproved      bool
	TotalAmount     float64
	Rate            float32
	MonthlyPayment  float64
	CashToClose     float64
	InitialAmount   float64
	HoldbackAmount  float64
	LoanTerm        int32
	Ltv             float32
	Tlc             float32
	Ltarv           float32
	ExperienceLevel ExperienceLevel
}

type HighValueState struct {
	State  string
	Cities []string
}

type LoanTypeRule struct {
	LoanType       constants.LoanType // loan type
	GeneralRules   Rules              // general rules
	HighValueRules Rules              // high value rules
}

type BaseRate struct {
	ExperienceLevelNoExperienced float32
	ExperienceLevelExperienced   float32
	ExperienceLevelInstitutional float32
}

type PreApproveAmountRange struct {
	StandardAmountRange  AmountRange
	LargeUnitAmountRange AmountRange
}

type AmountRange struct {
	Min float64
	Max float64
}

type ExperienceLevelPriceCap struct {
	ExperienceLevel       ExperienceLevel
	ExtraFactor           float32                              // extra factor for calculating ltc/ltaiv/ltarv
	PriceCapRuleMap       map[constants.LoanType]*LoanTypeRule // loan type rule map
	HighValueStates       []*HighValueState                    // high value states
	BaseRate              BaseRate                             // base rate
	PreApproveAmountRange PreApproveAmountRange                // pre_approved loan amount range
	MinLargeUnitCount     int32                                // if unit count >= minLargeUnitCount, use large unit logic,or use standard logic
}

// CalPreApprovedLoanAmount returns the pre-approved loan amount
func (e *ExperienceLevelPriceCap) CalPreApprovedLoanAmount(input InputLoanInfo) (result *CalculatePreApprovedLoanAmountResult, err error) {
	result = &CalculatePreApprovedLoanAmountResult{
		IsApproved: false,
	}

	// get experience level
	result.ExperienceLevel = e.calculateExperienceLevel(input.FlipPropertyCount, input.PurchaseRentalCount)

	// 检查拒绝贷款的限制条件
	if e.isLoanRejected(input, result.ExperienceLevel) {
		return result, nil
	}

	// 计算三个贷款上限值
	maxByLTAIV := input.EstCurrentPropertyValue * float64(e.getLTAIVCap(result.ExperienceLevel, input.UnitCount, input.IsHighValueArea))
	maxByLTC := (input.EstCurrentPropertyValue + input.OriginalBudget) * float64(e.getLTCCap(result.ExperienceLevel, input.UnitCount, input.IsHighValueArea))
	maxByLTARV := input.EstAfterRepairValue * float64(e.getLTARVCap(result.ExperienceLevel, input.UnitCount, input.IsHighValueArea))

	// 取三个值中的最小值
	totalLoanAmount := math.Min(math.Min(maxByLTAIV, maxByLTC), maxByLTARV)

	// 应用全局最小/最大值限制
	// if totalLoanAmount < e.MinLoanAmount {
	// 	totalLoanAmount = e.MinLoanAmount
	// }
	//
	// // 对于1-4单元物业，最大值为MaxLoanAmount1To4Unit
	// if input.UnitCount >= 1 && input.UnitCount <= 4 {
	// 	if totalLoanAmount > e.MaxLoanAmount1To4Unit {
	// 		totalLoanAmount = e.MaxLoanAmount1To4Unit
	// 	}
	// } else if totalLoanAmount > e.MaxLoanAmount {
	// 	totalLoanAmount = e.MaxLoanAmount
	// }

	// GUC规则：如果贷款类型为GUC，需要额外加上Interest reserve amount
	if input.LoanType == constants.LoanTypeConstruction {
		interestReserveAmount := e.calculateInterestReserveAmount(totalLoanAmount, result.ExperienceLevel, input.InterestReserveMonths)
		totalLoanAmount += interestReserveAmount
		result.HoldbackAmount = interestReserveAmount
	}

	// 设置计算结果
	result.TotalAmount = totalLoanAmount
	result.LoanTerm = input.LoanTerm

	// 计算LTV、TLC、LTARV比率
	if input.EstCurrentPropertyValue > 0 {
		result.Ltv = float32(totalLoanAmount / input.EstCurrentPropertyValue)
		result.Tlc = float32((input.EstCurrentPropertyValue + input.OriginalBudget) / input.EstCurrentPropertyValue)
		result.Ltarv = float32(totalLoanAmount / input.EstAfterRepairValue)
	}

	// 计算利率
	result.Rate = e.calculateRate(input, result)

	// 计算月供
	result.MonthlyPayment = e.calculateMonthlyPayment(result.TotalAmount, float64(result.Rate), float64(input.LoanTerm))

	// 计算现金结算
	result.CashToClose = e.calculateCashToClose(input, result)

	// 计算初始贷款和修复贷款
	result.InitialAmount, result.HoldbackAmount = e.calculateInitialAndHoldbackAmount(input, result)

	result.IsApproved = true
	return result, nil
}

// calculateExperienceLevel 计算经验等级
func (e *ExperienceLevelPriceCap) calculateExperienceLevel(flipCount, purchaseRentalCount int32) ExperienceLevel {
	experienceValue := flipCount + purchaseRentalCount
	if experienceValue >= 10 {
		return ExperienceLevelInstitutional
	} else if experienceValue >= 3 {
		return ExperienceLevelExperienced
	}
	return ExperienceLevelNoExperienced
}

// isLoanRejected 检查是否应该拒绝贷款
func (e *ExperienceLevelPriceCap) isLoanRejected(input InputLoanInfo, experienceLevel ExperienceLevel) bool {
	// 检查单元数最大值
	if input.UnitCount > 30 {
		return true
	}

	// 检查original budget不能大于现在的房屋价值的100%
	if input.OriginalBudget > input.EstCurrentPropertyValue {
		return true
	}

	// 检查贷款期限范围（12-24个月）
	if input.LoanTerm < 12 || input.LoanTerm > 24 {
		return true
	}

	// 检查borrower equity是否大于等于10%的total cost
	borrowerEquity := e.calculateBorrowerEquity(input)
	totalCost := e.calculateTotalCost(input)
	if borrowerEquity < totalCost*0.1 {
		return true
	}

	// 检查特殊情况
	// 当loan type为ground up Construction，并且loan purpose为cash out refi
	if input.LoanType == constants.LoanTypeConstruction && input.LoanPurpose == constants.LoanPurposeCashOut {
		// lender端工作人员会视情况做决定，这里我们暂时不允许
		return true
	}

	// 当loan type为rehab，并且loan purpose为cash out refi
	if input.LoanType == constants.LoanTypeRehab && input.LoanPurpose == constants.LoanPurposeCashOut {
		// 根据文档，这种情况应该被拒绝
		return true
	}

	return false
}

// calculateBorrowerEquity 计算borrower equity
func (e *ExperienceLevelPriceCap) calculateBorrowerEquity(input InputLoanInfo) float64 {
	if input.LoanPurpose == constants.LoanPurposePurchase {
		return input.PurchasePrice // purchase时等于purchase price
	}
	// refi时等于est current value - total loan amount
	// 注意：由于我们正在计算total loan amount，所以这里使用预估值
	maxByLTAIV := input.EstCurrentPropertyValue * float64(e.getLTAIVCap(e.calculateExperienceLevel(input.FlipPropertyCount, input.PurchaseRentalCount), input.UnitCount, input.IsHighValueArea))
	maxByLTC := (input.EstCurrentPropertyValue + input.OriginalBudget) * float64(e.getLTCCap(e.calculateExperienceLevel(input.FlipPropertyCount, input.PurchaseRentalCount), input.UnitCount, input.IsHighValueArea))
	maxByLTARV := input.EstAfterRepairValue * float64(e.getLTARVCap(e.calculateExperienceLevel(input.FlipPropertyCount, input.PurchaseRentalCount), input.UnitCount, input.IsHighValueArea))
	totalLoanAmount := math.Min(math.Min(maxByLTAIV, maxByLTC), maxByLTARV)
	return input.EstCurrentPropertyValue - totalLoanAmount
}

// calculateTotalCost 计算total cost
func (e *ExperienceLevelPriceCap) calculateTotalCost(input InputLoanInfo) float64 {
	switch input.LoanType {
	case constants.LoanTypeRehab:
		return input.EstCurrentPropertyValue + input.OriginalBudget
	case constants.LoanTypeConstruction, constants.LoanTypeBridge:
		return input.EstCurrentPropertyValue
	default:
		return input.EstCurrentPropertyValue + input.OriginalBudget
	}
}

// getLTAIVCap 获取LTAIV上限
func (e *ExperienceLevelPriceCap) getLTAIVCap(experienceLevel ExperienceLevel, unitCount int32, isHighValueArea bool) float32 {
	var capValue float32
	switch experienceLevel {
	case ExperienceLevelNoExperienced:
		if unitCount >= 5 {
			capValue = 0.60 // (4+unit)规则
		} else {
			capValue = 0.65
		}
	case ExperienceLevelExperienced:
		if unitCount >= 5 {
			capValue = 0.65
		} else {
			capValue = 0.70
		}
	case ExperienceLevelInstitutional:
		if unitCount >= 5 {
			capValue = 0.70
		} else {
			capValue = 0.75
		}
	}

	// 高价值地区规则，只对institutional的borrower有效
	if isHighValueArea && experienceLevel == ExperienceLevelInstitutional {
		capValue = 0.80
	}

	return capValue + e.ExtraFactor/100
}

// getLTCCap 获取LTC上限
func (e *ExperienceLevelPriceCap) getLTCCap(experienceLevel ExperienceLevel, unitCount int32, isHighValueArea bool) float32 {
	var capValue float32
	switch experienceLevel {
	case ExperienceLevelNoExperienced:
		if unitCount >= 5 {
			capValue = 0.60
		} else {
			capValue = 0.65
		}
	case ExperienceLevelExperienced:
		if unitCount >= 5 {
			capValue = 0.65
		} else {
			capValue = 0.70
		}
	case ExperienceLevelInstitutional:
		if unitCount >= 5 {
			capValue = 0.70
		} else {
			capValue = 0.75
		}
	}

	// 高价值地区规则，只对institutional的borrower有效
	if isHighValueArea && experienceLevel == ExperienceLevelInstitutional {
		capValue = 0.80
	}

	return capValue + e.ExtraFactor/100
}

// getLTARVCap 获取LTARV上限
func (e *ExperienceLevelPriceCap) getLTARVCap(experienceLevel ExperienceLevel, unitCount int32, isHighValueArea bool) float32 {
	var capValue float32
	switch experienceLevel {
	case ExperienceLevelNoExperienced:
		if unitCount >= 5 {
			capValue = 0.65
		} else {
			capValue = 0.70
		}
	case ExperienceLevelExperienced:
		if unitCount >= 5 {
			capValue = 0.70
		} else {
			capValue = 0.75
		}
	case ExperienceLevelInstitutional:
		if unitCount >= 5 {
			capValue = 0.75
		} else {
			capValue = 0.80
		}
	}

	// 高价值地区规则，只对institutional的borrower有效
	if isHighValueArea && experienceLevel == ExperienceLevelInstitutional {
		capValue = 0.65
	}

	return capValue + e.ExtraFactor/100
}

// calculateInterestReserveAmount 计算Interest reserve amount
func (e *ExperienceLevelPriceCap) calculateInterestReserveAmount(totalLoanAmount float64, experienceLevel ExperienceLevel, months int32) float64 {
	defaultRate := e.getDefaultRate(experienceLevel)
	interestReserveAmount := float64(months) * (0.7 * defaultRate * totalLoanAmount)
	return interestReserveAmount
}

// getDefaultRate 获取默认利率
func (e *ExperienceLevelPriceCap) getDefaultRate(experienceLevel ExperienceLevel) float64 {
	switch experienceLevel {
	case ExperienceLevelInstitutional:
		return 0.10 // 10%
	case ExperienceLevelExperienced:
		return 0.09 // 9%
	default:
		return 0.0875 // 8.75%
	}
}

// calculateRate 计算利率
func (e *ExperienceLevelPriceCap) calculateRate(input InputLoanInfo, result *CalculatePreApprovedLoanAmountResult) float32 {
	// 计算新的LTC和LTARV
	var ltc, ltarv float32
	if input.EstCurrentPropertyValue+input.OriginalBudget > 0 {
		ltc = float32(result.TotalAmount / (input.EstCurrentPropertyValue + input.OriginalBudget))
	}
	if input.EstAfterRepairValue > 0 {
		ltarv = float32(result.TotalAmount / input.EstAfterRepairValue)
	}

	// 根据LTC和LTARV查找对应的利率，并取其中最小值作为基准利率
	ltcRate := e.getLTCRate(ltc, result.ExperienceLevel)
	ltarvRate := e.getLTARVRate(ltarv, result.ExperienceLevel)
	baseRate := float32(math.Min(float64(ltcRate), float64(ltarvRate)))

	// 计算调整因子
	adjusters := e.calculateAdjusters(input, result.ExperienceLevel)

	return baseRate + adjusters
}

// getLTCRate 获取LTC利率
func (e *ExperienceLevelPriceCap) getLTCRate(ltc float32, experienceLevel ExperienceLevel) float32 {
	// 这里需要根据实际的利率表来实现
	// 暂时使用简化逻辑
	if ltc <= 0.7 {
		return 0.08
	} else {
		return 0.09
	}
}

// getLTARVRate 获取LTARV利率
func (e *ExperienceLevelPriceCap) getLTARVRate(ltarv float32, experienceLevel ExperienceLevel) float32 {
	// 这里需要根据实际的利率表来实现
	// 暂时使用简化逻辑
	if ltarv <= 0.7 {
		return 0.08
	} else {
		return 0.09
	}
}

// calculateAdjusters 计算调整因子
func (e *ExperienceLevelPriceCap) calculateAdjusters(input InputLoanInfo, experienceLevel ExperienceLevel) float32 {
	var adjusters float32

	// FICO信用分调整
	if input.FicoScore < 650 {
		adjusters += 0.005 // +0.5%
	} else if input.FicoScore >= 700 {
		adjusters -= 0.0025 // -0.25%
	}

	// Ground up Construction调整
	if input.LoanType == constants.LoanTypeConstruction {
		adjusters += 0.005 // +0.5%
	}

	// Bridge cash out调整
	if input.LoanPurpose == constants.LoanPurposeCashOut {
		adjusters += 0.005 // +0.5%
	}

	// 经验等级调整
	if experienceLevel == ExperienceLevelNoExperienced {
		adjusters += 0.005 // +0.5%
	}

	return adjusters
}

// calculateMonthlyPayment 计算月供
func (e *ExperienceLevelPriceCap) calculateMonthlyPayment(totalLoanAmount float64, rate float64, loanTerm float64) float64 {
	if loanTerm <= 0 {
		return 0
	}
	monthlyRate := rate / 12
	return monthlyRate * totalLoanAmount / loanTerm
}

// calculateCashToClose 计算现金结算（手续费）
func (e *ExperienceLevelPriceCap) calculateCashToClose(input InputLoanInfo, result *CalculatePreApprovedLoanAmountResult) float64 {
	originationFee := 0.01 * result.TotalAmount    // 1% * total loan amount
	underwritingFee := 0.0025 * result.TotalAmount // 0.25% * total loan amount
	outOfPocket := e.calculateOutOfPocket(input, result)
	return originationFee + underwritingFee + outOfPocket
}

// calculateOutOfPocket 计算Out of pocket
func (e *ExperienceLevelPriceCap) calculateOutOfPocket(input InputLoanInfo, result *CalculatePreApprovedLoanAmountResult) float64 {
	switch input.LoanType {
	case constants.LoanTypeRehab:
		if input.LoanPurpose == constants.LoanPurposePurchase {
			assignmentFee := float64(0)
			if input.HasAssignmentFee {
				assignmentFee = input.AssignmentFeeAmount
			}
			return input.PurchasePrice + assignmentFee + input.OriginalBudget - result.TotalAmount
		} else { // refinance
			return input.OutstandingLoanAmount + input.OriginalBudget - result.TotalAmount
		}
	case constants.LoanTypeConstruction:
		assignmentFee := float64(0)
		if input.HasAssignmentFee {
			assignmentFee = input.AssignmentFeeAmount
		}
		interestReserveAmount := e.calculateInterestReserveAmount(result.TotalAmount, result.ExperienceLevel, input.InterestReserveMonths)
		if input.LoanPurpose == constants.LoanPurposePurchase {
			return input.PurchasePrice + assignmentFee + input.OriginalBudget + interestReserveAmount - result.TotalAmount
		} else { // refinance
			return input.OutstandingLoanAmount + input.OriginalBudget + interestReserveAmount - result.TotalAmount
		}
	case constants.LoanTypeBridge:
		assignmentFee := float64(0)
		if input.HasAssignmentFee {
			assignmentFee = input.AssignmentFeeAmount
		}
		if input.LoanPurpose == constants.LoanPurposePurchase {
			return input.PurchasePrice + assignmentFee - result.TotalAmount
		} else if input.LoanPurpose == constants.LoanPurposeCashOut {
			return input.OutstandingLoanAmount - result.TotalAmount
		} else { // refinance
			return input.OutstandingLoanAmount - result.TotalAmount
		}
	default:
		return 0
	}
}

// calculateInitialAndHoldbackAmount 计算初始贷款和修复贷款
func (e *ExperienceLevelPriceCap) calculateInitialAndHoldbackAmount(input InputLoanInfo, result *CalculatePreApprovedLoanAmountResult) (initialAmount, holdbackAmount float64) {
	// 当rehab loan的值=original budget
	holdbackAmount = input.OriginalBudget

	// Initial loan = total loan amount - original budget - interest reserve amount
	interestReserveAmount := float64(0)
	if input.LoanType == constants.LoanTypeConstruction {
		interestReserveAmount = e.calculateInterestReserveAmount(result.TotalAmount, result.ExperienceLevel, input.InterestReserveMonths)
	}
	initialAmount = result.TotalAmount - input.OriginalBudget - interestReserveAmount

	return initialAmount, holdbackAmount
}

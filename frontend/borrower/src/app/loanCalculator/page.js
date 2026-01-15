'use client';
import React, { useState } from 'react';
import EstimateLoan from './components/EstimateLoan';
import PreApproved from './components/PreApproved';

// 贷款计算器（预估，不准确，需要吸引客户申请） - 用于ZeeBee官网
const LoanCalculator = () => {
	const [isResult, setIsResult] = useState(false);
	const [result, setResult] = useState({});

	return isResult ? (
		<PreApproved result={result} setIsResult={setIsResult} />
	) : (
		<EstimateLoan setResult={setResult} setIsResult={setIsResult} />
	);
};
export default LoanCalculator;

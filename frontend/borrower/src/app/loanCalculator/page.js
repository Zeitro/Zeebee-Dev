'use client';
import React, { useState } from 'react';
import EstimateLoan from './components/EstimateLoan';
import PreApproved from './components/PreApproved';

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

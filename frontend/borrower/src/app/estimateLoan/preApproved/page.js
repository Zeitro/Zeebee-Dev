'use client';
import { Button, Form, Steps } from 'antd';
import React, { useEffect, useMemo, useState } from 'react';
// import styles from './index.module.scss';
import { useRouter } from 'next/navigation';
import { createRequest } from '@/lib/fetch';
import { formatAmount } from '@/lib/format';

// 预审计算loan
const PreApproved = () => {
	const [list, setList] = useState([{ name: '1' }]);
	const zeebeeList = [
		{
			label: 'Total loan amount',
			value: 'total_amount',
			render: (v) => formatAmount(v),
		},
		{
			label: 'Interest Rate',
			value: 'rate',
			render: (v) => (v ? `${v}%` : '--'),
		},
		{
			label: 'Monthly Payment',
			value: 'Monthly Payment',
			render: (v) => formatAmount(v),
		},
		{
			label: 'Cash to Close',
			value: '1',
			render: (v) => formatAmount(v),
		},
	];
	const otherList = [
		{
			label: 'Initial Loan',
			value: 'total_amount',
			render: (v) => formatAmount(v),
			className: '!w-1/2',
		},
		{
			label: 'Term',
			value: 'rate',
			render: (v) => (v ? `${v} Months` : '--'),
			className: '!w-1/2',
		},
		{
			label: 'Loan-to-Value',
			value: 'Monthly Payment',
			render: (v) => (v ? `${v}%` : '--'),
		},
		{
			label: 'Loan-to-Cost',
			value: '1',
			render: (v) => (v ? `${v}%` : '--'),
		},
		{
			label: 'Loan-to-ARV',
			value: '1',
			render: (v) => (v ? `${v}%` : '--'),
		},
	];

	const handleClick = () => {};
	return (
		<div className="max-w-[700px] my-10 overflow-hidden overflow-y-auto flex flex-col gap-4 [&::-webkit-scrollbar-thumb]:border-[#E4F3FF] [&::-webkit-scrollbar]:bg-[#E4F3FF]">
			<div className={['bg-white rounded-lg px-8 py-6 h-max text-center'].join(' ')}>
				<div className="text-3xl/10 font-medium text-[#0958D9]">Pre - Approved!</div>
				<div className="text-base font-medium text-black-dark pt-2">
					You're pre-approved for a bridge loan based on your project details.
				</div>
				<div className="py-12 flex flex-wrap items-center gap-y-10">
					{zeebeeList.map((item, index) => (
						<div
							key={index}
							className="flex flex-col flex-wrap items-center pt-4 w-1/2 gap-1"
						>
							<div className="text-[38px]/[46px] font-medium text-base-primary">
								{item?.render?.() || '--'}
							</div>
							<div className="text-base font-medium text-black-dark opacity-65">
								{item.label}
							</div>
						</div>
					))}
				</div>
				<div className="flex flex-wrap items-center gap-y-4">
					{otherList.map((item, index) => (
						<div key={index} className={`overflow-hidden w-1/3 ${item.className}`}>
							<div className="font-medium text-black-dark opacity-65 pt-3 pb-2 border-b border-[#C9CDD4]">
								{item.label}
							</div>
							<div className="text-base font-medium py-2">
								{item?.render?.() || '--'}
							</div>
						</div>
					))}
				</div>
				<div className="text-sm text-[#86909C] pb-4 pt-10 text-left">
					*The final loan amount will be determined based on the appraised value provided
					by an authorized third-party appraiser.
				</div>
				<div className="text-sm text-[#86909C] pb-10 text-left">
					*Monthly payment is estimated based on the total loan amount. Actual interest
					payments may vary depending on the draw schedule.
				</div>
				<Button type="primary" onClick={handleClick}>
					Accept loan terms
				</Button>
			</div>
			{list.map((item, index) => (
				<div key={index} className="bg-white rounded-lg px-8 py-6">
					<div>11</div>
				</div>
			))}
		</div>
	);
};
export default PreApproved;

'use client';
import { Button } from 'antd';
import { formatAmount } from '@repo/shared-methods/utils';

// 批准
const PreApprovedSuccess = ({ result, handleClick, back }) => {
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
			value: 'monthly_payment',
			render: (v) => formatAmount(v),
		},
		{
			label: 'Cash to Close',
			value: 'cash_to_close',
			render: (v) => formatAmount(v),
		},
	];
	const otherList = [
		{
			label: 'Initial Loan',
			value: 'initial_amount',
			render: (v) => formatAmount(v),
			className: 'md:w-1/2',
		},
		{
			label: 'Term',
			value: 'loan_term',
			render: (v) => (v ? `${v} Months` : '--'),
			className: 'md:w-1/2',
		},
		{
			label: 'Loan-to-Value',
			value: 'ltv',
			render: (v) => (v ? `${v}%` : '--'),
		},
		{
			label: 'Loan-to-Cost',
			value: 'ltc',
			render: (v) => (v ? `${v}%` : '--'),
		},
		{
			label: 'Loan-to-ARV',
			value: 'ltarv',
			render: (v) => (v ? `${v}%` : '--'),
		},
	];

	return (
		<div className={['bg-white rounded-lg px-8 py-6 h-max text-center'].join(' ')}>
			<div className="text-3xl/10 font-medium text-[#0958D9] relative">
				{back}Pre - Approved!
			</div>
			<div className="text-base font-medium text-black-dark pt-2">
				You're pre-approved for a bridge loan based on your project details.
			</div>
			<div className="py-12 flex flex-wrap items-center gap-y-10">
				{zeebeeList.map((item, index) => (
					<div
						key={index}
						className="flex flex-col flex-wrap items-center pt-4 w-full md:w-1/2 gap-1"
					>
						<div className="text-[38px]/[46px] font-medium text-base-primary">
							{item?.render?.(result?.primary?.[item.value]) || '--'}
						</div>
						<div className="text-base font-medium text-black-dark opacity-65">
							{item.label}
						</div>
					</div>
				))}
			</div>
			<div className="flex flex-wrap items-center gap-y-4">
				{otherList.map((item, index) => (
					<div
						key={index}
						className={`overflow-hidden w-full ${item.className || 'md:w-1/3'}`}
					>
						<div className="font-medium text-black-dark opacity-65 pt-3 pb-2 border-b border-[#C9CDD4]">
							{item.label}
						</div>
						<div className="text-base font-medium py-2">
							{item?.render?.(result?.primary?.[item.value]) || '--'}
						</div>
					</div>
				))}
			</div>
			<div className="text-sm text-[#86909C] pb-4 pt-10 text-left">
				*The final loan amount will be determined based on the appraised value provided by
				an authorized third-party appraiser.
			</div>
			<div className="text-sm text-[#86909C] pb-10 text-left">
				*Monthly payment is estimated based on the total loan amount. Actual interest
				payments may vary depending on the draw schedule.
			</div>
			<Button type="primary" onClick={() => handleClick(result?.primary?.lender_company_id)}>
				Accept loan terms
			</Button>
		</div>
	);
};
export default PreApprovedSuccess;

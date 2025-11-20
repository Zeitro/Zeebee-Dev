'use client';
import { Button } from 'antd';
import { formatAmount } from '@repo/shared-methods/utils';

const OtherLender = ({ result, handleClick }) => {
	const lenderList = [
		{
			label: 'Total loan amount',
			value: 'total_amount',
			render: (v) => `${formatAmount(v.min)} - ${formatAmount(v.max)}`,
		},
		{
			label: 'Interest Rate',
			value: 'rate',
			render: (v) => (v ? `${v.min}% - ${v.max}%` : '--'),
		},
		{
			label: 'Loan term (Months)',
			value: 'loan_term',
			render: (v) => (v ? `${v.min} - ${v.max}` : '--'),
		},
	];

	return (
		<>
			<div className="pt-12 pb-6">
				<p className="text-xl font-bold text-[#1D2129]">
					Not the right fit? We've got more options for you.
				</p>
				<p className="text-base text-black-dark pt-4">
					Based on your project and documents, we found offers from other lenders that
					might better match your goals.
				</p>
			</div>
			{(result?.secondaries || [])?.map((item, index) => (
				<div
					key={index}
					className="bg-white rounded-lg px-6 py-6 mb-4 [&:last-child]:mb-0]"
				>
					<div className="flex justify-between items-center pb-12">
						<p className="text-xl font-bold text-black-dark">
							{item?.lender_company_name}
						</p>
						<Button
							type="primary"
							onClick={() => handleClick(item?.lender_company_id, true)}
						>
							Send loan request
						</Button>
					</div>
					<div className="flex items-center flex-col md:flex-row gap-4">
						{lenderList.map((i, j) => (
							<div key={j} className="w-full lg:w-1/3 flex items-center flex-col">
								<p className="text-2xl font-bold text-[#0E42D2]">
									{i?.render?.(item?.[i.value]) || '--'}
								</p>
								<div className="text-base text-[#4E5969] font-medium">
									{i.label}
								</div>
							</div>
						))}
					</div>
				</div>
			))}
		</>
	);
};
export default OtherLender;

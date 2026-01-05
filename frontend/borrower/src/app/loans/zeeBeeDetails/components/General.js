'use client';
import { Col, Row } from 'antd';
import React from 'react';
import { formatAmount } from '@repo/shared-methods/utils';
export const generalTermsFields = [
	{
		label: 'Total loan amount',
		value: 'total_amount',
		render: (v) => formatAmount(v),
	},
	{
		label: 'Rate',
		value: 'rate',
		render: (v) => (v ? `${v}%` : '--'),
	},
	{
		label: 'Outstanding loan amount',
		value: 'principal_amount',
		render: (v) => formatAmount(v),
	},
	{ label: 'Maturity date', value: 'maturity_date' },
	{
		label: 'Monthly Payment',
		value: 'monthly_payment',
		render: (v) => formatAmount(v),
	},
	{
		label: 'Term',
		value: 'loan_term',
		render: (v) => `${v || '--'} Months`,
	},
];

export const generalParametersFields = [
	{
		label: 'Rate type',
		value: 'lien_position',
	},
	{
		label: 'Due date',
		value: 'due_date',
	},
	{
		label: 'Frequency',
		value: 'frequency',
	},
	{
		label: 'Amortization type',
		value: 'amortization_type',
	},
	{
		label: 'Amortization Term',
		value: 'amortization_term',
	},
	{
		label: 'First payment',
		value: 'first_payment',
	},
];

export const loanMatrixFields = [
	{
		label: 'Loan To Value (LTV)',
		value: 'ltv1',
	},
	{
		label: 'Loan To Cost (LTC)',
		value: 'ltc',
	},
	{
		label: 'After Repair LTV',
		value: 'after_repair_ltv',
	},
	{
		label: 'Blended LTC',
		value: 'blended_ltc',
	},
];

export const collateralSummaryFields = [
	{
		label: 'Purchase price',
		value: 'purchase_amount',
		render: (v) => formatAmount(v),
	},
	{
		label: 'Rehab cost ',
		value: 'rehab_cost',
	},
	{
		label: 'Estimated current value',
		value: 'est_property_value_amount',
		render: (v) => formatAmount(v),
	},
	{
		label: 'Estimated after Repair value',
		value: 'arv',
		render: (v) => formatAmount(v),
	},
];
export const generalList = [
	{
		title: 'Terms',
		fields: generalTermsFields,
	},
	{
		title: 'Parameters',
		fields: generalParametersFields,
	},
	{
		title: 'Loan Matrix',
		fields: loanMatrixFields,
	},
	{
		title: 'Collateral Summary',
		fields: collateralSummaryFields,
	},
];
const General = ({ rowData = {}, className }) => {
	return (
		<Row gutter={[40, 24]} className={className}>
			{generalList.map((i, j) => {
				return (
					<Col key={j} sm={24} xs={24} md={12}>
						<div className="border border-base-border rounded-2xl pt-6 pb-4 px-6">
							<div className="text-xl font-medium text-black-dark pl-3 pb-4 border-b border-base-border">
								{i.title}
							</div>
							<Row gutter={[24, 24]} className="px-4 py-4">
								{i.fields?.map((item, index) => (
									<Col sm={24} xs={24} md={12} key={index}>
										<div
											className="text-sm/[22px] pb-2"
											style={{ color: 'rgba(0, 0, 0, 0.65)' }}
										>
											{item.label}
										</div>
										<div className="text-xl font-medium text-black-dark">
											{item?.render
												? item?.render?.(rowData[item.value])
												: rowData[item.value] || '--'}
										</div>
									</Col>
								))}
							</Row>
						</div>
					</Col>
				);
			})}
		</Row>
	);
};

export default General;

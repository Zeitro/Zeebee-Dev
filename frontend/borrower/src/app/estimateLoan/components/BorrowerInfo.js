import React from 'react';
// import styles from './index.module.scss';
import { Row, Col, Form, Select, Input, InputNumber } from 'antd';
import {
	emailRegex,
	creditScoreRegex,
	minCreditScore,
	maxCreditScore,
	validateRule,
	capitalizeFirstLetter,
} from '@/lib/format';
import { titleCol } from '@/enum/common';
import { useGetOptions } from '@/hooks/useGetOptions';
const BorrowerInfo = ({ gutter = 24, children }) => {
	const { loanTypeOptions, loanPurposeOptions } = useGetOptions([
		{
			key: 'loan_type',
			valueName: 'loanTypeOptions',
		},
		{
			key: 'loan_purpose',
			valueName: 'loanPurposeOptions',
		},
	]);
	return (
		<Row gutter={gutter}>
			<Col {...titleCol} className="!text-3xl/10 font-medium pb-8">
				Borrower information
			</Col>
			<Col {...titleCol}>
				<Form.Item
					label="Name"
					name="name"
					rules={[
						{
							required: true,
						},
					]}
					getValueFromEvent={(v) => capitalizeFirstLetter(v.target.value)}
				>
					<Input placeholder="Enter Borrower name" />
				</Form.Item>
			</Col>
			<Col {...titleCol}>
				<Form.Item
					label="Email"
					name="email"
					rules={[
						{
							required: true,
						},
						{
							validator: (_, value) =>
								validateRule(value, emailRegex, 'Invalid email format'),
						},
					]}
				>
					<Input placeholder="Enter email" />
				</Form.Item>
			</Col>
			<Col {...titleCol}>
				<Form.Item label="Loan type" name="loan_type" rules={[{ required: true }]}>
					<Select options={loanTypeOptions} placeholder="Loan type" />
				</Form.Item>
			</Col>
			<Col {...titleCol}>
				<Form.Item label="Loan purpose" name="loan_purpose" rules={[{ required: true }]}>
					<Select options={loanPurposeOptions} placeholder="Loan purpose" />
				</Form.Item>
			</Col>
			<Col {...titleCol}>
				<Form.Item
					label="FICO credit score"
					name="est_credit_score"
					rules={[
						{
							required: true,
						},
						{
							validator: (_, value) =>
								validateRule(value, creditScoreRegex, 'Invalid Credit Score'),
						},
					]}
				>
					<InputNumber
						style={{ width: '100%' }}
						min={minCreditScore}
						max={maxCreditScore}
						controls={false}
						placeholder="FICO credit score"
					/>
				</Form.Item>
			</Col>
			<Col {...titleCol} className="text-base font-medium text-black-dark pt-6 pb-4">
				In the last 3 years, how many times have you...
			</Col>
			<Col {...titleCol}>
				<Form.Item label="Flipped a property?" name="property">
					<InputNumber
						style={{ width: '100%' }}
						// min={minCreditScore}
						// max={maxCreditScore}
						controls={false}
						placeholder="Number of properties"
					/>
				</Form.Item>
			</Col>
			<Col {...titleCol}>
				<Form.Item label="Purchased a rental property?" name="rental_property">
					<InputNumber
						style={{ width: '100%' }}
						// min={minCreditScore}
						// max={maxCreditScore}
						controls={false}
						placeholder="Number of rental properties"
					/>
				</Form.Item>
			</Col>
			{children}
		</Row>
	);
};
export default BorrowerInfo;

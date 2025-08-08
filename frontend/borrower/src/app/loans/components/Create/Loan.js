'use client';
import { Col, DatePicker, Form, Radio, Row, Select } from 'antd';
import React from 'react';
import styles from './index.module.scss';
import InputNumAmount from '@/components/InputNumAmount';
import { useGetOptions } from '@/hooks/useGetOptions';
import { loan_active, loan_archived, loan_pending } from '@/enum/common';
import { formatStr } from '@/lib/date';
const MAX_COUNT = 5;
const CreateLoan = ({ form, stepBtn, gutter, ...colProps }) => {
	const stage = Form.useWatch('stage', form);
	const {
		loanPendingStatus,
		loanPurposeOptions,
		loanTypeOptions,
		loanStatusOptions,
		businessStrategyOptions,
	} = useGetOptions([
		{
			key: 'loan_purpose',
			valueName: 'loanPurposeOptions',
		},
		{
			key: 'loan_type',
			valueName: 'loanTypeOptions',
		},
		{
			key: 'loan_pending_status',
			valueName: 'loanPendingStatus',
		},
		{
			key: 'loan_stage',
			valueName: 'loanStatusOptions',
		},
		{
			key: 'business_strategy',
			valueName: 'businessStrategyOptions',
		},
	]);
	return (
		<Row className={styles.create} gutter={gutter}>
			<Col md={24}>
				<Form.Item
					label="Select a loan stage"
					name="stage"
					className={styles.commonFormItem}
					rules={[{ required: true, message: 'Please select a loan stage!' }]}
				>
					<Radio.Group options={loanStatusOptions} />
				</Form.Item>
			</Col>
			{stage === loan_pending && (
				<Col md={24}>
					<Form.Item
						label="Select a pending status"
						name="pending_status"
						className={styles.commonFormItem}
						rules={[{ required: true, message: 'Please select a pending status!' }]}
					>
						<Radio.Group options={loanPendingStatus} />
					</Form.Item>
				</Col>
			)}
			<Col md={24} {...(colProps || {})}>
				<Form.Item
					label="Loan purpose"
					name="loan_purpose"
					rules={[{ required: true, message: 'Please select a loan purpose!' }]}
				>
					<Select options={loanPurposeOptions} placeholder="Select a loan purpose" />
				</Form.Item>
			</Col>
			<Col md={24} {...(colProps || {})}>
				<Form.Item
					label="Loan type"
					name="loan_type"
					rules={[{ required: true, message: 'Please select a loan type!' }]}
				>
					<Select options={loanTypeOptions} placeholder="Select a loan type" />
				</Form.Item>
			</Col>
			{stage !== loan_archived && (
				<>
					<Col md={24} {...(colProps || {})}>
						<Form.Item
							label="Total loan amount"
							name="total_amount"
							rules={[
								{ required: true, message: 'Please iEnter total loan amount!' },
							]}
						>
							<InputNumAmount placeholder="Enter total loan amount" />
						</Form.Item>
					</Col>
					<Col md={24} {...(colProps || {})}>
						<Form.Item label="Initial loan amount(optional)" name="initial_amount">
							<InputNumAmount placeholder="Enter initial loan amount" />
						</Form.Item>
					</Col>
				</>
			)}

			{stage === loan_active && (
				<>
					<Col md={24} {...(colProps || {})}>
						<Form.Item
							label="Outstanding loan amount"
							name="principal_amount"
							rules={[
								{
									required: true,
									message: 'Please Enter Outstanding loan amount!',
								},
							]}
						>
							<InputNumAmount placeholder="Enter Outstanding loan amount" />
						</Form.Item>
					</Col>
					<Col md={24} {...(colProps || {})}>
						<Form.Item
							label="Monthly paymant"
							name="monthly_payment"
							rules={[{ required: true, message: 'Please Enter Monthly paymant!' }]}
						>
							<InputNumAmount placeholder="Enter Monthly paymant" />
						</Form.Item>
					</Col>
				</>
			)}
			{stage === loan_archived && (
				<>
					<Col md={24} {...(colProps || {})}>
						<Form.Item
							label="Principal balance"
							name="principal_amount"
							rules={[
								{
									required: true,
									message: 'Please Enter Principal balance!',
								},
							]}
						>
							<InputNumAmount placeholder="Enter Principal balance" />
						</Form.Item>
					</Col>
					<Col md={24} {...(colProps || {})}>
						<Form.Item
							label="Paid off date"
							name="paid_off_date"
							rules={[{ required: true, message: 'Please select paid off date!' }]}
						>
							<DatePicker
								format={formatStr}
								placeholder="select"
								style={{ width: '100%' }}
							/>
						</Form.Item>
					</Col>
				</>
			)}
			{stage !== loan_archived && (
				<Col md={24} {...(colProps || {})}>
					<Form.Item
						label="Business strategy"
						name="business_strategy"
						rules={[{ required: true, message: 'Please select Business strategy!' }]}
					>
						<Select
							mode="multiple"
							maxCount={MAX_COUNT}
							style={{ width: '100%' }}
							placeholder="Select Business strategy"
							options={businessStrategyOptions}
						/>
					</Form.Item>
				</Col>
			)}

			{stepBtn}
		</Row>
	);
};
export default CreateLoan;

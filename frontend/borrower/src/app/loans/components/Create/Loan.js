'use client';
import { Col, Form, Radio, Row, Select, Input } from 'antd';
import React from 'react';
import styles from './index.module.scss';
// import InputNumAmount from '@/components/InputNumAmount';
import { InputNumAmount } from '@repo/shared-methods';
import { useGetOptions } from '@/hooks/useGetOptions';
import { loan_active, loan_archived, loan_pending } from '@/enum/common';
import { capitalizeFirstLetter, maxAmount } from '@/lib/format';
import CommonFormDate from '@/components/CommonDate';
const MAX_COUNT = 5;
const CreateLoan = ({ form, stepBtn, gutter, isRadioStyle, ...colProps }) => {
	const stage = Form.useWatch('stage', form);
	const total_amount = Form.useWatch('total_amount', form) || maxAmount;
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
	const layoutCol = {
		xs: 24,
		sm: 24,
	};
	return (
		<Row className={styles.create} gutter={gutter}>
			<Col {...layoutCol}>
				<Form.Item
					label="Select a loan stage"
					name="stage"
					className={styles.commonFormItem}
					rules={[{ required: true, message: 'Please select a loan stage!' }]}
				>
					<Radio.Group
						options={loanStatusOptions}
						className={isRadioStyle && styles.radio}
					/>
				</Form.Item>
			</Col>
			{stage === loan_pending && (
				<Col {...layoutCol}>
					<Form.Item
						label="Select a pending status"
						name="pending_status"
						className={styles.commonFormItem}
						rules={[{ required: true, message: 'Please select a pending status!' }]}
					>
						<Radio.Group
							options={loanPendingStatus}
							className={isRadioStyle && styles.radio}
						/>
					</Form.Item>
				</Col>
			)}
			<Col {...layoutCol} {...(colProps || {})}>
				<Form.Item
					label="Loan purpose"
					name="loan_purpose"
					rules={[{ required: true, message: 'Please select a loan purpose!' }]}
				>
					<Select options={loanPurposeOptions} placeholder="Select a loan purpose" />
				</Form.Item>
			</Col>
			<Col {...layoutCol} {...(colProps || {})}>
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
					<Col {...layoutCol} {...(colProps || {})}>
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
					<Col {...layoutCol} {...(colProps || {})}>
						<Form.Item label="Initial loan amount(optional)" name="initial_amount">
							<InputNumAmount
								placeholder="Enter initial loan amount"
								max={total_amount}
							/>
						</Form.Item>
					</Col>
				</>
			)}

			{stage === loan_active && (
				<>
					<Col {...layoutCol} {...(colProps || {})}>
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
							<InputNumAmount
								placeholder="Enter Outstanding loan amount"
								max={total_amount}
							/>
						</Form.Item>
					</Col>
					<Col {...layoutCol} {...(colProps || {})}>
						<Form.Item
							label="Monthly payment"
							name="monthly_payment"
							rules={[{ required: true, message: 'Please Enter Monthly payment!' }]}
						>
							<InputNumAmount
								placeholder="Enter Monthly payment"
								max={total_amount}
							/>
						</Form.Item>
					</Col>
				</>
			)}
			{stage === loan_archived && (
				<>
					<Col {...layoutCol} {...(colProps || {})}>
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
					<Col {...layoutCol} {...(colProps || {})}>
						<Form.Item
							label="Paid off date"
							name="paid_off_date"
							rules={[{ required: true, message: 'Please select paid off date!' }]}
						>
							<CommonFormDate />
						</Form.Item>
					</Col>
				</>
			)}
			<Col {...layoutCol} {...(colProps || {})}>
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
			{stage === loan_pending && (
				<>
					<Col {...layoutCol} {...(colProps || {})}>
						<Form.Item
							label="Existing lender"
							name="existing_lender_company"
							getValueFromEvent={(v) => capitalizeFirstLetter(v.target.value)}
						>
							<Input placeholder="Enter existing lender" />
						</Form.Item>
					</Col>
					<Col {...layoutCol} {...(colProps || {})}>
						<Form.Item label="Business line" name={['extend', 'business_line']}>
							<Input placeholder="Enter business line" />
						</Form.Item>
					</Col>
					<Col {...layoutCol}></Col>
					<Col {...layoutCol} {...(colProps || {})}>
						<Form.Item label="Pricing" name={['extend', 'pricing']}>
							<Input.TextArea
								placeholder="Enter pricing"
								autoSize={{ minRows: 2, maxRows: 3 }}
							/>
						</Form.Item>
					</Col>
					<Col {...layoutCol} {...(colProps || {})}>
						<Form.Item
							label="Loan status updates"
							name={['extend', 'loan_status_updates']}
						>
							<Input.TextArea
								placeholder="Enter loan status updates"
								autoSize={{ minRows: 2, maxRows: 3 }}
							/>
						</Form.Item>
					</Col>
					<Col {...layoutCol} {...(colProps || {})}>
						<Form.Item label="Meeting notes" name={['extend', 'meeting_notes']}>
							<Input.TextArea
								placeholder="Enter meeting notes"
								autoSize={{ minRows: 2, maxRows: 3 }}
							/>
						</Form.Item>
					</Col>
					<Col {...layoutCol} {...(colProps || {})}>
						<Form.Item label="Action list" name={['extend', 'action_list']}>
							<Input.TextArea
								placeholder="Enter action list"
								autoSize={{ minRows: 2, maxRows: 3 }}
							/>
						</Form.Item>
					</Col>
				</>
			)}

			{stepBtn}
		</Row>
	);
};
export default CreateLoan;

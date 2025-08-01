'use client';
import { Button, Col, DatePicker, Form, InputNumber, Radio, Row, Select, Steps, Tabs } from 'antd';
import React, { useState } from 'react';
import {
	loanStatusOptions,
	loanPendingStatus,
	loanPurposeOptions,
	loanTypeOptions,
} from '@/enum/loan';
import styles from './index.module.scss';
import { useRouter } from 'next/navigation';
const CreateLoan = ({ form, stepBtn, gutter, ...colProps }) => {
	const stage = Form.useWatch('stage', form);
	const router = useRouter();
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
			{stage === 'pending' && (
				<Col md={24}>
					<Form.Item
						label="Select a pending status"
						name="pendingStatus"
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
					name="loanPurpose"
					rules={[{ required: true, message: 'Please select a loan purpose!' }]}
				>
					<Select options={loanPurposeOptions} placeholder="Select a loan purpose" />
				</Form.Item>
			</Col>
			<Col md={24} {...(colProps || {})}>
				<Form.Item
					label="Loan type"
					name="loanType"
					rules={[{ required: true, message: 'Please select a loan type!' }]}
				>
					<Select options={loanTypeOptions} placeholder="Select a loan type" />
				</Form.Item>
			</Col>
			{stage !== 'archived' && (
				<>
					<Col md={24} {...(colProps || {})}>
						<Form.Item
							label="Total loan amount"
							name="totalLoanAmount"
							rules={[
								{ required: true, message: 'Please iEnter total loan amount!' },
							]}
						>
							<InputNumber
								formatter={(value) =>
									`$ ${value}`.replace(/\B(?=(\d{3})+(?!\d))/g, ',')
								}
								parser={(value) =>
									value === null || value === void 0
										? void 0
										: value.replace(/\$\s?|(,*)/g, '')
								}
								style={{ width: '100%' }}
								min={1}
								max={10}
								controls={false}
								placeholder="Enter total loan amount"
							/>
						</Form.Item>
					</Col>
					<Col md={24} {...(colProps || {})}>
						<Form.Item label="Initial loan amount(optional)" name="initialLoanAmount">
							<InputNumber
								style={{ width: '100%' }}
								formatter={(value) =>
									`$ ${value}`.replace(/\B(?=(\d{3})+(?!\d))/g, ',')
								}
								parser={(value) =>
									value === null || value === void 0
										? void 0
										: value.replace(/\$\s?|(,*)/g, '')
								}
								min={1}
								max={10}
								controls={false}
								placeholder="Enter initial loan amount"
							/>
						</Form.Item>
					</Col>
				</>
			)}

			{stage === 'active' && (
				<>
					<Col md={24} {...(colProps || {})}>
						<Form.Item
							label="Outstanding loan amount"
							name="OutstandingLoanAmount"
							rules={[
								{
									required: true,
									message: 'Please Enter Outstanding loan amount!',
								},
							]}
						>
							<InputNumber
								style={{ width: '100%' }}
								formatter={(value) =>
									`$ ${value}`.replace(/\B(?=(\d{3})+(?!\d))/g, ',')
								}
								parser={(value) =>
									value === null || value === void 0
										? void 0
										: value.replace(/\$\s?|(,*)/g, '')
								}
								min={1}
								max={10}
								controls={false}
								placeholder="Enter Outstanding loan amount"
							/>
						</Form.Item>
					</Col>
					<Col md={24} {...(colProps || {})}>
						<Form.Item
							label="Monthly paymant"
							name="monthlyPaymant"
							rules={[{ required: true, message: 'Please Enter Monthly paymant!' }]}
						>
							<InputNumber
								style={{ width: '100%' }}
								formatter={(value) =>
									`$ ${value}`.replace(/\B(?=(\d{3})+(?!\d))/g, ',')
								}
								parser={(value) =>
									value === null || value === void 0
										? void 0
										: value.replace(/\$\s?|(,*)/g, '')
								}
								min={1}
								max={10}
								controls={false}
								placeholder="Enter Monthly paymant"
							/>
						</Form.Item>
					</Col>
				</>
			)}
			{stage === 'archived' && (
				<>
					<Col md={24} {...(colProps || {})}>
						<Form.Item
							label="Principal balance"
							name="principalBalance"
							rules={[
								{
									required: true,
									message: 'Please Enter Principal balance!',
								},
							]}
						>
							<InputNumber
								style={{ width: '100%' }}
								formatter={(value) =>
									`$ ${value}`.replace(/\B(?=(\d{3})+(?!\d))/g, ',')
								}
								parser={(value) =>
									value === null || value === void 0
										? void 0
										: value.replace(/\$\s?|(,*)/g, '')
								}
								min={1}
								max={10}
								controls={false}
								placeholder="Enter Principal balance"
							/>
						</Form.Item>
					</Col>
					<Col md={24} {...(colProps || {})}>
						<Form.Item
							label="Paid off date"
							name="padiOffDate"
							rules={[{ required: true, message: 'Please select paid off date!' }]}
						>
							<DatePicker formatter="MM/DD/YYYY" placeholder="select" />
						</Form.Item>
					</Col>
				</>
			)}
			{stage !== 'archived' && (
				<Col md={24} {...(colProps || {})}>
					<Form.Item
						label="Business strategy"
						name="businessStrategy"
						rules={[{ required: true, message: 'Please select Business strategy!' }]}
					>
						<Select
							options={loanPurposeOptions}
							placeholder="Select Business strategy"
						/>
					</Form.Item>
				</Col>
			)}

			{stepBtn}
		</Row>
	);
};
export default CreateLoan;

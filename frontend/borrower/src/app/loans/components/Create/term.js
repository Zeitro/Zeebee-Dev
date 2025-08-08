'use client';
import { Col, DatePicker, Form, Input, InputNumber, Row, Select } from 'antd';
import React from 'react';
import styles from './index.module.scss';
import InputNumAmount from '@/components/InputNumAmount';
import { useGetOptions } from '@/hooks/useGetOptions';
import { loan_active, loan_pending } from '@/enum/common';
import { formatStr } from '@/lib/date';
const CreateTerm = ({ stage, stepBtn, gutter, ...colProps }) => {
	const { lienPositionOptions, prepaymentPenaltyOptions } = useGetOptions([
		{
			key: 'lien_position',
			valueName: 'lienPositionOptions',
		},
		{
			key: 'loan_prepayment_penalty',
			valueName: 'prepaymentPenaltyOptions',
		},
	]);
	const rateDom = (
		<Col md={24} {...(colProps || {})}>
			<Form.Item
				label="Rate"
				name="rate"
				rules={[{ required: true, message: 'Please input!' }]}
			>
				<InputNumber
					style={{ width: '100%' }}
					min={1}
					max={100}
					controls={false}
					placeholder="Enter Rate"
					suffix="%"
				/>
			</Form.Item>
		</Col>
	);
	const lenderDom = (
		<Col md={24} {...(colProps || {})}>
			<Form.Item
				label={stage === 'pending' ? 'TagetLender' : 'Lender'}
				name="lender_company"
				rules={[{ required: true, message: 'Please input!' }]}
			>
				<Input placeholder={`Enter ${stage === 'pending' ? 'taget ' : ''}lender`} />
			</Form.Item>
		</Col>
	);
	return (
		<Row className={styles.create} gutter={gutter}>
			<Form.List name="applicants">
				{(fields) => {
					return fields?.map((field, index) => (
						<Col
							md={24}
							{...(colProps || {})}
							key={field.key}
							fieldkey={field.fieldKey}
						>
							<Form.Item
								label={`Borrower ${index + 1}`}
								name={[field.name, 'name']}
								rules={[{ required: !index, message: 'Please input!' }]}
							>
								<Input placeholder="Enter Borrower" />
							</Form.Item>
						</Col>
					));
				}}
			</Form.List>
			{stage === loan_pending && (
				<>
					{lenderDom}
					<Col md={24} {...(colProps || {})}>
						<Form.Item
							label="Target closing date"
							name="closing_date"
							rules={[{ required: true, message: 'Please select date!' }]}
						>
							<DatePicker
								format={formatStr}
								placeholder="Select target closing date"
								style={{ width: '100%' }}
							/>
						</Form.Item>
					</Col>
					<Col md={24} {...(colProps || {})}>
						<Form.Item
							label="Holdback"
							name="holdback_amount"
							rules={[{ required: true, message: 'Please input!' }]}
						>
							<InputNumAmount placeholder="Enter holdback amount" />
						</Form.Item>
					</Col>
					{rateDom}
					<Col md={24} {...(colProps || {})}>
						<Form.Item
							label="Estimated closing fund"
							name="closing_fund"
							rules={[{ required: true, message: 'Please input!' }]}
						>
							<InputNumAmount placeholder="Enter Estimated closing fund" />
						</Form.Item>
					</Col>
				</>
			)}
			{stage !== loan_pending && (
				<>
					{lenderDom}
					{rateDom}
				</>
			)}
			{stage === loan_active && (
				<>
					{/* <Col md={24} {...(colProps||{})}>
						<Form.Item
							label="Point"
							name="point"
							rules={[{ required: true, message: 'Please input!' }]}
						>
							<InputNumber
								style={{ width: '100%' }}
								min={1}
								max={10}
								controls={false}
								placeholder="Enter point"
							/>
						</Form.Item>
					</Col> */}
					<Col md={24} {...(colProps || {})}>
						<Form.Item
							label="Lien position"
							name="lien_position"
							rules={[{ required: true, message: 'Please select Lien position!' }]}
						>
							<Select
								options={lienPositionOptions}
								placeholder="Select Lien position"
							/>
						</Form.Item>
					</Col>
					<Col md={24} {...(colProps || {})}>
						<Form.Item
							label="Maturity date"
							name="maturity_date"
							rules={[{ required: true, message: 'Please select date!' }]}
						>
							<DatePicker
								format={formatStr}
								placeholder="Select Maturity date"
								style={{ width: '100%' }}
							/>
						</Form.Item>
					</Col>
					<Col md={24} {...(colProps || {})}>
						<Form.Item
							label="Prepayment penalty"
							name="prepayment_penalty"
							rules={[{ required: true, message: 'Please input!' }]}
						>
							<Select
								options={prepaymentPenaltyOptions}
								placeholder="Enter Prepayment penalty"
							/>
							{/* <InputNumAmount placeholder="Enter Prepayment penalty" /> */}
						</Form.Item>
					</Col>
					<Col md={24} {...(colProps || {})}>
						<Form.Item
							label="Servicing company"
							name="service_company"
							rules={[{ required: true, message: 'Please input!' }]}
						>
							<Input placeholder="Enter Servicing company" />
						</Form.Item>
					</Col>
				</>
			)}
			{stepBtn}
		</Row>
	);
};
export default CreateTerm;

'use client';
import {
	Button,
	Col,
	DatePicker,
	Form,
	Input,
	InputNumber,
	Radio,
	Row,
	Select,
	Space,
	Steps,
	Tabs,
} from 'antd';
import React, { useState } from 'react';
import {
	loanStatusOptions,
	loanPendingStatus,
	loanPurposeOptions,
	loanTypeOptions,
} from '@/enum/loan';
import styles from './index.module.scss';
import { usePropertyStore } from '@/store/property';
import { useRouter } from 'next/navigation';
const CreateTerm = ({  stage,stepBtn,gutter,...colProps }) => {
	const { propertyTypeOptions } = usePropertyStore();
	const router = useRouter();
	const rateDom = (
		<Col md={24} {...(colProps||{})}>
			<Form.Item
				label="Rate"
				name="rate"
				rules={[{ required: true, message: 'Please input!' }]}
			>
				<InputNumber
					style={{ width: '100%' }}
					min={1}
					max={10}
					controls={false}
					placeholder="Enter Rate"
				/>
			</Form.Item>
		</Col>
	);
	const lenderDom = (
		<Col md={24} {...(colProps||{})}>
			<Form.Item
				label={stage === 'pending' ? 'TagetLender' : 'Lender'}
				name="lender"
				rules={[{ required: true, message: 'Please input!' }]}
			>
				<Input placeholder={`Enter ${stage === 'pending' ? 'taget ' : ''}lender`} />
			</Form.Item>
		</Col>
	);
	return (
		<Row className={styles.create} gutter={gutter}>
			<Col md={24} {...(colProps||{})}>
				<Form.Item
					label="Borrower"
					name="borrower"
					rules={[{ required: true, message: 'Please input!' }]}
				>
					<Input placeholder="Enter Borrower" />
				</Form.Item>
			</Col>
			{stage === 'pending' && (
				<>
					{lenderDom}
					<Col md={24} {...(colProps||{})}>
						<Form.Item
							label="Target closing date"
							name="targetClosingDate"
							rules={[{ required: true, message: 'Please select date!' }]}
						>
							<DatePicker
								formatter="MM/DD/YYYY"
								placeholder="Select target closing date"
								style={{ width: '100%' }}
							/>
						</Form.Item>
					</Col>
					<Col md={24} {...(colProps||{})}>
						<Form.Item
							label="Holdback"
							name="holdback"
							rules={[{ required: true, message: 'Please input!' }]}
						>
							<InputNumber
								style={{ width: '100%' }}
								min={1}
								max={10}
								controls={false}
								placeholder="Enter holdback amount"
								formatter={(value) =>
									`$ ${value}`.replace(/\B(?=(\d{3})+(?!\d))/g, ',')
								}
								parser={(value) =>
									value === null || value === void 0
										? void 0
										: value.replace(/\$\s?|(,*)/g, '')
								}
							/>
						</Form.Item>
					</Col>
					{rateDom}
					<Col md={24} {...(colProps||{})}>
						<Form.Item
							label="Estimated closing fund"
							name="estimatedClosingFund"
							rules={[{ required: true, message: 'Please input!' }]}
						>
							<InputNumber
								style={{ width: '100%' }}
								min={1}
								max={10}
								controls={false}
								placeholder="Enter Estimated closing fund"
								formatter={(value) =>
									`$ ${value}`.replace(/\B(?=(\d{3})+(?!\d))/g, ',')
								}
								parser={(value) =>
									value === null || value === void 0
										? void 0
										: value.replace(/\$\s?|(,*)/g, '')
								}
							/>
						</Form.Item>
					</Col>
				</>
			)}
			{stage !== 'pending' && (
				<>
					{lenderDom}
					{rateDom}
				</>
			)}
			{stage === 'active' && (
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
					<Col md={24} {...(colProps||{})}>
						<Form.Item
							label="Lien position"
							name="lienPosition"
							rules={[{ required: true, message: 'Please select Lien position!' }]}
						>
							<Select options={[]} placeholder="Select Lien position" />
						</Form.Item>
					</Col>
					<Col md={24} {...(colProps||{})}>
						<Form.Item
							label="Maturity date"
							name="maturityDate"
							rules={[{ required: true, message: 'Please select date!' }]}
						>
							<DatePicker
								formatter="MM/DD/YYYY"
								placeholder="Select Maturity date"
								style={{ width: '100%' }}
							/>
						</Form.Item>
					</Col>
					<Col md={24} {...(colProps||{})}>
						<Form.Item
							label="Prepayment penalty"
							name="prepaymentPenalty"
							rules={[{ required: true, message: 'Please input!' }]}
						>
							<InputNumber
								style={{ width: '100%' }}
								min={1}
								max={10}
								controls={false}
								placeholder="Enter Prepayment penalty"
							/>
						</Form.Item>
					</Col>
					<Col md={24} {...(colProps||{})}>
						<Form.Item
							label='Servicing company'
							name="servicingCompany"
							rules={[{ required: true, message: 'Please input!' }]}
						>
							<Input
								placeholder='Enter Servicing company'
							/>
						</Form.Item>
					</Col>
				</>
			)}
			{stepBtn}
		</Row>
	);
};
export default CreateTerm;

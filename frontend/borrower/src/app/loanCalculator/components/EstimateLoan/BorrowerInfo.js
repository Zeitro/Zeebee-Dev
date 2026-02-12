import React from 'react';
// import styles from './index.module.scss';
import { Row, Col, Form, Input, InputNumber } from 'antd';
import {
	emailRegex,
	creditScoreRegex,
	minCreditScore,
	maxCreditScore,
	validateRule,
	capitalizeFirstLetter,
	nameRegex,
} from '@/lib/format';
import { titleCol } from '@/enum/common';
import InputPhone from '@/components/InputPhone';
const BorrowerInfo = ({ gutter = 24, children }) => {
	return (
		<Row gutter={gutter}>
			<Col {...titleCol} className="!text-3xl/10 font-medium pb-8">
				Borrower information
			</Col>
			<Col {...titleCol}>
				<Form.Item
					label="Name"
					name={['borrower', 'name']}
					rules={[
						{
							required: true,
						},
						{
							pattern: nameRegex,
							message: 'Cannot start or end with a delimiter',
						},
					]}
					getValueFromEvent={(v) => capitalizeFirstLetter(v.target.value)}
				>
					<Input placeholder="Enter Borrower name" />
				</Form.Item>
			</Col>
			<Col {...titleCol}>
				<Form.Item
					label="Phone"
					name={['borrower', 'phone']}
					rules={[
						{
							required: true,
						},
					]}
				>
					<InputPhone />
				</Form.Item>
			</Col>
			<Col {...titleCol}>
				<Form.Item
					label="Email"
					name={['borrower', 'email']}
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
				<Form.Item
					label="FICO credit score"
					name={['borrower', 'est_credit_score']}
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
				<Form.Item
					label="Flipped a property?"
					name={['borrower', 'flip_property_count']}
					rules={[{ required: true }]}
				>
					<InputNumber
						style={{ width: '100%' }}
						min={0}
						// max={maxCreditScore}
						controls={false}
						placeholder="Number of properties"
					/>
				</Form.Item>
			</Col>
			<Col {...titleCol}>
				<Form.Item
					label="Purchased a rental property?"
					name={['borrower', 'purchase_rental_property_count']}
					rules={[{ required: true }]}
				>
					<InputNumber
						style={{ width: '100%' }}
						min={0}
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

'use client';
import { Col, Form, Input, InputNumber, Row, Select } from 'antd';
import React from 'react';
import styles from './index.module.scss';
import { useGetOptions } from '@/hooks/useGetOptions';
import AddressForm from '@/components/GoogleAutoComplete/AddressForm';
import { maxNumberOfUnits } from '@/lib/format';
import { loan_pending } from '@/enum/common';
import CommonFormDate from '@/components/CommonDate';
import InputNumAmount from '@/components/InputNumAmount';
const layoutCol = {
	xs: 24,
	sm: 24,
};
const CreateProperty = ({
	stage,
	stepBtn,
	gutter,
	addressProps = {},
	form,
	isRequestProperty,
	...colProps
}) => {
	const { propertyTypeOptions } = useGetOptions([
		{
			key: 'property_type',
			valueName: 'propertyTypeOptions',
		},
	]);
	return (
		<Row className={styles.create} gutter={gutter}>
			<Col {...layoutCol} {...(colProps || {})}>
				<Form.Item
					label="Property type"
					name={['loan_property', 'property_type']}
					rules={[{ required: true, message: 'Please select Property type!' }]}
				>
					<Select
						options={propertyTypeOptions}
						placeholder="Select Property type"
						// disabled={isRequestProperty}
					/>
				</Form.Item>
			</Col>
			<Col {...layoutCol} {...(colProps || {})}>
				<Form.Item
					label="Number of units"
					name={['loan_property', 'property_base', 'unit']}
					rules={[{ required: true, message: 'Please input!' }]}
				>
					<InputNumber
						style={{ width: '100%' }}
						min={1}
						max={maxNumberOfUnits}
						controls={false}
						placeholder="Enter number of unit"
						disabled={isRequestProperty}
					/>
				</Form.Item>
			</Col>
			<AddressForm
				{...addressProps}
				addressName={['loan_property', 'property_base', 'address', 'address1']}
				addressName2={['loan_property', 'property_base', 'address', 'address2']}
				cityName={['loan_property', 'property_base', 'address', 'city']}
				stateName={['loan_property', 'property_base', 'address', 'state']}
				zipName={['loan_property', 'property_base', 'address', 'zip_code']}
				form={form}
				disabled={isRequestProperty}
			/>
			{stage === loan_pending && (
				<>
					<Col {...layoutCol} {...(colProps || {})}>
						<Form.Item label="Purchase date" name="purchase_date">
							<CommonFormDate />
						</Form.Item>
					</Col>
					<Col {...layoutCol} {...(colProps || {})}>
						<Form.Item label="owner" name={['loan_property', 'owner']}>
							<Input placeholder="Enter owner" />
						</Form.Item>
					</Col>
					<Col {...layoutCol} {...(colProps || {})}>
						<Form.Item label="Purchase price" name="purchase_amount">
							<InputNumAmount placeholder="Enter purchase price" />
						</Form.Item>
					</Col>
				</>
			)}
			{stepBtn}
		</Row>
	);
};
export default CreateProperty;

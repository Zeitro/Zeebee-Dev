'use client';
import { Col, Form, InputNumber, Row, Select } from 'antd';
import React from 'react';
import styles from './index.module.scss';
import { useGetOptions } from '@/hooks/useGetOptions';
import AddressForm from '@/components/GoogleAutoComplete/AddressForm';
import { maxNumberOfUnits } from '@/lib/format';
const layoutCol = {
	xs: 24,
	sm: 24,
};
const CreateProperty = ({
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
			{stepBtn}
		</Row>
	);
};
export default CreateProperty;

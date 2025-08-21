import React from 'react';
import styles from './index.module.scss';
import { Row, Col, Form, Radio, Select, Input, InputNumber } from 'antd';
import AddressForm from '@/components/GoogleAutoComplete/AddressForm';
import { useGetOptions } from '@/hooks/useGetOptions';
const PropertyInfo = ({ gutter = 24, form, addressProps }) => {
	const { loanPurposeOptions, loanTypeOptions, propertyTypeOptions } = useGetOptions([
		{
			key: 'loan_purpose',
			valueName: 'loanPurposeOptions',
		},
		{
			key: 'loan_type',
			valueName: 'loanTypeOptions',
		},
		{
			key: 'property_type',
			valueName: 'propertyTypeOptions',
		},
	]);
	const borrowerType = Form.useWatch('borrower_type', form);
	const layoutCol = {
		xs: 24,
		sm: 24,
	};

	return (
		<Row className={styles.createZeeBee} gutter={gutter}>
			<Col {...layoutCol} className="!text-2xl/8 font-medium text-[#000000a6] pb-8">
				Property information
			</Col>
			<AddressForm
				{...addressProps}
				addressName={['loan_property', 'property_base', 'address', 'address1']}
				addressName2={['loan_property', 'property_base', 'address', 'address2']}
				cityName={['loan_property', 'property_base', 'address', 'city']}
				stateName={['loan_property', 'property_base', 'address', 'state']}
				zipName={['loan_property', 'property_base', 'address', 'zip_code']}
				addressLabel="Subject property address 1"
				addressLabel2="Subject property address 2"
				form={form}
			/>
			<Col {...layoutCol}>
				<Form.Item
					label="Loan purpose"
					name="loan_purpose"
					rules={[{ required: true, message: 'Please select loan purpose!' }]}
				>
					<Select options={loanPurposeOptions} placeholder="Select loan purpose" />
				</Form.Item>
			</Col>
			<Col {...layoutCol}>
				<Form.Item
					label="Loan type"
					name="loan_type"
					rules={[{ required: true, message: 'Please select loan type!' }]}
				>
					<Select options={loanTypeOptions} placeholder="Select loan type" />
				</Form.Item>
			</Col>
			<Col {...layoutCol}>
				<Form.Item
					label="Property type"
					name="property_type"
					rules={[{ required: true, message: 'Please select property type!' }]}
				>
					<Select options={propertyTypeOptions} placeholder="Select property type" />
				</Form.Item>
			</Col>
			<Col {...layoutCol} md={12}>
				<Form.Item
					label="Living area sqft"
					name="name"
					rules={[
						{
							required: true,
							message: 'Please enter Living area sqft!',
						},
					]}
				>
					<Input placeholder="Enter Living area sqft" />
				</Form.Item>
			</Col>
			<Col {...layoutCol} md={12}>
				<Form.Item
					label="Number of units"
					name="name"
					rules={[
						{
							required: true,
							message: 'Please enter Number of units!',
						},
					]}
				>
					<InputNumber
						style={{ width: '100%' }}
						min={1}
						max={10}
						controls={false}
						placeholder="Enter number of unit"
					/>
				</Form.Item>
			</Col>
			<Col {...layoutCol} md={12}>
				<Form.Item
					label="Number of bedrooms"
					name="name"
					rules={[
						{
							required: true,
							message: 'Please enter Number of bedrooms!',
						},
					]}
				>
					<InputNumber
						style={{ width: '100%' }}
						min={1}
						max={10}
						controls={false}
						placeholder="Enter Number of bedrooms"
					/>
				</Form.Item>
			</Col>
			<Col {...layoutCol} md={12}>
				<Form.Item
					label="Number of bathrooms"
					name="name"
					rules={[
						{
							required: true,
							message: 'Please enter Number of bathrooms!',
						},
					]}
				>
					<InputNumber
						style={{ width: '100%' }}
						min={1}
						max={10}
						controls={false}
						placeholder="Enter Number of bathrooms"
					/>
				</Form.Item>
			</Col>
		</Row>
	);
};
export default PropertyInfo;
